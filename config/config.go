package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	IPLimitPerSecond     int           `mapstructure:"RATE_LIMIT_IP_DEFAULT"`
	TokenLimitPerSecond  int           `mapstructure:"RATE_LIMIT_TOKENS"`
	IPLockoutDuration    time.Duration `mapstructure:"RATE_LIMIT_IP_LOCKOUT_DURATION"`
	TokenLockoutDuration time.Duration `mapstructure:"RATE_LIMIT_TOKEN_LOCKOUT_DURATION"`
	RedisHost            string        `mapstructure:"REDIS_HOST"`
	RedisPort            string        `mapstructure:"REDIS_PORT"`
	RedisPassword        string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB              int           `mapstructure:"REDIS_DB"`
}

// LoadConfig loads configuration from a .env file located at the specified path.
// It uses Viper to read and parse the configuration file, and automatically
// loads environment variables. If the .env file cannot be read or parsed,
// the function prints an error message to stderr and panics.
// Returns a pointer to the Config struct and an error if any occurred.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Add search paths
	viper.AddConfigPath(path)
	viper.AddConfigPath(".") // fallback to current directory
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: could not read .env file: %v\n", err)
		panic(err)
	}

	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to parse config: %v\n", err)
		panic(err)
	}

	return cfg, nil
}
