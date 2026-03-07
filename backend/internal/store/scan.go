package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"shieldpanel/backend/internal/domain"
)

func scanUser(row pgx.Row) (domain.User, error) {
	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Role,
		&user.Language,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func scanSessionWithUser(row pgx.Row) (domain.Session, domain.User, error) {
	var session domain.Session
	var user domain.User
	err := row.Scan(
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
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Role,
		&user.Language,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return session, user, err
}

func scanDomain(row pgx.Row) (domain.Domain, error) {
	var d domain.Domain
	err := row.Scan(
		&d.ID,
		&d.Name,
		&d.OriginHost,
		&d.OriginPort,
		&d.OriginProtocol,
		&d.OriginServerName,
		&d.Enabled,
		&d.ProtectionEnabled,
		&d.ProtectionMode,
		&d.ChallengeMode,
		&d.CloudflareMode,
		&d.SSLAutoIssue,
		&d.SSLEnabled,
		&d.ForceHTTPS,
		&d.RateLimitRPS,
		&d.RateLimitBurst,
		&d.BadBotMode,
		&d.HeaderValidation,
		&d.JSChallengeEnabled,
		&d.AllowedMethods,
		&d.Notes,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	return d, err
}

func scanDomainRule(row pgx.Row) (domain.DomainRule, error) {
	var rule domain.DomainRule
	err := row.Scan(
		&rule.ID,
		&rule.DomainID,
		&rule.Name,
		&rule.Type,
		&rule.Pattern,
		&rule.Action,
		&rule.Enabled,
		&rule.Priority,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	return rule, err
}

func scanIPEntry(row pgx.Row) (domain.IPEntry, error) {
	var entry domain.IPEntry
	err := row.Scan(
		&entry.ID,
		&entry.DomainID,
		&entry.ListType,
		&entry.IP,
		&entry.CIDR,
		&entry.Reason,
		&entry.ExpiresAt,
		&entry.CreatedBy,
		&entry.CreatedAt,
	)
	return entry, err
}

func scanTemporaryBan(row pgx.Row) (domain.TemporaryBan, error) {
	var ban domain.TemporaryBan
	err := row.Scan(
		&ban.ID,
		&ban.DomainID,
		&ban.IP,
		&ban.Reason,
		&ban.Source,
		&ban.ExpiresAt,
		&ban.CreatedAt,
	)
	return ban, err
}

func scanRequestLog(row pgx.Row) (domain.RequestLog, error) {
	var log domain.RequestLog
	err := row.Scan(
		&log.ID,
		&log.DomainID,
		&log.DomainName,
		&log.ClientIP,
		&log.CountryCode,
		&log.Method,
		&log.Path,
		&log.QueryString,
		&log.UserAgent,
		&log.RequestID,
		&log.Decision,
		&log.StatusCode,
		&log.BlockReason,
		&log.ResponseTimeMS,
		&log.Score,
		&log.ChallengeType,
		&log.CreatedAt,
	)
	return log, err
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return nil
	}
	return *id
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}
