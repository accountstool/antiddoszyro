package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/middleware"
	"shieldpanel/backend/internal/util"
)

func userFromContext(c *gin.Context) *domain.User {
	return middleware.CurrentUser(c)
}

func requestIDFromContext(c *gin.Context) string {
	value, _ := c.Get("request_id")
	if value == nil {
		return uuid.NewString()
	}
	return value.(string)
}

func mustUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	parsed, err := uuid.Parse(c.Param(name))
	if err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid id")
		return uuid.Nil, false
	}
	return parsed, true
}

func parseTimeRange(c *gin.Context) (time.Time, time.Time) {
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()
	if value := c.Query("from"); value != "" {
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			from = parsed
		}
	}
	if value := c.Query("to"); value != "" {
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			to = parsed
		}
	}
	return from, to
}
