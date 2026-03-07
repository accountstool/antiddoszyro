package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/protection"
	"shieldpanel/backend/internal/util"
)

func (a *API) ProtectionCheck(c *gin.Context) {
	headers := map[string]string{}
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[strings.ToLower(key)] = values[0]
		}
	}
	cookies := map[string]string{}
	for _, cookie := range c.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	decision, err := a.engine.Evaluate(c.Request.Context(), protection.RequestContext{
		Host:      c.GetHeader("X-Original-Host"),
		URI:       c.GetHeader("X-Original-URI"),
		Path:      c.GetHeader("X-Original-Path"),
		Query:     c.GetHeader("X-Original-Query"),
		Method:    c.GetHeader("X-Original-Method"),
		Scheme:    c.GetHeader("X-Forwarded-Proto"),
		RemoteIP:  c.GetHeader("X-Real-IP"),
		UserAgent: c.Request.UserAgent(),
		Headers:   headers,
		Cookies:   cookies,
		RequestID: requestIDFromContext(c),
	})
	if err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("X-ShieldPanel-Reason", decision.Reason)
	c.Header("X-ShieldPanel-Challenge", decision.ChallengeMode)
	c.Status(decision.StatusCode)
}

func (a *API) ProtectionChallenge(c *gin.Context) {
	lang := c.GetHeader("Accept-Language")
	mode := c.GetHeader("X-Challenge-Mode")
	if mode == "" {
		mode = "cookie"
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusUnauthorized, protection.ChallengePage(lang, c.GetHeader("X-Original-Host"), c.GetHeader("X-Original-URI"), mode, c.GetHeader("X-Reason")))
}

func (a *API) ProtectionBlock(c *gin.Context) {
	lang := c.GetHeader("Accept-Language")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusForbidden, protection.BlockPage(lang, c.GetHeader("X-Reason")))
}

func (a *API) ProtectionVerify(c *gin.Context) {
	host := strings.ToLower(strings.TrimSpace(c.PostForm("host")))
	if host == "" {
		host = strings.ToLower(strings.TrimSpace(c.GetHeader("X-Original-Host")))
	}
	redirectURI := c.PostForm("redirect_uri")
	if redirectURI == "" {
		redirectURI = "/"
	}
	item, err := a.engine.ResolveDomain(c.Request.Context(), host)
	if err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid host")
		return
	}
	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}
	token := a.engine.IssueClearance(item.Name, clientIP)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("shieldpanel_clearance", token, int(a.cfg.Security.ChallengeGraceTTL.Seconds()), "/", "", a.cfg.Auth.CookieSecure, false)
	a.sink.Publish(domain.RequestLogInput{
		DomainID:       &item.ID,
		DomainName:     item.Name,
		ClientIP:       clientIP,
		CountryCode:    c.GetHeader("CF-IPCountry"),
		Method:         "POST",
		Path:           "/__shieldpanel_verify",
		QueryString:    "",
		UserAgent:      c.Request.UserAgent(),
		RequestID:      requestIDFromContext(c),
		Decision:       domain.RequestDecisionChallengePassed,
		StatusCode:     http.StatusFound,
		BlockReason:    "challenge_passed",
		ResponseTimeMS: 1,
		ChallengeType:  "cookie",
		OccurredAt:     time.Now(),
	})
	c.Redirect(http.StatusFound, redirectURI)
}
