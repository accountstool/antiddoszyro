package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migrator struct {
	db            *pgxpool.Pool
	migrationsDir string
}

func NewMigrator(db *pgxpool.Pool, migrationsDir string) *Migrator {
	return &Migrator{db: db, migrationsDir: migrationsDir}
}

func (m *Migrator) Run(ctx context.Context) error {
	if _, err := m.db.Exec(ctx, `
		create table if not exists schema_migrations (
			version text primary key,
			applied_at timestamptz not null default now()
		)
	`); err != nil {
		return err
	}

	entries, err := os.ReadDir(m.migrationsDir)
	if err != nil {
		return err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	for _, file := range files {
		var exists bool
		if err := m.db.QueryRow(ctx, `select exists(select 1 from schema_migrations where version=$1)`, file).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		content, err := os.ReadFile(filepath.Join(m.migrationsDir, file))
		if err != nil {
			return err
		}
		tx, err := m.db.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migration %s failed: %w", file, err)
		}
		if _, err := tx.Exec(ctx, `insert into schema_migrations (version) values ($1)`, file); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}
