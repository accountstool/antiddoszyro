package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"shieldpanel/backend/internal/domain"
)

type LogFilters struct {
	DomainID *uuid.UUID
	From     time.Time
	To       time.Time
	Decision string
	Reason   string
}

func (s *Store) InsertRequestLogs(ctx context.Context, inputs []domain.RequestLogInput) error {
	if len(inputs) == 0 {
		return nil
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range inputs {
		if _, err := tx.Exec(ctx, `
			insert into request_logs (
				domain_id, domain_name, client_ip, country_code, method, path, query_string, user_agent,
				request_id, decision, status_code, block_reason, response_time_ms, score, challenge_type, created_at
			)
			values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		`, nullableUUID(item.DomainID), item.DomainName, item.ClientIP, item.CountryCode, item.Method, item.Path, item.QueryString, item.UserAgent,
			item.RequestID, item.Decision, item.StatusCode, item.BlockReason, item.ResponseTimeMS, item.Score, item.ChallengeType, item.OccurredAt); err != nil {
			return err
		}

	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	for _, item := range inputs {
		if item.DomainID == nil {
			continue
		}
		if _, err := s.db.Exec(ctx, `
			insert into traffic_rollups (
				window_start, granularity, domain_id, domain_name, decision, block_reason, request_count
			)
			values (date_trunc('hour', $1::timestamptz), 'hour', $2, $3, $4, $5, 1)
			on conflict (window_start, granularity, domain_id, decision, block_reason)
			do update set request_count = traffic_rollups.request_count + 1
		`, item.OccurredAt, *item.DomainID, item.DomainName, item.Decision, item.BlockReason); err != nil {
			s.logger.Warn("failed to update traffic rollup", "error", err, "domain", item.DomainName, "decision", item.Decision)
		}
	}

	return nil
}

func (s *Store) ListRequestLogs(ctx context.Context, filters LogFilters, limit int, offset int) ([]domain.RequestLog, int64, error) {
	where, args := buildLogWhere(filters)
	query := `
		select id, domain_id, domain_name, client_ip, country_code, method, path, query_string, user_agent, request_id, decision, status_code, block_reason, response_time_ms, score, challenge_type, created_at
		from request_logs
		` + where + `
		order by created_at desc
	`
	query += fmt.Sprintf(" limit $%d offset $%d", len(args)+1, len(args)+2)
	argsWithLimit := append(append([]any{}, args...), limit, offset)

	rows, err := s.db.Query(ctx, query, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []domain.RequestLog
	for rows.Next() {
		item, err := scanRequestLog(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}

	var total int64
	if err := s.db.QueryRow(ctx, `select count(*) from request_logs `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return result, total, rows.Err()
}

func buildLogWhere(filters LogFilters) (string, []any) {
	parts := []string{"created_at between $1 and $2"}
	args := []any{filters.From, filters.To}
	if filters.DomainID != nil {
		parts = append(parts, fmt.Sprintf("domain_id=$%d", len(args)+1))
		args = append(args, *filters.DomainID)
	}
	if filters.Decision != "" {
		parts = append(parts, fmt.Sprintf("decision=$%d", len(args)+1))
		args = append(args, filters.Decision)
	}
	if filters.Reason != "" {
		parts = append(parts, fmt.Sprintf("block_reason=$%d", len(args)+1))
		args = append(args, filters.Reason)
	}
	return `where ` + strings.Join(parts, " and "), args
}

func (s *Store) DashboardSummary(ctx context.Context, currentRPS int64, currentBlockedPS int64) (domain.DashboardSummary, error) {
	summary := domain.DashboardSummary{
		Healthy:          true,
		CurrentRPS:       currentRPS,
		CurrentBlockedPS: currentBlockedPS,
	}
	if err := s.db.QueryRow(ctx, `select count(*) from domains where enabled=true`).Scan(&summary.TotalDomains); err != nil {
		return summary, err
	}
	if err := s.db.QueryRow(ctx, `
		select count(*)
		from request_logs
		where decision='blocked' and created_at >= date_trunc('day', now())
	`).Scan(&summary.BlockedToday); err != nil {
		return summary, err
	}
	if err := s.db.QueryRow(ctx, `
		select count(*),
		       count(*) filter (where decision='allowed'),
		       count(*) filter (where decision='blocked'),
		       count(*) filter (where decision='challenged'),
		       coalesce(
		        (count(*) filter (where decision='challenge_passed'))::double precision /
		        nullif((count(*) filter (where decision='challenged'))::double precision, 0),
		        0
		       )
		from request_logs
		where created_at >= now() - interval '24 hours'
	`).Scan(&summary.TotalRequests24h, &summary.Allowed24h, &summary.Blocked24h, &summary.Challenged24h, &summary.ChallengePassRate); err != nil {
		return summary, err
	}

	err := s.db.QueryRow(ctx, `
		select domain_name
		from request_logs
		where created_at >= now() - interval '24 hours'
		group by domain_name
		order by count(*) filter (where decision='blocked') desc, count(*) desc
		limit 1
	`).Scan(&summary.TopAttackedDomain)
	if err != nil && err != pgx.ErrNoRows {
		return summary, err
	}

	err = s.db.QueryRow(ctx, `
		select client_ip
		from request_logs
		where created_at >= now() - interval '24 hours' and decision='blocked'
		group by client_ip
		order by count(*) desc
		limit 1
	`).Scan(&summary.TopAttackingIP)
	if err != nil && err != pgx.ErrNoRows {
		return summary, err
	}
	return summary, nil
}

func (s *Store) DashboardTimeSeries(ctx context.Context, hours int) ([]domain.TimePoint, error) {
	rows, err := s.db.Query(ctx, `
		select to_char(bucket, 'YYYY-MM-DD HH24:00') as label,
		       coalesce(sum(case when decision='allowed' then request_count else 0 end), 0) as allowed,
		       coalesce(sum(case when decision='blocked' then request_count else 0 end), 0) as blocked,
		       coalesce(sum(case when decision='challenged' then request_count else 0 end), 0) as challenged
		from (
			select date_trunc('hour', created_at) as bucket, decision, count(*) as request_count
			from request_logs
			where created_at >= now() - make_interval(hours => $1)
			group by bucket, decision
		) t
		group by bucket
		order by bucket asc
	`, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []domain.TimePoint
	for rows.Next() {
		var item domain.TimePoint
		if err := rows.Scan(&item.Label, &item.Allowed, &item.Blocked, &item.Challenge); err != nil {
			return nil, err
		}
		points = append(points, item)
	}
	return points, rows.Err()
}

func (s *Store) DomainStats(ctx context.Context, domainID uuid.UUID, from time.Time, to time.Time) (domain.StatsOverview, error) {
	return s.StatsOverview(ctx, LogFilters{
		DomainID: &domainID,
		From:     from,
		To:       to,
	})
}

func (s *Store) StatsOverview(ctx context.Context, filters LogFilters) (domain.StatsOverview, error) {
	overview := domain.StatsOverview{}
	where, args := buildLogWhere(filters)

	if err := s.db.QueryRow(ctx, `
		select
			count(*),
			count(*) filter (where decision='allowed'),
			count(*) filter (where decision='blocked'),
			count(*) filter (where decision='challenged'),
			coalesce(
				(count(*) filter (where decision='challenge_passed'))::double precision /
				nullif((count(*) filter (where decision='challenged'))::double precision, 0),
				0
			),
			count(distinct client_ip)
		from request_logs `+where, args...).Scan(
		&overview.IncomingRequests,
		&overview.AllowedRequests,
		&overview.BlockRequests,
		&overview.ChallengedRequests,
		&overview.ChallengePassRate,
		&overview.UniqueIPs,
	); err != nil {
		return overview, err
	}

	if err := s.db.QueryRow(ctx, `
		select coalesce(max(c), 0), coalesce(to_char(bucket, 'YYYY-MM-DD HH24:MI:SS'), '')
		from (
			select date_trunc('second', created_at) as bucket, count(*) as c
			from request_logs `+where+`
			group by bucket
		) t
	`, args...).Scan(&overview.PeakRPS, &overview.PeakTime); err != nil && err != pgx.ErrNoRows {
		return overview, err
	}

	loadRanked := func(query string) ([]domain.RankedMetric, error) {
		rows, err := s.db.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []domain.RankedMetric
		for rows.Next() {
			var item domain.RankedMetric
			if err := rows.Scan(&item.Name, &item.Value); err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, rows.Err()
	}

	var err error
	overview.TopIPs, err = loadRanked(`select client_ip, count(*) from request_logs ` + where + ` group by client_ip order by count(*) desc limit 10`)
	if err != nil {
		return overview, err
	}
	overview.TopUserAgents, err = loadRanked(`select coalesce(user_agent, 'unknown'), count(*) from request_logs ` + where + ` and decision='blocked' group by user_agent order by count(*) desc limit 10`)
	if err != nil {
		return overview, err
	}
	overview.TopDomains, err = loadRanked(`select coalesce(domain_name, 'unknown'), count(*) from request_logs ` + where + ` group by domain_name order by count(*) desc limit 10`)
	if err != nil {
		return overview, err
	}
	overview.TopReasons, err = loadRanked(`select coalesce(block_reason, 'allowed'), count(*) from request_logs ` + where + ` group by block_reason order by count(*) desc limit 10`)
	if err != nil {
		return overview, err
	}
	overview.TopCountries, err = loadRanked(`select coalesce(country_code, 'unknown'), count(*) from request_logs ` + where + ` group by country_code order by count(*) desc limit 10`)
	if err != nil {
		return overview, err
	}

	rows, err := s.db.Query(ctx, `
		select to_char(bucket, 'YYYY-MM-DD HH24:00') as label,
		       coalesce(sum(case when decision='allowed' then request_count else 0 end), 0) as allowed,
		       coalesce(sum(case when decision='blocked' then request_count else 0 end), 0) as blocked,
		       coalesce(sum(case when decision='challenged' then request_count else 0 end), 0) as challenged
		from (
			select date_trunc('hour', created_at) as bucket, decision, count(*) as request_count
			from request_logs `+where+`
			group by bucket, decision
		) t
		group by bucket
		order by bucket asc
	`, args...)
	if err != nil {
		return overview, err
	}
	defer rows.Close()

	for rows.Next() {
		var point domain.TimePoint
		if err := rows.Scan(&point.Label, &point.Allowed, &point.Blocked, &point.Challenge); err != nil {
			return overview, err
		}
		overview.RequestSeries = append(overview.RequestSeries, point)
	}
	return overview, rows.Err()
}
