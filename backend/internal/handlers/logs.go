package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

func (a *API) StatisticsOverview(c *gin.Context) {
	from, to := parseTimeRange(c)
	var domainID *uuid.UUID
	if value := c.Query("domainId"); value != "" {
		if parsed, err := uuid.Parse(value); err == nil {
			domainID = &parsed
		}
	}
	overview, err := a.logs.Overview(c.Request.Context(), store.LogFilters{
		DomainID: domainID,
		From:     from,
		To:       to,
		Decision: c.Query("decision"),
		Reason:   c.Query("reason"),
	})
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ok", overview)
}

func (a *API) RequestLogs(c *gin.Context) {
	limit, offset, meta := util.PaginationFromRequest(c, 25)
	from, to := parseTimeRange(c)
	var domainID *uuid.UUID
	if value := c.Query("domainId"); value != "" {
		if parsed, err := uuid.Parse(value); err == nil {
			domainID = &parsed
		}
	}
	items, total, err := a.logs.List(c.Request.Context(), store.LogFilters{
		DomainID: domainID,
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
