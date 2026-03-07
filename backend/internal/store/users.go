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

type CreateUserParams struct {
	Username     string
	Email        string
	PasswordHash string
	DisplayName  string
	Role         domain.UserRole
	Language     string
}

type UpdateUserParams struct {
	ID          uuid.UUID
	Username    string
	Email       string
	DisplayName string
	Role        domain.UserRole
	Language    string
}

func (s *Store) CreateUser(ctx context.Context, params CreateUserParams) (domain.User, error) {
	return scanUser(s.db.QueryRow(ctx, `
		insert into users (username, email, password_hash, display_name, role, language)
		values ($1, $2, $3, $4, $5, $6)
		returning id, username, email, password_hash, display_name, role, language, last_login_at, created_at, updated_at
	`, strings.ToLower(params.Username), strings.ToLower(params.Email), params.PasswordHash, params.DisplayName, params.Role, params.Language))
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return scanUser(s.db.QueryRow(ctx, `
		select id, username, email, password_hash, display_name, role, language, last_login_at, created_at, updated_at
		from users
		where id = $1
	`, id))
}

func (s *Store) FindUserByIdentifier(ctx context.Context, identifier string) (domain.User, error) {
	identifier = strings.ToLower(strings.TrimSpace(identifier))
	return scanUser(s.db.QueryRow(ctx, `
		select id, username, email, password_hash, display_name, role, language, last_login_at, created_at, updated_at
		from users
		where lower(username) = $1 or lower(email) = $1
	`, identifier))
}

func (s *Store) ListUsers(ctx context.Context, limit int, offset int) ([]domain.User, int64, error) {
	rows, err := s.db.Query(ctx, `
		select id, username, email, password_hash, display_name, role, language, last_login_at, created_at, updated_at
		from users
		order by created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]domain.User, 0, limit)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	var total int64
	if err := s.db.QueryRow(ctx, `select count(*) from users`).Scan(&total); err != nil {
		return nil, 0, err
	}
	return users, total, rows.Err()
}

func (s *Store) UpdateUser(ctx context.Context, params UpdateUserParams) (domain.User, error) {
	return scanUser(s.db.QueryRow(ctx, `
		update users
		set username=$2,
			email=$3,
			display_name=$4,
			role=$5,
			language=$6,
			updated_at=now()
		where id=$1
		returning id, username, email, password_hash, display_name, role, language, last_login_at, created_at, updated_at
	`, params.ID, strings.ToLower(params.Username), strings.ToLower(params.Email), params.DisplayName, params.Role, params.Language))
}

func (s *Store) UpdateUserPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	_, err := s.db.Exec(ctx, `
		update users
		set password_hash=$2, updated_at=now()
		where id=$1
	`, id, passwordHash)
	return err
}

func (s *Store) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tag, err := s.db.Exec(ctx, `delete from users where id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (s *Store) SetUserLastLogin(ctx context.Context, id uuid.UUID, when time.Time) error {
	_, err := s.db.Exec(ctx, `update users set last_login_at=$2 where id=$1`, id, when)
	return err
}

func isNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
