package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/middleware"
	"shieldpanel/backend/internal/services"
	"shieldpanel/backend/internal/util"
)

type loginRequest struct {
	Identifier string `json:"identifier" binding:"required,min=3"`
	Password   string `json:"password" binding:"required,min=8"`
	RememberMe bool   `json:"rememberMe"`
}

func (a *API) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid login payload")
		return
	}

	result, err := a.auth.Login(c.Request.Context(), request.Identifier, request.Password, request.RememberMe, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrRateLimited):
			util.AbortError(c, http.StatusTooManyRequests, "too many login attempts")
		default:
			util.AbortError(c, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	maxAge := int(a.auth.SessionDuration(request.RememberMe).Seconds())
	c.SetCookie(a.auth.SessionCookieName(), result.SessionToken, maxAge, "/", a.auth.CookieDomain(), a.auth.SetSessionCookiesSecure(), true)
	c.SetCookie(a.auth.CSRFCookieName(), result.CSRFToken, maxAge, "/", a.auth.CookieDomain(), a.auth.SetSessionCookiesSecure(), false)
	a.audit.Record(c.Request.Context(), &result.User, "login", "session", result.User.ID.String(), c.ClientIP(), c.Request.UserAgent(), "")
	util.JSON(c, http.StatusOK, "ok", gin.H{
		"user": result.User,
	})
}

func (a *API) Logout(c *gin.Context) {
	user := middleware.CurrentUser(c)
	token, _ := c.Cookie(a.auth.SessionCookieName())
	_ = a.auth.Logout(c.Request.Context(), token)
	c.SetCookie(a.auth.SessionCookieName(), "", -1, "/", a.auth.CookieDomain(), a.auth.SetSessionCookiesSecure(), true)
	c.SetCookie(a.auth.CSRFCookieName(), "", -1, "/", a.auth.CookieDomain(), a.auth.SetSessionCookiesSecure(), false)
	if user != nil {
		a.audit.Record(c.Request.Context(), user, "logout", "session", user.ID.String(), c.ClientIP(), c.Request.UserAgent(), "")
	}
	util.JSON(c, http.StatusOK, "ok", gin.H{"loggedOut": true})
}

func (a *API) Me(c *gin.Context) {
	user := middleware.CurrentUser(c)
	util.JSON(c, http.StatusOK, "ok", gin.H{
		"user": user,
	})
}
