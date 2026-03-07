package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/protection"
	"shieldpanel/backend/internal/services"
	"shieldpanel/backend/internal/stats"
	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

type API struct {
	cfg       config.Config
	store     *store.Store
	auth      *services.AuthService
	audit     *services.AuditService
	dashboard *services.DashboardService
	domains   *services.DomainService
	settings  *services.SettingsService
	users     *services.UsersService
	ipControl *services.IPControlService
	logs      *services.LogsService
	ssl       *services.SSLService
	engine    *protection.Engine
	sink      *stats.Sink
}

func NewAPI(
	cfg config.Config,
	repo *store.Store,
	auth *services.AuthService,
	audit *services.AuditService,
	dashboard *services.DashboardService,
	domains *services.DomainService,
	settings *services.SettingsService,
	users *services.UsersService,
	ipControl *services.IPControlService,
	logs *services.LogsService,
	ssl *services.SSLService,
	engine *protection.Engine,
	sink *stats.Sink,
) *API {
	return &API{
		cfg:       cfg,
		store:     repo,
		auth:      auth,
		audit:     audit,
		dashboard: dashboard,
		domains:   domains,
		settings:  settings,
		users:     users,
		ipControl: ipControl,
		logs:      logs,
		ssl:       ssl,
		engine:    engine,
		sink:      sink,
	}
}

func (a *API) Healthz(c *gin.Context) {
	if err := a.store.Ping(c.Request.Context()); err != nil {
		util.JSON(c, http.StatusServiceUnavailable, "unhealthy", gin.H{"status": "down"})
		return
	}
	util.JSON(c, http.StatusOK, "ok", gin.H{
		"status":  "up",
		"app":     a.cfg.App.Name,
		"version": a.cfg.App.Version,
		"time":    time.Now().UTC(),
	})
}

func (a *API) AuditLogs(c *gin.Context) {
	limit, offset, meta := util.PaginationFromRequest(c, 20)
	items, total, err := a.store.ListAuditLogs(c.Request.Context(), limit, offset)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := util.WithPagination(meta, total)
	util.JSONPage(c, http.StatusOK, "ok", items, &page)
}
