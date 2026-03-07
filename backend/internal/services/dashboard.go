package services

import (
	"context"
	"time"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type DashboardService struct {
	store *store.Store
}

func NewDashboardService(repo *store.Store) *DashboardService {
	return &DashboardService{store: repo}
}

func (s *DashboardService) Summary(ctx context.Context, currentRPS int64, currentBlockedPS int64) (domain.DashboardSummary, error) {
	return s.store.DashboardSummary(ctx, currentRPS, currentBlockedPS)
}

func (s *DashboardService) Charts(ctx context.Context) (map[string]any, error) {
	last24h, err := s.store.DashboardTimeSeries(ctx, 24)
	if err != nil {
		return nil, err
	}
	last7d, err := s.store.DashboardTimeSeries(ctx, 24*7)
	if err != nil {
		return nil, err
	}
	overview, err := s.store.StatsOverview(ctx, store.LogFilters{
		From: time.Now().Add(-7 * 24 * time.Hour),
		To:   time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"last24h": last24h,
		"last7d":  last7d,
		"topIps":  overview.TopIPs,
		"topDomains": overview.TopDomains,
		"topReasons": overview.TopReasons,
	}, nil
}
