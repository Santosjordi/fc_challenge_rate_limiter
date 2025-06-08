package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"github.com/santosjordi/fc_challenge_rate_limiter/config"
	db "github.com/santosjordi/fc_challenge_rate_limiter/internal/infra/ratelimiter"
	"github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/handlers"
	mw "github.com/santosjordi/fc_challenge_rate_limiter/internal/webserver/middleware"
	"github.com/stretchr/testify/suite"
)

type RateLimiterTestSuite struct {
	suite.Suite
	cfg       *config.Config
	client    *redis.Client
	store     *db.RedisStorage
	server    *httptest.Server
	ctx       context.Context
	testToken string
	testIP    string
}

func (suite *RateLimiterTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	suite.cfg = &config.Config{
		IPRequestPerSecond:    3,
		TokenRequestPerSecond: 10,
		IPLockoutDuration:     2 * time.Second,
		TokenLockoutDuration:  1 * time.Second,
		RedisHost:             "localhost",
		RedisPort:             "6379",
		RedisPassword:         "",
		RedisDB:               0,
		ServerPort:            "8080",
	}

	suite.client = redis.NewClient(&redis.Options{
		Addr:     suite.cfg.RedisHost + ":" + suite.cfg.RedisPort,
		Password: suite.cfg.RedisPassword,
		DB:       suite.cfg.RedisDB,
	})

	suite.store = db.NewRedisStorage(suite.client)
	rateLimiter := mw.NewRateLimitMiddleware(suite.store, suite.cfg)

	router := chi.NewRouter()
	router.Use(rateLimiter.Handler)
	router.Get("/generate", handlers.UuidHandler().ServeHTTP)

	suite.server = httptest.NewServer(router)

	suite.testToken = "test-token-123"
	suite.testIP = "127.0.0.1"

	// Reset before starting
	suite.store.Reset(suite.ctx, suite.testToken)
	suite.store.Reset(suite.ctx, suite.testIP)
}

func (suite *RateLimiterTestSuite) TearDownSuite() {
	suite.server.Close()
	suite.store.Reset(suite.ctx, suite.testToken)
	suite.store.Reset(suite.ctx, suite.testIP)
	suite.client.Close()
}

func (suite *RateLimiterTestSuite) sendRequestWithToken(token string) *http.Response {
	req, _ := http.NewRequest("GET", suite.server.URL+"/generate", nil)
	req.Header.Set("API_KEY", token)
	resp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	return resp
}

func (suite *RateLimiterTestSuite) sendRequestWithIP(ip string) *http.Response {
	req, _ := http.NewRequest("GET", suite.server.URL+"/generate", nil)
	req.RemoteAddr = ip + ":1234"
	resp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	return resp
}

func (suite *RateLimiterTestSuite) TestTokenRateLimitExceeded() {
	for i := 0; i < suite.cfg.TokenRequestPerSecond; i++ {
		resp := suite.sendRequestWithToken(suite.testToken)
		suite.Equal(http.StatusOK, resp.StatusCode)
	}

	resp := suite.sendRequestWithToken(suite.testToken)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)
}

func (suite *RateLimiterTestSuite) TestIPRateLimitExceeded() {
	for i := 0; i < suite.cfg.IPRequestPerSecond; i++ {
		resp := suite.sendRequestWithIP(suite.testIP)
		suite.Equal(http.StatusOK, resp.StatusCode)
	}

	resp := suite.sendRequestWithIP(suite.testIP)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)
}

func (suite *RateLimiterTestSuite) TestIPLockoutExpiry() {
	// Arrange
	for i := 0; i < suite.cfg.IPRequestPerSecond; i++ {
		resp := suite.sendRequestWithIP(suite.testIP)
		suite.Equal(http.StatusOK, resp.StatusCode)
	}

	// This should trigger lockout
	resp := suite.sendRequestWithIP(suite.testIP)
	suite.Equal(http.StatusTooManyRequests, resp.StatusCode)

	// Wait for lockout to expire
	waitTime := suite.cfg.IPLockoutDuration + time.Second
	suite.T().Logf("Waiting for %s before retrying", waitTime)
	time.Sleep(waitTime)

	// Try again
	resp = suite.sendRequestWithIP(suite.testIP)
	suite.Equal(http.StatusOK, resp.StatusCode)
}

func TestRateLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}
