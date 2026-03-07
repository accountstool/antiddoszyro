package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/database"
	"shieldpanel/backend/internal/handlers"
	"shieldpanel/backend/internal/middleware"
	"shieldpanel/backend/internal/nginx"
	"shieldpanel/backend/internal/protection"
	"shieldpanel/backend/internal/router"
	"shieldpanel/backend/internal/services"
	"shieldpanel/backend/internal/stats"
	"shieldpanel/backend/internal/store"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	cfg := config.Load()

	db, err := database.NewPostgres(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer db.Close()

	redisClient, err := database.NewRedis(ctx, cfg.Redis)
	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer redisClient.Close()

	repo := store.New(db, redisClient, logger)
	migrator := store.NewMigrator(db, cfg.Paths.MigrationsDir)
	if err := migrator.Run(ctx); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	_ = repo.CleanupExpiredSessions(ctx)

	sink := stats.NewSink(repo, logger)
	go sink.Start(ctx)

	auditService := services.NewAuditService(repo)
	authService := services.NewAuthService(repo, logger, cfg)
	nginxManager := nginx.NewManager(cfg.Nginx, logger)
	domainService := services.NewDomainService(repo, nginxManager, auditService, logger)
	dashboardService := services.NewDashboardService(repo)
	settingsService := services.NewSettingsService(repo, auditService)
	usersService := services.NewUsersService(repo, auditService, cfg)
	ipControlService := services.NewIPControlService(repo, auditService)
	logsService := services.NewLogsService(repo)
	sslService := services.NewSSLService(repo, cfg, auditService, domainService)
	protectionEngine := protection.NewEngine(repo, logger, cfg, sink)

	api := handlers.NewAPI(cfg, repo, authService, auditService, dashboardService, domainService, settingsService, usersService, ipControlService, logsService, sslService, protectionEngine, sink)
	engine := router.New(api, cfg.Server.FrontendDistDir, middleware.RequireAuth(authService), middleware.RequireCSRF(authService))

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	logger.Info("shieldpanel backend listening", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
