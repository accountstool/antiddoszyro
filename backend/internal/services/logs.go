package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type LogsService struct {
	store *store.Store
}

func NewLogsService(repo *store.Store) *LogsService {
	return &LogsService{store: repo}
}

func (s *LogsService) List(ctx context.Context, filters store.LogFilters, limit int, offset int) ([]domain.RequestLog, int64, error) {
	return s.store.ListRequestLogs(ctx, filters, limit, offset)
}

func (s *LogsService) Overview(ctx context.Context, filters store.LogFilters) (domain.StatsOverview, error) {
	return s.store.StatsOverview(ctx, filters)
}

func (s *LogsService) DomainOverview(ctx context.Context, domainID uuid.UUID, from, to time.Time) (domain.StatsOverview, error) {
	return s.store.DomainStats(ctx, domainID, from, to)
}
