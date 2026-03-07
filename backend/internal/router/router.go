package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/handlers"
	"shieldpanel/backend/internal/middleware"
)

func New(api *handlers.API, frontendDistDir string, authMiddleware gin.HandlerFunc, csrfMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.SecurityHeaders())

	r.GET("/healthz", api.Healthz)
	r.POST("/api/auth/login", api.Login)
	r.GET("/internal/protection/check", api.ProtectionCheck)
	r.GET("/internal/protection/challenge", api.ProtectionChallenge)
	r.GET("/internal/protection/block", api.ProtectionBlock)
	r.POST("/internal/protection/verify", api.ProtectionVerify)

	authenticated := r.Group("/api")
	authenticated.Use(authMiddleware)
	authenticated.GET("/auth/me", api.Me)
	authenticated.POST("/auth/logout", csrfMiddleware, api.Logout)
	authenticated.GET("/dashboard/summary", api.DashboardSummary)
	authenticated.GET("/dashboard/charts", api.DashboardCharts)
	authenticated.GET("/domains", api.ListDomains)
	authenticated.GET("/domains/:id", api.GetDomain)
	authenticated.GET("/domains/:id/stats", api.DomainStats)
	authenticated.GET("/domains/:id/logs", api.DomainLogs)
	authenticated.GET("/logs", api.RequestLogs)
	authenticated.GET("/statistics/overview", api.StatisticsOverview)
	authenticated.GET("/settings", api.ListSettings)
	authenticated.GET("/users", api.ListUsers)
	authenticated.GET("/audit-logs", api.AuditLogs)
	authenticated.GET("/ip-control/blacklist", func(c *gin.Context) { api.ListIPEntries(c, "blacklist") })
	authenticated.GET("/ip-control/whitelist", func(c *gin.Context) { api.ListIPEntries(c, "whitelist") })
	authenticated.GET("/ip-control/bans", api.TemporaryBans)

	mutating := authenticated.Group("")
	mutating.Use(csrfMiddleware)
	mutating.POST("/domains", api.CreateDomain)
	mutating.PUT("/domains/:id", api.UpdateDomain)
	mutating.DELETE("/domains/:id", api.DeleteDomain)
	mutating.POST("/nginx/reload", api.ReloadNginx)
	mutating.POST("/ssl/issue", api.IssueSSL)
	mutating.POST("/ssl/renew", api.RenewSSL)
	mutating.PUT("/settings", api.UpdateSettings)
	mutating.POST("/users", api.CreateUser)
	mutating.PUT("/users/:id", api.UpdateUser)
	mutating.DELETE("/users/:id", api.DeleteUser)
	mutating.POST("/ip-control/blacklist", func(c *gin.Context) { api.CreateIPEntry(c, "blacklist") })
	mutating.POST("/ip-control/whitelist", func(c *gin.Context) { api.CreateIPEntry(c, "whitelist") })
	mutating.DELETE("/ip-control/:id", api.DeleteIPEntry)

	serveFrontend(r, frontendDistDir)
	return r
}

func serveFrontend(r *gin.Engine, frontendDistDir string) {
	if stat, err := os.Stat(frontendDistDir); err != nil || !stat.IsDir() {
		return
	}
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || strings.HasPrefix(c.Request.URL.Path, "/internal/") {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
			return
		}
		target := filepath.Join(frontendDistDir, strings.TrimPrefix(filepath.Clean(c.Request.URL.Path), string(filepath.Separator)))
		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			c.File(target)
			return
		}
		c.File(filepath.Join(frontendDistDir, "index.html"))
	})
}
