package services

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type SSLService struct {
	store   *store.Store
	cfg     config.Config
	audit   *AuditService
	domains *DomainService
}

func NewSSLService(repo *store.Store, cfg config.Config, audit *AuditService, domains *DomainService) *SSLService {
	return &SSLService{
		store:   repo,
		cfg:     cfg,
		audit:   audit,
		domains: domains,
	}
}

func (s *SSLService) Issue(ctx context.Context, actor *domain.User, domainID uuid.UUID, ipAddress string, userAgent string) error {
	item, err := s.store.GetDomainByID(ctx, domainID)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "bash", s.cfg.Paths.IssueCertScript, item.Name, s.cfg.Nginx.ACMEWebroot)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("issue certificate: %w: %s", err, string(output))
	}
	now := time.Now()
	certPath := filepath.Join("/etc/letsencrypt/live", item.Name, "fullchain.pem")
	keyPath := filepath.Join("/etc/letsencrypt/live", item.Name, "privkey.pem")
	if err := s.store.UpsertSSLCertificate(ctx, item.ID, "Let's Encrypt", "issued", certPath, keyPath, nil, &now); err != nil {
		return err
	}
	if err := s.store.SetDomainSSLEnabled(ctx, item.ID, true); err != nil {
		return err
	}
	s.audit.Record(ctx, actor, "issue_ssl", "domain", item.ID.String(), ipAddress, userAgent, item.Name)
	return s.domains.Sync(ctx, actor, ipAddress, userAgent)
}

func (s *SSLService) Renew(ctx context.Context, actor *domain.User, domainID uuid.UUID, ipAddress string, userAgent string) error {
	item, err := s.store.GetDomainByID(ctx, domainID)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "bash", s.cfg.Paths.RenewCertScript, item.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("renew certificate: %w: %s", err, string(output))
	}
	now := time.Now()
	certPath := filepath.Join("/etc/letsencrypt/live", item.Name, "fullchain.pem")
	keyPath := filepath.Join("/etc/letsencrypt/live", item.Name, "privkey.pem")
	if err := s.store.UpsertSSLCertificate(ctx, item.ID, "Let's Encrypt", "renewed", certPath, keyPath, nil, &now); err != nil {
		return err
	}
	if err := s.store.SetDomainSSLEnabled(ctx, item.ID, true); err != nil {
		return err
	}
	s.audit.Record(ctx, actor, "renew_ssl", "domain", item.ID.String(), ipAddress, userAgent, item.Name)
	return s.domains.Sync(ctx, actor, ipAddress, userAgent)
}
