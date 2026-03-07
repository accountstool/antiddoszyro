package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/database"
	"shieldpanel/backend/internal/services"
	"shieldpanel/backend/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.Load()
	ctx := context.Background()

	db, err := database.NewPostgres(ctx, cfg.Database)
	if err != nil {
		logger.Error("failed to connect postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := database.NewRedis(ctx, cfg.Redis)
	if err != nil {
		logger.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	repo := store.New(db, redisClient, logger)
	seeder := services.NewSeeder(repo, logger, cfg)
	result, err := seeder.Seed(ctx)
	if err != nil {
		logger.Error("failed to seed system", "error", err)
		os.Exit(1)
	}

	fmt.Printf("seed completed: admin=%s sample_domain=%s\n", result.AdminIdentifier, result.SampleDomain)
}
