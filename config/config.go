package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	AppConfig      AppConfig      `mapstructure:"app"`
	DatabaseConfig DatabaseConfig `mapstructure:"db"`
	LoggerConfig   LoggerConfig   `mapstructure:"logger"`
}

type AppConfig struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type LoggerConfig struct {
	Level       string `mapstructure:"level"`
	Encoding    string `mapstructure:"encoding"`
	Development bool   `mapstructure:"development"`
}

var c Config

func LoadConfig() {
	logger := zap.NewExample()
	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "local"
	}
	logger = logger.With(zap.String("environment", env))
	logger.Info("[Config] loading configs")

	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(filename)

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Environment variables take precedence
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.BindEnv("env")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("[Config] error reading config file", zap.Error(err),
			zap.String("filename", filename),
			zap.String("config_path", configPath),
		)
	}

	if err := viper.Unmarshal(&c); err != nil {
		logger.Fatal("[Config] unable to decode config into struct", zap.Error(err),
			zap.String("filename", filename),
			zap.String("config_path", configPath),
		)
	}

	logger.Info("[Config] loaded configuration")
}

func GetConfigs() Config {
	return c
}
