package protection

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/security"
	"shieldpanel/backend/internal/stats"
	"shieldpanel/backend/internal/store"
)

type cachedDomain struct {
	Domain    domain.Domain
	Rules     []domain.DomainRule
	ExpiresAt time.Time
}

type Engine struct {
	store   *store.Store
	logger  *slog.Logger
	cfg     config.Config
	signer  *security.ChallengeSigner
	sink    *stats.Sink
	cacheMu sync.RWMutex
	cache   map[string]cachedDomain
}

func NewEngine(repo *store.Store, logger *slog.Logger, cfg config.Config, sink *stats.Sink) *Engine {
	return &Engine{
		store:  repo,
		logger: logger,
		cfg:    cfg,
		signer: security.NewChallengeSigner(cfg.Security.ChallengeSecret),
		sink:   sink,
		cache:  map[string]cachedDomain{},
	}
}

func (e *Engine) Evaluate(ctx context.Context, input RequestContext) (Decision, error) {
	started := time.Now()
	item, rules, err := e.loadDomain(ctx, input.Host)
	if err != nil {
		return Decision{
			Action:         ActionBlock,
			Reason:         "unknown_domain",
			StatusCode:     http.StatusForbidden,
			ResponseTimeMS: int(time.Since(started).Milliseconds()),
			OccurredAt:     time.Now(),
		}, nil
	}

	clientIP := e.resolveClientIP(item, input)
	countryCode := input.Headers["cf-ipcountry"]
	decision := Decision{
		Action:         ActionAllow,
		Reason:         "ok",
		StatusCode:     http.StatusNoContent,
		Domain:         item,
		ClientIP:       clientIP,
		CountryCode:    countryCode,
		OccurredAt:     time.Now(),
	}
	defer func() {
		decision.ResponseTimeMS = int(time.Since(started).Milliseconds())
		e.publish(input, decision)
	}()

	if !item.Enabled {
		decision.Action = ActionBlock
		decision.Reason = "domain_disabled"
		decision.StatusCode = http.StatusForbidden
		return decision, nil
	}

	if entry, err := e.store.FindMatchingIPEntry(ctx, &item.ID, domain.ListTypeWhitelist, clientIP); err == nil && entry.ID != uuid.Nil {
		decision.Reason = "whitelist"
		return decision, nil
	}
	if entry, err := e.store.FindMatchingIPEntry(ctx, &item.ID, domain.ListTypeBlacklist, clientIP); err == nil && entry.ID != uuid.Nil {
		decision.Action = ActionBlock
		decision.Reason = "blacklist"
		decision.StatusCode = http.StatusForbidden
		return decision, nil
	}

	if banned, reason := e.isTemporarilyBanned(ctx, item.Name, clientIP); banned {
		decision.Action = ActionBlock
		decision.Reason = reason
		decision.StatusCode = http.StatusForbidden
		return decision, nil
	}

	if !methodAllowed(item.AllowedMethods, input.Method) {
		decision.Action = ActionBlock
		decision.Reason = "invalid_method"
		decision.StatusCode = http.StatusForbidden
		e.registerOffense(ctx, &item.ID, item.Name, clientIP, decision.Reason)
		return decision, nil
	}

	if action, reason := matchCustomRules(rules, input); action != "" {
		switch action {
		case "block":
			decision.Action = ActionBlock
			decision.Reason = reason
			decision.StatusCode = http.StatusForbidden
			e.registerOffense(ctx, &item.ID, item.Name, clientIP, reason)
		case "challenge":
			decision.Action = ActionChallenge
			decision.Reason = reason
			decision.StatusCode = http.StatusUnauthorized
			decision.ChallengeMode = e.resolveChallengeMode(item)
		}
		return decision, nil
	}

	allowed, err := allowByTokenBucket(ctx, e.store.Redis(), e.cfg, item.Name, clientIP, item.RateLimitRPS, item.RateLimitBurst)
	if err != nil {
		e.logger.Warn("rate limiter failed", "error", err, "domain", item.Name)
	}
	if !allowed {
		e.registerOffense(ctx, &item.ID, item.Name, clientIP, "rate_limit")
		if item.ProtectionMode == domain.ProtectionUnderAttack || item.ProtectionMode == domain.ProtectionAggressive {
			decision.Action = ActionBlock
			decision.Reason = "rate_limit"
			decision.StatusCode = http.StatusForbidden
			return decision, nil
		}
		decision.Action = ActionChallenge
		decision.Reason = "rate_limit"
		decision.StatusCode = http.StatusUnauthorized
		decision.ChallengeMode = e.resolveChallengeMode(item)
		return decision, nil
	}

	if item.BadBotMode && detectBadBot(input.UserAgent) {
		decision.Score += 4
		decision.Reason = "bad_bot_signature"
	}
	if item.HeaderValidation {
		decision.Score += scoreHeaders(input.Headers, input.UserAgent)
	}
	if suspicious, reason := detectSuspiciousPath(input.Path, input.Query); suspicious {
		decision.Score += 5
		decision.Reason = reason
	}

	if item.ProtectionMode == domain.ProtectionUnderAttack && !e.hasValidClearance(item, clientIP, input.Cookies) {
		decision.Action = ActionChallenge
		decision.Reason = "under_attack"
		decision.StatusCode = http.StatusUnauthorized
		decision.ChallengeMode = e.resolveChallengeMode(item)
		return decision, nil
	}

	if decision.Score >= 7 {
		decision.Action = ActionBlock
		decision.StatusCode = http.StatusForbidden
		if decision.Reason == "ok" {
			decision.Reason = "heuristic_block"
		}
		e.registerOffense(ctx, &item.ID, item.Name, clientIP, decision.Reason)
		return decision, nil
	}

	if decision.Score >= 4 || (item.ChallengeMode != domain.ChallengeOff && !e.hasValidClearance(item, clientIP, input.Cookies)) {
		decision.Action = ActionChallenge
		decision.StatusCode = http.StatusUnauthorized
		decision.ChallengeMode = e.resolveChallengeMode(item)
		if decision.Reason == "ok" {
			decision.Reason = "challenge_required"
		}
		e.registerOffense(ctx, &item.ID, item.Name, clientIP, decision.Reason)
		return decision, nil
	}

	return decision, nil
}

