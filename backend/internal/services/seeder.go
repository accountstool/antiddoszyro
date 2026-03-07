package services

import (
	"context"
	"log/slog"
	"time"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

type Seeder struct {
	store  *store.Store
	cfg    config.Config
	logger *slog.Logger
}

type SeedResult struct {
	AdminIdentifier string
	SampleDomain    string
}

func NewSeeder(repo *store.Store, logger *slog.Logger, cfg config.Config) *Seeder {
	return &Seeder{
		store:  repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Seeder) Seed(ctx context.Context) (SeedResult, error) {
	if err := s.store.UpsertSetting(ctx, "language", s.cfg.Defaults.Language, "string"); err != nil {
		return SeedResult{}, err
	}
	if err := s.store.UpsertSetting(ctx, "panel_port", "8080", "number"); err != nil {
		return SeedResult{}, err
	}
	if err := s.store.UpsertSetting(ctx, "trusted_proxy_mode", "false", "bool"); err != nil {
		return SeedResult{}, err
	}

	adminIdentifier := "admin@shieldpanel.local"
	if _, err := s.store.FindUserByIdentifier(ctx, adminIdentifier); err != nil {
		hash, err := util.HashPassword("ChangeMe123!", s.cfg.Auth.BcryptCost)
		if err != nil {
			return SeedResult{}, err
		}
		if _, err := s.store.CreateUser(ctx, store.CreateUserParams{
			Username:     "admin",
			Email:        adminIdentifier,
			PasswordHash: hash,
			DisplayName:  "ShieldPanel Admin",
			Role:         domain.UserRoleOwner,
			Language:     s.cfg.Defaults.Language,
		}); err != nil {
			return SeedResult{}, err
		}
	}

	sampleDomain := "demo.local"
	if _, err := s.store.FindDomainByName(ctx, sampleDomain); err != nil {
		item, err := s.store.CreateDomain(ctx, store.DomainInput{
			Name:               sampleDomain,
			OriginHost:         "127.0.0.1",
			OriginPort:         8081,
			OriginProtocol:     "http",
			OriginServerName:   sampleDomain,
			Enabled:            true,
			ProtectionEnabled:  true,
			ProtectionMode:     domain.ProtectionBasic,
			ChallengeMode:      domain.ChallengeCookie,
			CloudflareMode:     false,
			SSLAutoIssue:       false,
			SSLEnabled:         false,
			ForceHTTPS:         false,
			RateLimitRPS:       s.cfg.Defaults.DefaultRateLimitRPS,
			RateLimitBurst:     s.cfg.Defaults.DefaultRateLimitBurst,
			BadBotMode:         true,
			HeaderValidation:   true,
			JSChallengeEnabled: false,
			AllowedMethods:     s.cfg.Security.AllowedMethods,
			Notes:              "Seeded sample domain",
		})
		if err != nil {
			return SeedResult{}, err
		}
		_ = s.store.ReplaceDomainRules(ctx, item.ID, []store.DomainRuleInput{
			{Name: "Block traversal", Type: "path", Pattern: "../", Action: "block", Enabled: true, Priority: 10},
			{Name: "Challenge exploit probes", Type: "query", Pattern: "union select", Action: "challenge", Enabled: true, Priority: 20},
		})
		now := time.Now()
		for i := 0; i < 32; i++ {
			_ = s.store.InsertRequestLogs(ctx, []domain.RequestLogInput{{
				DomainID:       &item.ID,
				DomainName:     item.Name,
				ClientIP:       "203.0.113.10",
				CountryCode:    "VN",
				Method:         "GET",
				Path:           "/",
				QueryString:    "",
				UserAgent:      "Mozilla/5.0",
				RequestID:      "seed-allowed",
				Decision:       domain.RequestDecisionAllowed,
				StatusCode:     200,
				ResponseTimeMS: 4,
				OccurredAt:     now.Add(time.Duration(-i) * time.Hour),
			}, {
				DomainID:       &item.ID,
				DomainName:     item.Name,
				ClientIP:       "198.51.100.9",
				CountryCode:    "US",
				Method:         "GET",
				Path:           "/../../etc/passwd",
				QueryString:    "",
				UserAgent:      "curl/8.0",
				RequestID:      "seed-blocked",
				Decision:       domain.RequestDecisionBlocked,
				StatusCode:     403,
				BlockReason:    "path_traversal",
				ResponseTimeMS: 2,
				OccurredAt:     now.Add(time.Duration(-i) * time.Hour),
			}})
		}
	}

	return SeedResult{
		AdminIdentifier: adminIdentifier,
		SampleDomain:    sampleDomain,
	}, nil
}
