package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
)

type CreateSessionParams struct {
	UserID     uuid.UUID
	TokenHash  string
	CSRFToken  string
	IPAddress  string
	UserAgent  string
	RememberMe bool
	ExpiresAt  time.Time
}

func (s *Store) CreateSession(ctx context.Context, params CreateSessionParams) (domain.Session, error) {
	var session domain.Session
	err := s.db.QueryRow(ctx, `
		insert into sessions (user_id, token_hash, csrf_token, ip_address, user_agent, remember_me, last_seen_at, expires_at)
		values ($1, $2, $3, $4, $5, $6, now(), $7)
		returning id, user_id, token_hash, csrf_token, ip_address, user_agent, remember_me, last_seen_at, expires_at, revoked_at, created_at
	`, params.UserID, params.TokenHash, params.CSRFToken, params.IPAddress, params.UserAgent, params.RememberMe, params.ExpiresAt).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.CSRFToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.RememberMe,
		&session.LastSeenAt,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)
	return session, err
}

func (s *Store) GetSessionByTokenHash(ctx context.Context, tokenHash string) (domain.Session, domain.User, error) {
	return scanSessionWithUser(s.db.QueryRow(ctx, `
		select
			s.id, s.user_id, s.token_hash, s.csrf_token, s.ip_address, s.user_agent, s.remember_me, s.last_seen_at, s.expires_at, s.revoked_at, s.created_at,
			u.id, u.username, u.email, u.password_hash, u.display_name, u.role, u.language, u.last_login_at, u.created_at, u.updated_at
		from sessions s
		join users u on u.id = s.user_id
		where s.token_hash = $1 and s.revoked_at is null and s.expires_at > now()
	`, tokenHash))
}

func (s *Store) TouchSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `update sessions set last_seen_at=now() where id=$1`, sessionID)
	return err
}

func (s *Store) RevokeSessionByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := s.db.Exec(ctx, `update sessions set revoked_at=now() where token_hash=$1 and revoked_at is null`, tokenHash)
	return err
}

func (s *Store) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `update sessions set revoked_at=now() where id=$1 and revoked_at is null`, sessionID)
	return err
}

func (s *Store) CleanupExpiredSessions(ctx context.Context) error {
	_, err := s.db.Exec(ctx, `delete from sessions where expires_at < now() or revoked_at is not null`)
	return err
}
