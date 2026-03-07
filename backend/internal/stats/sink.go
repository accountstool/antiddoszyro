package stats

import (
	"context"
	"log/slog"
	"time"

	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
)

type Sink struct {
	store     *store.Store
	logger    *slog.Logger
	batchSize int
	flushEvery time.Duration
	events    chan domain.RequestLogInput
}

func NewSink(repo *store.Store, logger *slog.Logger) *Sink {
	return &Sink{
		store:      repo,
		logger:     logger,
		batchSize:  100,
		flushEvery: 2 * time.Second,
		events:     make(chan domain.RequestLogInput, 5000),
	}
}

func (s *Sink) Start(ctx context.Context) {
	ticker := time.NewTicker(s.flushEvery)
	defer ticker.Stop()

	buffer := make([]domain.RequestLogInput, 0, s.batchSize)
	flush := func() {
		if len(buffer) == 0 {
			return
		}
		if err := s.store.InsertRequestLogs(context.Background(), buffer); err != nil {
			s.logger.Error("failed to flush request logs", "error", err, "count", len(buffer))
		}
		buffer = buffer[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return
		case event := <-s.events:
			buffer = append(buffer, event)
			if len(buffer) >= s.batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (s *Sink) Publish(event domain.RequestLogInput) {
	select {
	case s.events <- event:
	default:
		s.logger.Warn("dropping request log event because sink buffer is full")
	}
}
