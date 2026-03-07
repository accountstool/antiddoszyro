package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"shieldpanel/backend/internal/config"
)

func NewPostgres(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdle
	poolConfig.MaxConnLifetime = cfg.MaxConnLife
	poolConfig.HealthCheckPeriod = cfg.HealthPeriod
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second
	return pgxpool.NewWithConfig(ctx, poolConfig)
}
