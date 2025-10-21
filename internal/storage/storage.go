package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aleszilagyi/prosig-blog/config"
	log "github.com/aleszilagyi/prosig-blog/internal/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func Connect(cfg config.Config) (*sql.DB, error) {
	logger := log.GetLogger()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("[PGConnection] failed to open postgres connection", zap.Error(err),
			zap.String("db_name", cfg.DatabaseConfig.Name),
			zap.String("db_host", cfg.DatabaseConfig.Host),
			zap.Int("db_port", cfg.DatabaseConfig.Port),
		)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("[PGConnection] failed to ping postgres", zap.Error(err),
			zap.String("db_name", cfg.DatabaseConfig.Name),
			zap.String("db_host", cfg.DatabaseConfig.Host),
			zap.Int("db_port", cfg.DatabaseConfig.Port),
		)
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}
