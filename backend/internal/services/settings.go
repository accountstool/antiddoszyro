package services

import (
	"context"
	"encoding/json"
	"sort"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type SettingsService struct {
	store *store.Store
	audit *AuditService
}

func NewSettingsService(repo *store.Store, audit *AuditService) *SettingsService {
	return &SettingsService{
		store: repo,
		audit: audit,
	}
}

func (s *SettingsService) List(ctx context.Context) ([]domain.SystemSetting, error) {
	settings, err := s.store.ListSettings(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(settings, func(i, j int) bool { return settings[i].Key < settings[j].Key })
	return settings, nil
}

func (s *SettingsService) Update(ctx context.Context, actor *domain.User, values map[string]string, ipAddress string, userAgent string) error {
	for key, value := range values {
		if err := s.store.UpsertSetting(ctx, key, value, "string"); err != nil {
			return err
		}
	}
	details, _ := json.Marshal(values)
	s.audit.Record(ctx, actor, "change_settings", "system", "settings", ipAddress, userAgent, string(details))
	return nil
}
