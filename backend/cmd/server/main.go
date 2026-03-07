package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"shieldpanel/backend/internal/app"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx, logger); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
