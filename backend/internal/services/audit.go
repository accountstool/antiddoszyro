package services

import (
	"context"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type AuditService struct {
	store *store.Store
}

func NewAuditService(repo *store.Store) *AuditService {
	return &AuditService{store: repo}
}

func (s *AuditService) Record(ctx context.Context, actor *domain.User, action string, entityType string, entityID string, ipAddress string, userAgent string, details string) {
	var userID *uuid.UUID
	username := "system"
	if actor != nil {
		userID = &actor.ID
		username = actor.Username
	}
	_ = s.store.CreateAuditLog(ctx, store.CreateAuditLogParams{
		UserID:     userID,
		Username:   username,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    details,
	})
}
