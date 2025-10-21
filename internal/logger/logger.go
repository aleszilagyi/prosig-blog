package logger

import (
	"strings"

	"github.com/aleszilagyi/prosig-blog/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger initializes a zap logger using configs
func GetLogger() *zap.Logger {
	cfg := config.GetConfigs().LoggerConfig
	var zapCfg zap.Config
	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	zapCfg.Encoding = cfg.Encoding
	zapCfg.Level = zap.NewAtomicLevelAt(parseLevel(cfg.Level))

	logger, err := zapCfg.Build()
	if err != nil {
		zap.NewExample().Fatal("[Logger] Failed to build zap logger", zap.Error(err),
			zap.String("logger_encoding", cfg.Encoding),
			zap.String("logger_level", cfg.Level),
			zap.Bool("is_development", cfg.Development),
		)
	}
	return logger
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
