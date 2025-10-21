package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigAndGetConfigs(t *testing.T) {
	// Determine the directory of the current file (config.go)
	dir, err := os.Getwd()
	assert.NoError(t, err)

	// Ensure the config path includes the project directory
	viper.Reset()
	viper.AddConfigPath(dir)
	os.Setenv("APP_ENV", "local") // or "prod" depending on the file
	defer os.Unsetenv("APP_ENV")

	// Load the config
	LoadConfig()
	cfg := GetConfigs()

	// Validate app config
	assert.Equal(t, 8080, cfg.AppConfig.Port)
	assert.Equal(t, "local", cfg.AppConfig.Env)

	// Validate database config
	assert.Equal(t, "postgres", cfg.DatabaseConfig.Host)
	assert.Equal(t, 5432, cfg.DatabaseConfig.Port)
	assert.Equal(t, "postgres", cfg.DatabaseConfig.User)
	assert.Equal(t, "postgres", cfg.DatabaseConfig.Password)
	assert.Equal(t, "postgres", cfg.DatabaseConfig.Name)
	assert.Equal(t, "disable", cfg.DatabaseConfig.SSLMode)

	// Validate logger config
	assert.Equal(t, "info", cfg.LoggerConfig.Level)
	assert.Equal(t, "json", cfg.LoggerConfig.Encoding)
	assert.Equal(t, true, cfg.LoggerConfig.Development)
}
