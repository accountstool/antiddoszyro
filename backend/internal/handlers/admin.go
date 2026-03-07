package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

type settingsPayload struct {
	Values map[string]string `json:"values" binding:"required"`
}

type userPayload struct {
	Username    string          `json:"username" binding:"required,min=3"`
	Email       string          `json:"email" binding:"required,email"`
	DisplayName string          `json:"displayName" binding:"required"`
	Role        domain.UserRole `json:"role" binding:"required"`
	Language    string          `json:"language" binding:"required,oneof=en vi"`
	Password    string          `json:"password"`
}

type ipEntryPayload struct {
	DomainID  *uuid.UUID       `json:"domainId"`
	IP        string           `json:"ip"`
	CIDR      string           `json:"cidr"`
	Reason    string           `json:"reason"`
	ExpiresAt *time.Time       `json:"expiresAt"`
}

func (a *API) ListSettings(c *gin.Context) {
	items, err := a.settings.List(c.Request.Context())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "ok", items)
}

func (a *API) UpdateSettings(c *gin.Context) {
	var payload settingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid settings payload")
		return
	}
	if err := a.settings.Update(c.Request.Context(), userFromContext(c), payload.Values, c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "updated", payload.Values)
}

func (a *API) ListUsers(c *gin.Context) {
	limit, offset, meta := util.PaginationFromRequest(c, 20)
	items, total, err := a.users.List(c.Request.Context(), limit, offset)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := util.WithPagination(meta, total)
	util.JSONPage(c, http.StatusOK, "ok", items, &page)
}

func (a *API) CreateUser(c *gin.Context) {
	var payload userPayload
	if err := c.ShouldBindJSON(&payload); err != nil || len(payload.Password) < 8 {
		util.AbortError(c, http.StatusBadRequest, "invalid user payload")
		return
	}
	item, err := a.users.Create(c.Request.Context(), userFromContext(c), payload.Username, payload.Email, payload.DisplayName, payload.Role, payload.Language, payload.Password, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusCreated, "created", item)
}

func (a *API) UpdateUser(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	var payload userPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid user payload")
		return
	}
	item, err := a.users.Update(c.Request.Context(), userFromContext(c), store.UpdateUserParams{
		ID:          id,
		Username:    payload.Username,
		Email:       payload.Email,
		DisplayName: payload.DisplayName,
		Role:        payload.Role,
		Language:    payload.Language,
	}, payload.Password, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "updated", item)
}

func (a *API) DeleteUser(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	if err := a.users.Delete(c.Request.Context(), userFromContext(c), id, c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "deleted", gin.H{"id": id})
}

func (a *API) ListIPEntries(c *gin.Context, listType domain.ListType) {
	limit, offset, meta := util.PaginationFromRequest(c, 20)
	var domainID *uuid.UUID
	if value := c.Query("domainId"); value != "" {
		if parsed, err := uuid.Parse(value); err == nil {
			domainID = &parsed
		}
	}
	items, total, err := a.ipControl.ListEntries(c.Request.Context(), listType, domainID, c.Query("search"), limit, offset)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := util.WithPagination(meta, total)
	util.JSONPage(c, http.StatusOK, "ok", items, &page)
}

func (a *API) CreateIPEntry(c *gin.Context, listType domain.ListType) {
	var payload ipEntryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid ip entry payload")
		return
	}
	user := userFromContext(c)
	var createdBy *uuid.UUID
	if user != nil {
		createdBy = &user.ID
	}
	item, err := a.ipControl.CreateEntry(c.Request.Context(), user, store.IPEntryInput{
		DomainID:  payload.DomainID,
		ListType:  listType,
		IP:        payload.IP,
		CIDR:      payload.CIDR,
		Reason:    payload.Reason,
		ExpiresAt: payload.ExpiresAt,
		CreatedBy: createdBy,
	}, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusCreated, "created", item)
}

func (a *API) DeleteIPEntry(c *gin.Context) {
	id, ok := mustUUID(c, "id")
	if !ok {
		return
	}
	if err := a.ipControl.DeleteEntry(c.Request.Context(), userFromContext(c), id, c.ClientIP(), c.Request.UserAgent()); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.JSON(c, http.StatusOK, "deleted", gin.H{"id": id})
}

func (a *API) TemporaryBans(c *gin.Context) {
	limit, offset, meta := util.PaginationFromRequest(c, 20)
	items, total, err := a.ipControl.ListBans(c.Request.Context(), limit, offset)
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	page := util.WithPagination(meta, total)
	util.JSONPage(c, http.StatusOK, "ok", items, &page)
}
