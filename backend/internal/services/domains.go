package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/nginx"
	"shieldpanel/backend/internal/store"
)

type DomainService struct {
	store   *store.Store
	nginx   *nginx.Manager
	audit   *AuditService
	logger  *slog.Logger
}

type DomainDetail struct {
	Domain      domain.Domain       `json:"domain"`
	Rules       []domain.DomainRule `json:"rules"`
	NginxStatus string              `json:"nginxStatus"`
}

func NewDomainService(repo *store.Store, nginxManager *nginx.Manager, audit *AuditService, logger *slog.Logger) *DomainService {
	return &DomainService{
		store:  repo,
		nginx:  nginxManager,
		audit:  audit,
		logger: logger,
	}
}

func (s *DomainService) List(ctx context.Context, limit int, offset int, search string) ([]domain.Domain, int64, error) {
	return s.store.ListDomains(ctx, limit, offset, search)
}

func (s *DomainService) Get(ctx context.Context, id uuid.UUID) (DomainDetail, error) {
	item, err := s.store.GetDomainByID(ctx, id)
	if err != nil {
		return DomainDetail{}, err
	}
	rules, err := s.store.ListDomainRules(ctx, id)
	if err != nil {
		return DomainDetail{}, err
	}
	return DomainDetail{
		Domain:      item,
		Rules:       rules,
		NginxStatus: "managed",
	}, nil
}

func (s *DomainService) Create(ctx context.Context, actor *domain.User, input store.DomainInput, rules []store.DomainRuleInput, ipAddress string, userAgent string) (domain.Domain, error) {
	item, err := s.store.CreateDomain(ctx, input)
	if err != nil {
		return domain.Domain{}, err
	}
	if err := s.store.ReplaceDomainRules(ctx, item.ID, rules); err != nil {
		return domain.Domain{}, err
	}
	if err := s.syncNginx(ctx); err != nil {
		return domain.Domain{}, err
	}
	details, _ := json.Marshal(map[string]any{"domain": item.Name, "origin": item.OriginHost, "port": item.OriginPort})
	s.audit.Record(ctx, actor, "add_domain", "domain", item.ID.String(), ipAddress, userAgent, string(details))
	return item, nil
}

func (s *DomainService) Update(ctx context.Context, actor *domain.User, id uuid.UUID, input store.DomainInput, rules []store.DomainRuleInput, ipAddress string, userAgent string) (domain.Domain, error) {
	item, err := s.store.UpdateDomain(ctx, id, input)
	if err != nil {
		return domain.Domain{}, err
	}
	if err := s.store.ReplaceDomainRules(ctx, id, rules); err != nil {
		return domain.Domain{}, err
	}
	if err := s.syncNginx(ctx); err != nil {
		return domain.Domain{}, err
	}
	details, _ := json.Marshal(map[string]any{"domain": item.Name, "origin": item.OriginHost, "port": item.OriginPort})
	s.audit.Record(ctx, actor, "edit_domain", "domain", item.ID.String(), ipAddress, userAgent, string(details))
	return item, nil
}

func (s *DomainService) Delete(ctx context.Context, actor *domain.User, id uuid.UUID, ipAddress string, userAgent string) error {
	item, err := s.store.GetDomainByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.store.DeleteDomain(ctx, id); err != nil {
		return err
	}
	if err := s.syncNginx(ctx); err != nil {
		return err
	}
	s.audit.Record(ctx, actor, "delete_domain", "domain", item.ID.String(), ipAddress, userAgent, item.Name)
	return nil
}

func (s *DomainService) Sync(ctx context.Context, actor *domain.User, ipAddress string, userAgent string) error {
	if err := s.syncNginx(ctx); err != nil {
		return err
	}
	s.audit.Record(ctx, actor, "reload_nginx", "system", "nginx", ipAddress, userAgent, "manual reload")
	return nil
}

func (s *DomainService) syncNginx(ctx context.Context) error {
	domainsList, _, err := s.store.ListDomains(ctx, 1000, 0, "")
	if err != nil {
		return err
	}
	return s.nginx.Sync(ctx, domainsList)
}
