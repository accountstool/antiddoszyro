package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
)

type DomainInput struct {
	Name               string
	OriginHost         string
	OriginPort         int
	OriginProtocol     string
	OriginServerName   string
	Enabled            bool
	ProtectionEnabled  bool
	ProtectionMode     domain.ProtectionMode
	ChallengeMode      domain.ChallengeMode
	CloudflareMode     bool
	SSLAutoIssue       bool
	SSLEnabled         bool
	ForceHTTPS         bool
	RateLimitRPS       int
	RateLimitBurst     int
	BadBotMode         bool
	HeaderValidation   bool
	JSChallengeEnabled bool
	AllowedMethods     []string
	Notes              string
}

func (s *Store) ListDomains(ctx context.Context, limit int, offset int, search string) ([]domain.Domain, int64, error) {
	search = strings.TrimSpace(search)
	query := `
		select id, name, origin_host, origin_port, origin_protocol, origin_server_name, enabled, protection_enabled, protection_mode, challenge_mode, cloudflare_mode,
		       ssl_auto_issue, ssl_enabled, force_https, rate_limit_rps, rate_limit_burst, bad_bot_mode, header_validation, js_challenge_enabled, allowed_methods,
		       notes, created_at, updated_at
		from domains
	`
	args := []any{}
	if search != "" {
		query += ` where name ilike $1 or origin_host ilike $1`
		args = append(args, "%"+search+"%")
	}
	query += ` order by created_at desc`
	query += fmt.Sprintf(" limit $%d offset $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var domainsList []domain.Domain
	for rows.Next() {
		item, err := scanDomain(rows)
		if err != nil {
			return nil, 0, err
		}
		domainsList = append(domainsList, item)
	}

	countQuery := `select count(*) from domains`
	countArgs := []any{}
	if search != "" {
		countQuery += ` where name ilike $1 or origin_host ilike $1`
		countArgs = append(countArgs, "%"+search+"%")
	}
	var total int64
	if err := s.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return domainsList, total, rows.Err()
}

func (s *Store) CreateDomain(ctx context.Context, input DomainInput) (domain.Domain, error) {
	return scanDomain(s.db.QueryRow(ctx, `
		insert into domains (
			name, origin_host, origin_port, origin_protocol, origin_server_name, enabled, protection_enabled, protection_mode, challenge_mode, cloudflare_mode,
			ssl_auto_issue, ssl_enabled, force_https, rate_limit_rps, rate_limit_burst, bad_bot_mode, header_validation, js_challenge_enabled, allowed_methods, notes
		)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		returning id, name, origin_host, origin_port, origin_protocol, origin_server_name, enabled, protection_enabled, protection_mode, challenge_mode, cloudflare_mode,
		          ssl_auto_issue, ssl_enabled, force_https, rate_limit_rps, rate_limit_burst, bad_bot_mode, header_validation, js_challenge_enabled, allowed_methods, notes,
		          created_at, updated_at
	`, strings.ToLower(input.Name), input.OriginHost, input.OriginPort, input.OriginProtocol, input.OriginServerName, input.Enabled, input.ProtectionEnabled, input.ProtectionMode, input.ChallengeMode, input.CloudflareMode,
		input.SSLAutoIssue, input.SSLEnabled, input.ForceHTTPS, input.RateLimitRPS, input.RateLimitBurst, input.BadBotMode, input.HeaderValidation, input.JSChallengeEnabled, input.AllowedMethods, input.Notes))
}

func (s *Store) UpdateDomain(ctx context.Context, id uuid.UUID, input DomainInput) (domain.Domain, error) {
	return scanDomain(s.db.QueryRow(ctx, `
		update domains
		set name=$2,
			origin_host=$3,
			origin_port=$4,
			origin_protocol=$5,
			origin_server_name=$6,
			enabled=$7,
			protection_enabled=$8,
			protection_mode=$9,
			challenge_mode=$10,
			cloudflare_mode=$11,
			ssl_auto_issue=$12,
			ssl_enabled=$13,
			force_https=$14,
			rate_limit_rps=$15,
			rate_limit_burst=$16,
			bad_bot_mode=$17,
			header_validation=$18,
			js_challenge_enabled=$19,
			allowed_methods=$20,
			notes=$21,
			updated_at=now()
		where id=$1
		returning id, name, origin_host, origin_port, origin_protocol, origin_server_name, enabled, protection_enabled, protection_mode, challenge_mode, cloudflare_mode,
		          ssl_auto_issue, ssl_enabled, force_https, rate_limit_rps, rate_limit_burst, bad_bot_mode, header_validation, js_challenge_enabled, allowed_methods, notes,
		          created_at, updated_at
	`, id, strings.ToLower(input.Name), input.OriginHost, input.OriginPort, input.OriginProtocol, input.OriginServerName, input.Enabled, input.ProtectionEnabled, input.ProtectionMode, input.ChallengeMode, input.CloudflareMode,
		input.SSLAutoIssue, input.SSLEnabled, input.ForceHTTPS, input.RateLimitRPS, input.RateLimitBurst, input.BadBotMode, input.HeaderValidation, input.JSChallengeEnabled, input.AllowedMethods, input.Notes))
}

func (s *Store) GetDomainByID(ctx context.Context, id uuid.UUID) (domain.Domain, error) {
	return scanDomain(s.db.QueryRow(ctx, `
		select id, name, origin_host, origin_port, origin_protocol, origin_server_name, enabled, protection_enabled, protection_mode, challenge_mode, cloudflare_mode,
		       ssl_auto_issue, ssl_enabled, force_https, rate_limit_rps, rate_limit_burst, bad_bot_mode, header_validation, js_challenge_enabled, allowed_methods,
		       notes, created_at, updated_at
		from domains
		where id=$1
	`, id))
}

func (s *Store) FindDomainByName(ctx context.Context, name string) (domain.Domain, error) {
	return scanDomain(s.db.QueryRow(ctx, `
		select id, name, origin_host, origin_port, origin_protocol, origin_server_name, enabled, protection_enabled, protection_mode, challenge_mode, cloudflare_mode,
		       ssl_auto_issue, ssl_enabled, force_https, rate_limit_rps, rate_limit_burst, bad_bot_mode, header_validation, js_challenge_enabled, allowed_methods,
		       notes, created_at, updated_at
		from domains
		where name=$1
	`, strings.ToLower(strings.TrimSpace(name))))
}

func (s *Store) DeleteDomain(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `delete from domains where id=$1`, id)
	return err
}

