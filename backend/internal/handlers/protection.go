package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/protection"
	"shieldpanel/backend/internal/util"
)

const (
	minChallengeInteractionMS = 900
	minChallengePointerMoves  = 8
)

type browserChallengeSession struct {
	Domain       string `json:"domain"`
	ClientIP     string `json:"clientIp"`
	UserAgentSum string `json:"userAgentSum"`
	RedirectURI  string `json:"redirectUri"`
	IssuedAtUnix int64  `json:"issuedAtUnix"`
	Mode         string `json:"mode"`
}

type browserChallengeProof struct {
	Nonce           string
	Mode            string
	SliderCompleted bool
	InteractionMS   int
	PointerMoves    int
	WebDriver       bool
	Language        string
	Timezone        string
	Screen          string
	Viewport        string
	Platform        string
	DeviceMemory    string
	CPUCount        string
	Honeypot        string
}

func (a *API) ProtectionCheck(c *gin.Context) {
	headers := requestHeaders(c)
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
	headers := requestHeaders(c)
	host := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Original-Host")))
	redirectURI := sanitizeRedirectURI(c.GetHeader("X-Original-URI"))
	mode := normalizeChallengeMode(c.GetHeader("X-Challenge-Mode"))
	if err := a.renderBrowserChallenge(c, headers, host, redirectURI, mode, c.GetHeader("X-Reason")); err != nil {
		util.AbortError(c, http.StatusInternalServerError, err.Error())
	}
}

func (a *API) ProtectionBlock(c *gin.Context) {
	lang := c.GetHeader("Accept-Language")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusForbidden, protection.BlockPage(lang, c.GetHeader("X-Reason")))
}

func (a *API) ProtectionVerify(c *gin.Context) {
	headers := requestHeaders(c)
	host := strings.ToLower(strings.TrimSpace(c.PostForm("host")))
	if host == "" {
		host = strings.ToLower(strings.TrimSpace(c.GetHeader("X-Original-Host")))
	}
	redirectURI := sanitizeRedirectURI(c.PostForm("redirect_uri"))
	mode := normalizeChallengeMode(c.PostForm("challenge_mode"))
	item, clientIP, err := a.resolveChallengeTarget(c, headers, host)
	if err != nil {
		util.AbortError(c, http.StatusBadRequest, "invalid host")
		return
	}

	proof := readChallengeProof(c)
	if reason := a.verifyBrowserChallenge(c, item, clientIP, redirectURI, proof); reason != "" {
		a.logChallengeAttempt(c, item, clientIP, http.StatusUnauthorized, reason, mode, domain.RequestDecisionChallenged)
		a.registerChallengeFailure(c, item, clientIP, reason)
		if err := a.renderBrowserChallenge(c, headers, host, redirectURI, mode, reason); err != nil {
			util.AbortError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	token := a.engine.IssueClearance(item.Name, clientIP)
	a.engine.ClearMitigationState(c.Request.Context(), item.Name, clientIP)
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
		ChallengeType:  "slider",
		OccurredAt:     time.Now(),
	})
	c.Redirect(http.StatusFound, redirectURI)
}

func requestHeaders(c *gin.Context) map[string]string {
	headers := map[string]string{}
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[strings.ToLower(key)] = values[0]
		}
	}
	return headers
}

func (a *API) renderBrowserChallenge(c *gin.Context, headers map[string]string, host string, redirectURI string, mode string, reason string) error {
	item, clientIP, err := a.resolveChallengeTarget(c, headers, host)
	if err != nil {
		return err
	}
	nonce, err := util.RandomToken(24)
	if err != nil {
		return err
	}
	session := browserChallengeSession{
		Domain:       item.Name,
		ClientIP:     clientIP,
		UserAgentSum: util.SHA256Hex(c.Request.UserAgent()),
		RedirectURI:  redirectURI,
		IssuedAtUnix: time.Now().UnixMilli(),
		Mode:         mode,
	}
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	if err := a.store.Redis().Set(c.Request.Context(), challengeSessionKey(a.cfg.Redis.KeyPrefix, nonce), payload, a.cfg.Security.ChallengeTTL).Err(); err != nil {
		return err
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusUnauthorized, protection.ChallengePage(c.GetHeader("Accept-Language"), item.Name, redirectURI, mode, reason, nonce))
	return nil
}

func (a *API) resolveChallengeTarget(c *gin.Context, headers map[string]string, host string) (domain.Domain, string, error) {
	item, err := a.engine.ResolveDomain(c.Request.Context(), strings.ToLower(strings.TrimSpace(host)))
	if err != nil {
		return domain.Domain{}, "", err
	}
	clientIP := c.GetHeader("X-Real-IP")
	if clientIP == "" {
		clientIP = c.ClientIP()
	}
	clientIP = a.engine.ResolveClientIP(item, protection.RequestContext{
		Headers:  headers,
		RemoteIP: clientIP,
	})
	return item, clientIP, nil
}