func (e *Engine) IssueClearance(domainName string, clientIP string) string {
	return e.signer.Sign(domainName, clientIP, time.Now().Add(e.cfg.Security.ChallengeGraceTTL))
}

func (e *Engine) VerifyClearance(domainName string, clientIP string, token string) bool {
	return e.signer.Verify(token, domainName, clientIP, time.Now())
}

func (e *Engine) ResolveDomain(ctx context.Context, host string) (domain.Domain, error) {
	item, _, err := e.loadDomain(ctx, host)
	return item, err
}

func (e *Engine) loadDomain(ctx context.Context, host string) (domain.Domain, []domain.DomainRule, error) {
	host = strings.ToLower(strings.TrimSpace(host))
	e.cacheMu.RLock()
	if cached, ok := e.cache[host]; ok && cached.ExpiresAt.After(time.Now()) {
		e.cacheMu.RUnlock()
		return cached.Domain, cached.Rules, nil
	}
	e.cacheMu.RUnlock()

	item, err := e.store.FindDomainByName(ctx, host)
	if err != nil {
		return domain.Domain{}, nil, err
	}
	rules, err := e.store.ListDomainRules(ctx, item.ID)
	if err != nil {
		return domain.Domain{}, nil, err
	}
	e.cacheMu.Lock()
	e.cache[host] = cachedDomain{Domain: item, Rules: rules, ExpiresAt: time.Now().Add(30 * time.Second)}
	e.cacheMu.Unlock()
	return item, rules, nil
}

func (e *Engine) resolveClientIP(item domain.Domain, input RequestContext) string {
	if item.CloudflareMode {
		if value := strings.TrimSpace(input.Headers["cf-connecting-ip"]); value != "" {
			return value
		}
	}
	if e.cfg.Security.TrustedProxyMode {
		for _, header := range e.cfg.Security.TrustedProxyHeaders {
			if value := strings.TrimSpace(input.Headers[strings.ToLower(header)]); value != "" {
				if strings.Contains(value, ",") {
					return strings.TrimSpace(strings.Split(value, ",")[0])
				}
				return value
			}
		}
	}
	return input.RemoteIP
}

func (e *Engine) hasValidClearance(item domain.Domain, clientIP string, cookies map[string]string) bool {
	token := cookies["shieldpanel_clearance"]
	if token == "" {
		return false
	}
	return e.VerifyClearance(item.Name, clientIP, token)
}

