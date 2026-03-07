package services

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type IPControlService struct {
	store *store.Store
	audit *AuditService
}

func NewIPControlService(repo *store.Store, audit *AuditService) *IPControlService {
	return &IPControlService{
		store: repo,
		audit: audit,
	}
}

func (s *IPControlService) ListEntries(ctx context.Context, listType domain.ListType, domainID *uuid.UUID, search string, limit int, offset int) ([]domain.IPEntry, int64, error) {
	return s.store.ListIPEntries(ctx, listType, domainID, search, limit, offset)
}

func (s *IPControlService) CreateEntry(ctx context.Context, actor *domain.User, input store.IPEntryInput, ipAddress string, userAgent string) (domain.IPEntry, error) {
	item, err := s.store.UpsertIPEntry(ctx, input)
	if err != nil {
		return domain.IPEntry{}, err
	}
	details, _ := json.Marshal(item)
	s.audit.Record(ctx, actor, "change_rules", "ip_list", item.ID.String(), ipAddress, userAgent, string(details))
	return item, nil
}

func (s *IPControlService) DeleteEntry(ctx context.Context, actor *domain.User, id uuid.UUID, ipAddress string, userAgent string) error {
	if err := s.store.DeleteIPEntry(ctx, id); err != nil {
		return err
	}
	s.audit.Record(ctx, actor, "change_rules", "ip_list", id.String(), ipAddress, userAgent, "deleted")
	return nil
}

func (s *IPControlService) ListBans(ctx context.Context, limit int, offset int) ([]domain.TemporaryBan, int64, error) {
	return s.store.ListTemporaryBans(ctx, limit, offset)
}
