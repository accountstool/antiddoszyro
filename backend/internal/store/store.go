package store

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	db     *pgxpool.Pool
	redis  *redis.Client
	logger *slog.Logger
}

func New(db *pgxpool.Pool, redisClient *redis.Client, logger *slog.Logger) *Store {
	return &Store{
		db:     db,
		redis:  redisClient,
		logger: logger,
	}
}

func (s *Store) DB() *pgxpool.Pool {
	return s.db
}

func (s *Store) Redis() *redis.Client {
	return s.redis
}

func (s *Store) Ping(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return err
	}
	return s.redis.Ping(ctx).Err()
}