func (e *Engine) resolveChallengeMode(item domain.Domain) string {
	if item.JSChallengeEnabled || item.ChallengeMode == domain.ChallengeJS {
		return "js"
	}
	return "cookie"
}

func (e *Engine) isTemporarilyBanned(ctx context.Context, domainName string, clientIP string) (bool, string) {
	key := e.cfg.Redis.KeyPrefix + "ban:" + domainName + ":" + clientIP
	reason, err := e.store.Redis().Get(ctx, key).Result()
	if err != nil {
		return false, ""
	}
	return true, reason
}

func (e *Engine) registerOffense(ctx context.Context, domainID *uuid.UUID, domainName string, clientIP string, reason string) {
	key := e.cfg.Redis.KeyPrefix + "offense:" + domainName + ":" + clientIP
	pipe := e.store.Redis().Pipeline()
	count := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, e.cfg.Security.AutoBanWindow)
	_, err := pipe.Exec(ctx)
	if err != nil {
		e.logger.Warn("failed to record offense", "error", err)
		return
	}
	if count.Val() >= int64(e.cfg.Security.AutoBanThreshold) {
		banKey := e.cfg.Redis.KeyPrefix + "ban:" + domainName + ":" + clientIP
		if err := e.store.Redis().Set(ctx, banKey, reason, e.cfg.Security.AutoBanTTL).Err(); err != nil {
			e.logger.Warn("failed to set temporary ban", "error", err)
			return
		}
		if err := e.store.RecordTemporaryBan(ctx, domainID, clientIP, reason, "auto_ban", time.Now().Add(e.cfg.Security.AutoBanTTL)); err != nil {
			e.logger.Warn("failed to persist temporary ban", "error", err)
		}
	}
}

func (e *Engine) publish(input RequestContext, decision Decision) {
	if decision.Domain.Name == "" {
		return
	}
	logDecision := domain.RequestDecisionAllowed
	switch decision.Action {
	case ActionBlock:
		logDecision = domain.RequestDecisionBlocked
	case ActionChallenge:
		logDecision = domain.RequestDecisionChallenged
	}
	e.sink.Publish(domain.RequestLogInput{
		DomainID:       &decision.Domain.ID,
		DomainName:     decision.Domain.Name,
		ClientIP:       decision.ClientIP,
		CountryCode:    decision.CountryCode,
		Method:         input.Method,
		Path:           input.Path,
		QueryString:    input.Query,
		UserAgent:      input.UserAgent,
		RequestID:      input.RequestID,
		Decision:       logDecision,
		StatusCode:     decision.StatusCode,
		BlockReason:    decision.Reason,
		ResponseTimeMS: decision.ResponseTimeMS,
		Score:          decision.Score,
		ChallengeType:  decision.ChallengeMode,
		OccurredAt:     decision.OccurredAt,
	})
	e.bumpRealtimeCounter(logDecision)
}

func (e *Engine) bumpRealtimeCounter(decision domain.RequestDecision) {
	now := time.Now().Unix()
	reqKey := e.cfg.Redis.KeyPrefix + "realtime:req:" + strconv.FormatInt(now, 10)
	pipe := e.store.Redis().Pipeline()
	pipe.Incr(context.Background(), reqKey)
	pipe.Expire(context.Background(), reqKey, 5*time.Second)
	if decision == domain.RequestDecisionBlocked {
		blockKey := e.cfg.Redis.KeyPrefix + "realtime:block:" + strconv.FormatInt(now, 10)
		pipe.Incr(context.Background(), blockKey)
		pipe.Expire(context.Background(), blockKey, 5*time.Second)
	}
	_, _ = pipe.Exec(context.Background())
}

func (e *Engine) CurrentRates(ctx context.Context) (int64, int64) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	reqKey := e.cfg.Redis.KeyPrefix + "realtime:req:" + now
	blockKey := e.cfg.Redis.KeyPrefix + "realtime:block:" + now
	req, _ := e.store.Redis().Get(ctx, reqKey).Int64()
	blocked, _ := e.store.Redis().Get(ctx, blockKey).Int64()
	return req, blocked
}

func methodAllowed(allowed []string, method string) bool {
	for _, candidate := range allowed {
		if strings.EqualFold(candidate, method) {
			return true
		}
	}
	return false
}
