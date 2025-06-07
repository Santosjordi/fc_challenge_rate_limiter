package config_test

import (
	"path/filepath"
	"testing"

	"github.com/santosjordi/fc_challenge_rate_limiter/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromDotEnv(t *testing.T) {
	// Set up the path to your test .env file
	envPath := filepath.Join("..", "test", "testdata")
	envFile := ".env.test"

	viper.Reset()                // Ensure clean state
	viper.SetConfigName(envFile) // without extension
	viper.SetConfigType("env")   // needed since .env doesn't imply a type
	viper.AddConfigPath(envPath) // where to look
	viper.AutomaticEnv()         // allow override with real env vars

	err := viper.ReadInConfig()
	require.NoError(t, err)

	cfg, err := config.LoadConfig("") // if LoadConfig handles viper directly
	require.NoError(t, err)

	require.Equal(t, 5, cfg.IPLimitPerSecond)
}
