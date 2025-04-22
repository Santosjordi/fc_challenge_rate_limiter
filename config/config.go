package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type Config struct {
	// Rate Limiting
	RateLimitIP          int `mapstructure:"RATE_LIMIT_IP_DEFAULT"`
	RateLimitToken       int `mapstructure:"RATE_LIMIT_TOKENS"`
	LockoutDurationIP    int `mapstructure:"RATE_LIMIT_IP_LOCKOUT_DURATION"`
	LockoutDurationToken int `mapstructure:"RATE_LIMIT_TOKEN_LOCKOUT_DURATION"`

	// Auth
	JWTSecret string           `mapstructure:"JWT_SECRET"`
	TokenAuth *jwtauth.JWTAuth `mapstructure:"-"`

	// Persistence
	RateLimitBackend string `mapstructure:"RATE_LIMIT_BACKEND"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// Server
	ServerPort int `mapstructure:"SERVER_PORT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Add search paths
	viper.AddConfigPath(path)
	viper.AddConfigPath(".") // fallback to current directory
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not read .env file: %v\n", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Normalize backend (lowercase to simplify usage)
	cfg.RateLimitBackend = strings.ToLower(cfg.RateLimitBackend)

	// Set up JWT token auth
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return &cfg, nil
}
