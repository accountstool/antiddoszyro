package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
)

type IPEntryInput struct {
	DomainID  *uuid.UUID
	ListType  domain.ListType
	IP        string
	CIDR      string
	Reason    string
	ExpiresAt *time.Time
	CreatedBy *uuid.UUID
}

func (s *Store) ListIPEntries(ctx context.Context, listType domain.ListType, domainID *uuid.UUID, search string, limit int, offset int) ([]domain.IPEntry, int64, error) {
	parts := []string{"list_type=$1"}
	args := []any{listType}
	if domainID != nil {
		parts = append(parts, fmt.Sprintf("domain_id=$%d", len(args)+1))
		args = append(args, *domainID)
	}
	if search = strings.TrimSpace(search); search != "" {
		parts = append(parts, fmt.Sprintf("(ip ilike $%d or cidr ilike $%d or reason ilike $%d)", len(args)+1, len(args)+1, len(args)+1))
		args = append(args, "%"+search+"%")
	}
	query := `
		select id, domain_id, list_type, ip, cidr, reason, expires_at, created_by, created_at
		from ip_lists
		where ` + strings.Join(parts, " and ") + `
		order by created_at desc
	`
	query += fmt.Sprintf(" limit $%d offset $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []domain.IPEntry
	for rows.Next() {
		item, err := scanIPEntry(rows)
		if err != nil {
			return nil, 0, err
		}
		entries = append(entries, item)
	}

	countQuery := `select count(*) from ip_lists where ` + strings.Join(parts, " and ")
	var total int64
	if err := s.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return entries, total, rows.Err()
}

func (s *Store) UpsertIPEntry(ctx context.Context, input IPEntryInput) (domain.IPEntry, error) {
	var cidr any
	if strings.TrimSpace(input.CIDR) != "" {
		cidr = input.CIDR
	}
	return scanIPEntry(s.db.QueryRow(ctx, `
		insert into ip_lists (domain_id, list_type, ip, cidr, reason, expires_at, created_by)
		values ($1,$2,$3,$4,$5,$6,$7)
		returning id, domain_id, list_type, ip, cidr, reason, expires_at, created_by, created_at
	`, nullableUUID(input.DomainID), input.ListType, input.IP, cidr, input.Reason, nullableTime(input.ExpiresAt), nullableUUID(input.CreatedBy)))
}

func (s *Store) DeleteIPEntry(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `delete from ip_lists where id=$1`, id)
	return err
}

func (s *Store) FindMatchingIPEntry(ctx context.Context, domainID *uuid.UUID, listType domain.ListType, ip string) (domain.IPEntry, error) {
	if domainID != nil {
		return scanIPEntry(s.db.QueryRow(ctx, `
			select id, domain_id, list_type, ip, cidr, reason, expires_at, created_by, created_at
			from ip_lists
			where list_type=$1
			  and (domain_id is null or domain_id=$2)
			  and (ip=$3 or (cidr is not null and $3::inet <<= cidr::cidr))
			  and (expires_at is null or expires_at > now())
			order by domain_id nulls last, created_at desc
			limit 1
		`, listType, *domainID, ip))
	}
	return scanIPEntry(s.db.QueryRow(ctx, `
		select id, domain_id, list_type, ip, cidr, reason, expires_at, created_by, created_at
		from ip_lists
		where list_type=$1
		  and domain_id is null
		  and (ip=$2 or (cidr is not null and $2::inet <<= cidr::cidr))
		  and (expires_at is null or expires_at > now())
		order by created_at desc
		limit 1
	`, listType, ip))
}

func (s *Store) RecordTemporaryBan(ctx context.Context, domainID *uuid.UUID, ip string, reason string, source string, expiresAt time.Time) error {
	_, err := s.db.Exec(ctx, `
		insert into temporary_bans (domain_id, ip, reason, source, expires_at)
		values ($1,$2,$3,$4,$5)
	`, nullableUUID(domainID), ip, reason, source, expiresAt)
	return err
}

func (s *Store) ListTemporaryBans(ctx context.Context, limit int, offset int) ([]domain.TemporaryBan, int64, error) {
	rows, err := s.db.Query(ctx, `
		select id, domain_id, ip, reason, source, expires_at, created_at
		from temporary_bans
		where expires_at > now()
		order by expires_at asc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []domain.TemporaryBan
	for rows.Next() {
		item, err := scanTemporaryBan(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}
	var total int64
	if err := s.db.QueryRow(ctx, `select count(*) from temporary_bans where expires_at > now()`).Scan(&total); err != nil {
		return nil, 0, err
	}
	return result, total, rows.Err()
}
