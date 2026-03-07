package store

import (
	"context"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
)

type CreateAuditLogParams struct {
	UserID     *uuid.UUID
	Username   string
	Action     string
	EntityType string
	EntityID   string
	IPAddress  string
	UserAgent  string
	Details    string
}

func (s *Store) CreateAuditLog(ctx context.Context, params CreateAuditLogParams) error {
	_, err := s.db.Exec(ctx, `
		insert into audit_logs (user_id, username, action, entity_type, entity_id, ip_address, user_agent, details)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, nullableUUID(params.UserID), params.Username, params.Action, params.EntityType, params.EntityID, params.IPAddress, params.UserAgent, params.Details)
	return err
}

func (s *Store) ListAuditLogs(ctx context.Context, limit int, offset int) ([]domain.AuditLog, int64, error) {
	rows, err := s.db.Query(ctx, `
		select id, user_id, username, action, entity_type, entity_id, ip_address, user_agent, details, created_at
		from audit_logs
		order by created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []domain.AuditLog
	for rows.Next() {
		var item domain.AuditLog
		if err := rows.Scan(&item.ID, &item.UserID, &item.Username, &item.Action, &item.EntityType, &item.EntityID, &item.IPAddress, &item.UserAgent, &item.Details, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}
	var total int64
	if err := s.db.QueryRow(ctx, `select count(*) from audit_logs`).Scan(&total); err != nil {
		return nil, 0, err
	}
	return result, total, rows.Err()
}
