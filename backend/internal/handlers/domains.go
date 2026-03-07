package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

type domainRulePayload struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Pattern  string `json:"pattern" binding:"required"`
	Action   string `json:"action" binding:"required"`
	Enabled  bool   `json:"enabled"`
	Priority int    `json:"priority"`
}

type domainPayload struct {
	Name               string              `json:"name" binding:"required,hostname_rfc1123"`
	OriginHost         string              `json:"originHost" binding:"required"`
	OriginPort         int                 `json:"originPort" binding:"required,min=1,max=65535"`
	OriginProtocol     string              `json:"originProtocol" binding:"required,oneof=http https"`
	OriginServerName   string              `json:"originServerName"`
	Enabled            bool                `json:"enabled"`
	ProtectionEnabled  bool                `json:"protectionEnabled"`
	ProtectionMode     domain.ProtectionMode `json:"protectionMode" binding:"required"`
	ChallengeMode      domain.ChallengeMode  `json:"challengeMode" binding:"required"`
	CloudflareMode     bool                `json:"cloudflareMode"`
	SSLAutoIssue       bool                `json:"sslAutoIssue"`
	SSLEnabled         bool                `json:"sslEnabled"`
	ForceHTTPS         bool                `json:"forceHttps"`
	RateLimitRPS       int                 `json:"rateLimitRps" binding:"required,min=1,max=10000"`
	RateLimitBurst     int                 `json:"rateLimitBurst" binding:"required,min=1,max=10000"`
	BadBotMode         bool                `json:"badBotMode"`
	HeaderValidation   bool                `json:"headerValidation"`
	JSChallengeEnabled bool                `json:"jsChallengeEnabled"`
	AllowedMethods     []string            `json:"allowedMethods" binding:"required,min=1"`
	Notes              string              `json:"notes"`
	Rules              []domainRulePayload `json:"rules"`
}

func (a *API) ListDomains(c *gin.Context) {
	limit, offset, meta := util.PaginationFromRequest(c, 20)
	items, total, err := a.domains.List(c.Request.Context(), limit, offset, c.Query("search"))
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := util.WithPagination(meta, total)
	util.JSONPage(c, http.StatusOK, "ok", items, &page)
}

func (a *API) GetDomain(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	item, err := a.domains.Get(c.Request.Context(), id)
	if err != nil {
		util.AbortError(c, http.StatusNotFound, "domain not found")
		return
	}
	util.JSON(c, http.StatusOK, "ok", item)
}

func (a *API) CreateDomain(c *gin.Context) {
	var payload domainPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid domain payload")
		return
	}
	item, err := a.domains.Create(c.Request.Context(), userFromContext(c), store.DomainInput{
		Name:               payload.Name,
		OriginHost:         payload.OriginHost,
		OriginPort:         payload.OriginPort,
		OriginProtocol:     payload.OriginProtocol,
		OriginServerName:   payload.OriginServerName,
		Enabled:            payload.Enabled,
		ProtectionEnabled:  payload.ProtectionEnabled,
		ProtectionMode:     payload.ProtectionMode,
		ChallengeMode:      payload.ChallengeMode,
		CloudflareMode:     payload.CloudflareMode,
		SSLAutoIssue:       payload.SSLAutoIssue,
		SSLEnabled:         payload.SSLEnabled,
		ForceHTTPS:         payload.ForceHTTPS,
		RateLimitRPS:       payload.RateLimitRPS,
		RateLimitBurst:     payload.RateLimitBurst,
		BadBotMode:         payload.BadBotMode,
		HeaderValidation:   payload.HeaderValidation,
		JSChallengeEnabled: payload.JSChallengeEnabled,
		AllowedMethods:     payload.AllowedMethods,
		Notes:              payload.Notes,
	}, mapDomainRules(payload.Rules), c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusCreated, "created", item)
}

func (a *API) UpdateDomain(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	var payload domainPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid domain payload")
		return
	}
	item, err := a.domains.Update(c.Request.Context(), userFromContext(c), id, store.DomainInput{
		Name:               payload.Name,
		OriginHost:         payload.OriginHost,
		OriginPort:         payload.OriginPort,
		OriginProtocol:     payload.OriginProtocol,
		OriginServerName:   payload.OriginServerName,
		Enabled:            payload.Enabled,
		ProtectionEnabled:  payload.ProtectionEnabled,
		ProtectionMode:     payload.ProtectionMode,
		ChallengeMode:      payload.ChallengeMode,
		CloudflareMode:     payload.CloudflareMode,
		SSLAutoIssue:       payload.SSLAutoIssue,
		SSLEnabled:         payload.SSLEnabled,
		ForceHTTPS:         payload.ForceHTTPS,
		RateLimitRPS:       payload.RateLimitRPS,
		RateLimitBurst:     payload.RateLimitBurst,
		BadBotMode:         payload.BadBotMode,
		HeaderValidation:   payload.HeaderValidation,
		JSChallengeEnabled: payload.JSChallengeEnabled,
		AllowedMethods:     payload.AllowedMethods,
		Notes:              payload.Notes,
	}, mapDomainRules(payload.Rules), c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "updated", item)
}

func (a *API) DeleteDomain(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	if err := a.domains.Delete(c.Request.Context(), userFromContext(c), id, c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "deleted", gin.H{"id": id})
}

func (a *API) DomainStats(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	from, to := parseTimeRange(c)
	stats, err := a.logs.DomainOverview(c.Request.Context(), id, from, to)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ok", stats)
}

func (a *API) DomainLogs(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	limit, offset, meta := util.PaginationFromRequest(c, 20)
	from, to := parseTimeRange(c)
	items, total, err := a.logs.List(c.Request.Context(), store.LogFilters{
		DomainID: &id,
		From:     from,
		To:       to,
		Decision: c.Query("decision"),
		Reason:   c.Query("reason"),
	}, limit, offset)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := util.WithPagination(meta, total)
	util.JSONPage(c, http.StatusOK, "ok", items, &page)
}

func (a *API) ReloadNginx(c *gin.Context) {
	if err := a.domains.Sync(c.Request.Context(), userFromContext(c), c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "reloaded", gin.H{"reloaded": true})
}

func mapDomainRules(rules []domainRulePayload) []store.DomainRuleInput {
	items := make([]store.DomainRuleInput, 0, len(rules))
	for _, rule := range rules {
		items = append(items, store.DomainRuleInput{
			Name:     rule.Name,
			Type:     rule.Type,
			Pattern:  rule.Pattern,
			Action:   rule.Action,
			Enabled:  rule.Enabled,
			Priority: rule.Priority,
		})
	}
	return items
}
