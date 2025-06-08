package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration parameters for the rate limiter service.
// It includes settings for request limits, lockout durations, window size, and Redis connection details.
// Fields:
//   - IPRequestPerSecond: Maximum number of requests allowed per IP address per second.
//   - TokenRequestPerSecond: Maximum number of requests allowed per token per second.
//   - IPLockoutDuration: Duration (in seconds) for which an IP address is locked out after exceeding its request limit.
//   - TokenLockoutDuration: Duration (in seconds) for which a token is locked out after exceeding its request limit.
//   - RedisHost: Host address of the Redis server.
//   - RedisPort: Port number of the Redis server.
//   - RedisPassword: Password for authenticating with the Redis server.
//   - RedisDB: Redis database number to use.
//   - ServerPort: Port on which the server listens for incoming requests.
type Config struct {
	IPRequestPerSecond    int           `mapstructure:"RATE_LIMIT_IP_DEFAULT"`
	TokenRequestPerSecond int           `mapstructure:"RATE_LIMIT_TOKENS"`
	IPLockoutDuration     time.Duration `mapstructure:"RATE_LIMIT_IP_LOCKOUT_DURATION"`
	TokenLockoutDuration  time.Duration `mapstructure:"RATE_LIMIT_TOKEN_LOCKOUT_DURATION"`
	RedisHost             string        `mapstructure:"REDIS_HOST"`
	RedisPort             string        `mapstructure:"REDIS_PORT"`
	RedisPassword         string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB               int           `mapstructure:"REDIS_DB"`
	ServerPort            string        `mapstructure:"SERVER_PORT"`
}

// LoadConfig loads configuration from a .env file located at the specified path.
// It uses Viper to read and parse the configuration file, and automatically
// loads environment variables. If the .env file cannot be read or parsed,
// the function prints an error message to stderr and panics.
// Returns a pointer to the Config struct and an error if any occurred.
func LoadConfig(envFilePath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(envFilePath) // ✅ Use full path like "/project/.env"
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read .env file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
