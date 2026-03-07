package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/database"
	"shieldpanel/backend/internal/store"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.Load()

	db, err := database.NewPostgres(context.Background(), cfg.Database)
	if err != nil {
		logger.Error("failed to connect postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrations := store.NewMigrator(db, cfg.Paths.MigrationsDir)
	if err := migrations.Run(context.Background()); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	fmt.Println("migrations applied successfully")
}