func (a *API) verifyBrowserChallenge(c *gin.Context, item domain.Domain, clientIP string, redirectURI string, proof browserChallengeProof) string {
	if proof.Honeypot != "" {
		return "bot_signal_detected"
	}
	if proof.Nonce == "" {
		return "challenge_expired"
	}

	key := challengeSessionKey(a.cfg.Redis.KeyPrefix, proof.Nonce)
	raw, err := a.store.Redis().Get(c.Request.Context(), key).Bytes()
	if err != nil {
		return "challenge_expired"
	}
	_ = a.store.Redis().Del(c.Request.Context(), key).Err()

	var session browserChallengeSession
	if err := json.Unmarshal(raw, &session); err != nil {
		return "challenge_expired"
	}
	if session.Domain != item.Name || session.ClientIP != clientIP || session.UserAgentSum != util.SHA256Hex(c.Request.UserAgent()) {
		return "challenge_mismatch"
	}
	if session.RedirectURI != "" && session.RedirectURI != redirectURI {
		return "challenge_mismatch"
	}
	if time.Now().UnixMilli()-session.IssuedAtUnix < minChallengeInteractionMS {
		return "challenge_too_fast"
	}
	if proof.WebDriver {
		return "webdriver_detected"
	}
	if !proof.SliderCompleted {
		return "slider_incomplete"
	}
	if proof.InteractionMS < minChallengeInteractionMS {
		return "challenge_too_fast"
	}
	if proof.PointerMoves < minChallengePointerMoves {
		return "interaction_too_low"
	}
	if proof.Language == "" || proof.Timezone == "" || proof.Screen == "" || proof.Viewport == "" || proof.Platform == "" {
		return "browser_signals_missing"
	}
	return ""
}

func (a *API) registerChallengeFailure(c *gin.Context, item domain.Domain, clientIP string, reason string) {
	key := a.cfg.Redis.KeyPrefix + "challenge_fail:" + item.Name + ":" + clientIP
	pipe := a.store.Redis().Pipeline()
	count := pipe.Incr(c.Request.Context(), key)
	pipe.Expire(c.Request.Context(), key, a.cfg.Security.AutoBanWindow)
	if _, err := pipe.Exec(c.Request.Context()); err != nil {
		return
	}
	if count.Val() < 4 {
		return
	}

	banKey := a.cfg.Redis.KeyPrefix + "ban:" + item.Name + ":" + clientIP
	if err := a.store.Redis().Set(c.Request.Context(), banKey, reason, a.cfg.Security.AutoBanTTL).Err(); err != nil {
		return
	}
	_ = a.store.RecordTemporaryBan(c.Request.Context(), &item.ID, clientIP, reason, "challenge_failure", time.Now().Add(a.cfg.Security.AutoBanTTL))
}

func (a *API) logChallengeAttempt(c *gin.Context, item domain.Domain, clientIP string, statusCode int, reason string, mode string, decision domain.RequestDecision) {
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
		Decision:       decision,
		StatusCode:     statusCode,
		BlockReason:    reason,
		ResponseTimeMS: 1,
		ChallengeType:  mode,
		OccurredAt:     time.Now(),
	})
}

func readChallengeProof(c *gin.Context) browserChallengeProof {
	return browserChallengeProof{
		Nonce:           strings.TrimSpace(c.PostForm("challenge_nonce")),
		Mode:            normalizeChallengeMode(c.PostForm("challenge_mode")),
		SliderCompleted: parseBoolForm(c.PostForm("slider_completed")),
		InteractionMS:   parseIntForm(c.PostForm("interaction_ms")),
		PointerMoves:    parseIntForm(c.PostForm("pointer_moves")),
		WebDriver:       parseBoolForm(c.PostForm("webdriver")),
		Language:        strings.TrimSpace(c.PostForm("language")),
		Timezone:        strings.TrimSpace(c.PostForm("timezone")),
		Screen:          strings.TrimSpace(c.PostForm("screen")),
		Viewport:        strings.TrimSpace(c.PostForm("viewport")),
		Platform:        strings.TrimSpace(c.PostForm("platform")),
		DeviceMemory:    strings.TrimSpace(c.PostForm("device_memory")),
		CPUCount:        strings.TrimSpace(c.PostForm("hardware_concurrency")),
		Honeypot:        strings.TrimSpace(c.PostForm("website")),
	}
}

func sanitizeRedirectURI(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") {
		return "/"
	}
	return value
}

func normalizeChallengeMode(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "js") {
		return "js"
	}
	return "cookie"
}

func parseBoolForm(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func parseIntForm(value string) int {
	parsed, _ := strconv.Atoi(strings.TrimSpace(value))
	return parsed
}

func challengeSessionKey(prefix string, nonce string) string {
	return prefix + "challenge:browser:" + nonce
}
