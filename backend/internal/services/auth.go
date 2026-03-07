package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

type AuthService struct {
	store  *store.Store
	logger *slog.Logger
	cfg    config.Config
}

type LoginResult struct {
	User         domain.User
	SessionToken string
	CSRFToken    string
	ExpiresAt    time.Time
}

func NewAuthService(repo *store.Store, logger *slog.Logger, cfg config.Config) *AuthService {
	return &AuthService{
		store:  repo,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, identifier string, password string, rememberMe bool, ipAddress string, userAgent string) (LoginResult, error) {
	if err := s.checkLoginRateLimit(ctx, ipAddress); err != nil {
		return LoginResult{}, err
	}

	user, err := s.store.FindUserByIdentifier(ctx, identifier)
	if err != nil {
		if storeErr := s.bumpLoginRateLimit(ctx, ipAddress); storeErr != nil {
			s.logger.Warn("failed to increment login limiter", "error", storeErr)
		}
		return LoginResult{}, ErrInvalidCredentials
	}
	if err := util.CheckPassword(password, user.PasswordHash); err != nil {
		if storeErr := s.bumpLoginRateLimit(ctx, ipAddress); storeErr != nil {
			s.logger.Warn("failed to increment login limiter", "error", storeErr)
		}
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := util.RandomToken(32)
	if err != nil {
		return LoginResult{}, err
	}
	csrfToken, err := util.RandomToken(24)
	if err != nil {
		return LoginResult{}, err
	}
	expiresAt := time.Now().Add(s.cfg.Auth.SessionTTL)
	if rememberMe {
		expiresAt = time.Now().Add(s.cfg.Auth.RememberMeTTL)
	}

	if _, err := s.store.CreateSession(ctx, store.CreateSessionParams{
		UserID:     user.ID,
		TokenHash:  util.SHA256Hex(token),
		CSRFToken:  csrfToken,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		RememberMe: rememberMe,
		ExpiresAt:  expiresAt,
	}); err != nil {
		return LoginResult{}, err
	}

	if err := s.store.SetUserLastLogin(ctx, user.ID, time.Now()); err != nil {
		s.logger.Warn("failed to update last login", "error", err, "user", user.ID)
	}
	return LoginResult{
		User:         user,
		SessionToken: token,
		CSRFToken:    csrfToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *AuthService) AuthenticateSession(ctx context.Context, rawToken string) (domain.Session, domain.User, error) {
	if strings.TrimSpace(rawToken) == "" {
		return domain.Session{}, domain.User{}, ErrUnauthorized
	}
	session, user, err := s.store.GetSessionByTokenHash(ctx, util.SHA256Hex(rawToken))
	if err != nil {
		return domain.Session{}, domain.User{}, ErrUnauthorized
	}
	if err := s.store.TouchSession(ctx, session.ID); err != nil {
		s.logger.Warn("failed to touch session", "error", err, "session", session.ID)
	}
	return session, user, nil
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	if strings.TrimSpace(rawToken) == "" {
		return nil
	}
	return s.store.RevokeSessionByTokenHash(ctx, util.SHA256Hex(rawToken))
}

func (s *AuthService) checkLoginRateLimit(ctx context.Context, ipAddress string) error {
	key := s.cfg.Redis.KeyPrefix + "login:" + ipAddress
	value, err := s.store.Redis().Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}
	if value >= 10 {
		return ErrRateLimited
	}
	return nil
}

func (s *AuthService) bumpLoginRateLimit(ctx context.Context, ipAddress string) error {
	key := s.cfg.Redis.KeyPrefix + "login:" + ipAddress
	pipe := s.store.Redis().Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 15*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *AuthService) SetSessionCookiesSecure() bool {
	return s.cfg.Auth.CookieSecure
}

func (s *AuthService) SessionCookieName() string {
	return s.cfg.Auth.SessionCookieName
}

func (s *AuthService) CSRFCookieName() string {
	return s.cfg.Auth.CSRFCookieName
}

func (s *AuthService) CookieDomain() string {
	return s.cfg.Auth.CookieDomain
}

func (s *AuthService) SessionDuration(rememberMe bool) time.Duration {
	if rememberMe {
		return s.cfg.Auth.RememberMeTTL
	}
	return s.cfg.Auth.SessionTTL
}

func (s *AuthService) DebugString() string {
	return fmt.Sprintf("sessionCookie=%s csrfCookie=%s", s.cfg.Auth.SessionCookieName, s.cfg.Auth.CSRFCookieName)
}
