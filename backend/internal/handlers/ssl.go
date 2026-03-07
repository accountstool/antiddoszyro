package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shieldpanel/backend/internal/util"
)

type sslPayload struct {
	DomainID string `json:"domainId" binding:"required"`
}

func (a *API) IssueSSL(c *gin.Context) {
	var payload sslPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid ssl payload")
		return
	}
	id, err := uuid.Parse(payload.DomainID)
	if err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid domain id")
		return
	}
	if err := a.ssl.Issue(c.Request.Context(), userFromContext(c), id, c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ssl issued", gin.H{"domainId": id})
}

func (a *API) RenewSSL(c *gin.Context) {
	var payload sslPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid ssl payload")
		return
	}
	id, err := uuid.Parse(payload.DomainID)
	if err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid domain id")
		return
	}
	if err := a.ssl.Renew(c.Request.Context(), userFromContext(c), id, c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ssl renewed", gin.H{"domainId": id})
}
