package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/services"
	"shieldpanel/backend/internal/util"
)

const (
	ContextUserKey    = "auth_user"
	ContextSessionKey = "auth_session"
)

func RequireAuth(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(auth.SessionCookieName())
		if err != nil || token == "" {
			util.AbortError(c, http.StatusUnauthorized, "authentication required")
			return
		}
		session, user, err := auth.AuthenticateSession(c.Request.Context(), token)
		if err != nil {
			util.AbortError(c, http.StatusUnauthorized, "invalid session")
			return
		}
		c.Set(ContextUserKey, user)
		c.Set(ContextSessionKey, session)
		c.Next()
	}
}

func RequireCSRF(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}
		sessionValue, ok := c.Get(ContextSessionKey)
		if !ok {
			util.AbortError(c, http.StatusUnauthorized, "invalid session")
			return
		}
		session := sessionValue.(domain.Session)
		token := c.GetHeader("X-CSRF-Token")
		if token == "" || token != session.CSRFToken {
			util.AbortError(c, http.StatusForbidden, "invalid csrf token")
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *domain.User {
	value, ok := c.Get(ContextUserKey)
	if !ok {
		return nil
	}
	user := value.(domain.User)
	return &user
}
