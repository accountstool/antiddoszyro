package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/util"
)

func (a *API) DashboardSummary(c *gin.Context) {
	currentRPS, currentBlockedPS := a.engine.CurrentRates(c.Request.Context())
	summary, err := a.dashboard.Summary(c.Request.Context(), currentRPS, currentBlockedPS)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ok", summary)
}

func (a *API) DashboardCharts(c *gin.Context) {
	charts, err := a.dashboard.Charts(c.Request.Context())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ok", charts)
}
