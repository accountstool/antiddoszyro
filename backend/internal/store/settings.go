package store

import (
	"context"

	"shieldpanel/backend/internal/domain"
)

func (s *Store) ListSettings(ctx context.Context) ([]domain.SystemSetting, error) {
	rows, err := s.db.Query(ctx, `
		select key, value, type, created_at, updated_at
		from system_settings
		order by key asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []domain.SystemSetting
	for rows.Next() {
		var item domain.SystemSetting
		if err := rows.Scan(&item.Key, &item.Value, &item.Type, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, item)
	}
	return settings, rows.Err()
}

func (s *Store) UpsertSetting(ctx context.Context, key string, value string, valueType string) error {
	_, err := s.db.Exec(ctx, `
		insert into system_settings (key, value, type)
		values ($1, $2, $3)
		on conflict (key)
		do update set value=excluded.value, type=excluded.type, updated_at=now()
	`, key, value, valueType)
	return err
}