func (s *Store) SetDomainSSLEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	_, err := s.db.Exec(ctx, `
		update domains
		set ssl_enabled=$2,
			updated_at=now()
		where id=$1
	`, id, enabled)
	return err
}

func (s *Store) ListDomainRules(ctx context.Context, domainID uuid.UUID) ([]domain.DomainRule, error) {
	rows, err := s.db.Query(ctx, `
		select id, domain_id, name, type, pattern, action, enabled, priority, created_at, updated_at
		from domain_rules
		where domain_id=$1
		order by priority asc, created_at asc
	`, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []domain.DomainRule
	for rows.Next() {
		rule, err := scanDomainRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

type DomainRuleInput struct {
	Name     string
	Type     string
	Pattern  string
	Action   string
	Enabled  bool
	Priority int
}

func (s *Store) ReplaceDomainRules(ctx context.Context, domainID uuid.UUID, rules []DomainRuleInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `delete from domain_rules where domain_id=$1`, domainID); err != nil {
		return err
	}
	for _, rule := range rules {
		if _, err := tx.Exec(ctx, `
			insert into domain_rules (domain_id, name, type, pattern, action, enabled, priority)
			values ($1,$2,$3,$4,$5,$6,$7)
		`, domainID, rule.Name, rule.Type, rule.Pattern, rule.Action, rule.Enabled, rule.Priority); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
