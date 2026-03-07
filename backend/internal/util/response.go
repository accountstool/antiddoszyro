package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/domain"
)

type Envelope struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message,omitempty"`
	Data       any               `json:"data,omitempty"`
	Pagination *domain.Pagination `json:"pagination,omitempty"`
}

func JSON(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Envelope{
		Success: status < http.StatusBadRequest,
		Message: message,
		Data:    data,
	})
}

func JSONPage(c *gin.Context, status int, message string, data any, pagination *domain.Pagination) {
	c.JSON(status, Envelope{
		Success:    status < http.StatusBadRequest,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func AbortError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Envelope{
		Success: false,
		Message: message,
	})
}
