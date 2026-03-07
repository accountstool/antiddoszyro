package services

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"shieldpanel/backend/internal/config"
	"shieldpanel/backend/internal/domain"
	"shieldpanel/backend/internal/store"
	"shieldpanel/backend/internal/util"
)

type UsersService struct {
	store *store.Store
	audit *AuditService
	cfg   config.Config
}

func NewUsersService(repo *store.Store, audit *AuditService, cfg config.Config) *UsersService {
	return &UsersService{
		store: repo,
		audit: audit,
		cfg:   cfg,
	}
}

func (s *UsersService) List(ctx context.Context, limit int, offset int) ([]domain.User, int64, error) {
	return s.store.ListUsers(ctx, limit, offset)
}

func (s *UsersService) Create(ctx context.Context, actor *domain.User, username string, email string, displayName string, role domain.UserRole, language string, password string, ipAddress string, userAgent string) (domain.User, error) {
	hash, err := util.HashPassword(password, s.cfg.Auth.BcryptCost)
	if err != nil {
		return domain.User{}, err
	}
	item, err := s.store.CreateUser(ctx, store.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		DisplayName:  displayName,
		Role:         role,
		Language:     language,
	})
	if err != nil {
		return domain.User{}, err
	}
	details, _ := json.Marshal(map[string]any{"username": item.Username, "role": item.Role})
	s.audit.Record(ctx, actor, "create_user", "user", item.ID.String(), ipAddress, userAgent, string(details))
	return item, nil
}

func (s *UsersService) Update(ctx context.Context, actor *domain.User, input store.UpdateUserParams, password string, ipAddress string, userAgent string) (domain.User, error) {
	item, err := s.store.UpdateUser(ctx, input)
	if err != nil {
		return domain.User{}, err
	}
	if password != "" {
		hash, err := util.HashPassword(password, s.cfg.Auth.BcryptCost)
		if err != nil {
			return domain.User{}, err
		}
		if err := s.store.UpdateUserPassword(ctx, input.ID, hash); err != nil {
			return domain.User{}, err
		}
	}
	details, _ := json.Marshal(map[string]any{"username": item.Username, "role": item.Role})
	s.audit.Record(ctx, actor, "update_user", "user", item.ID.String(), ipAddress, userAgent, string(details))
	return s.store.GetUserByID(ctx, input.ID)
}

func (s *UsersService) Delete(ctx context.Context, actor *domain.User, id uuid.UUID, ipAddress string, userAgent string) error {
	if err := s.store.DeleteUser(ctx, id); err != nil {
		return err
	}
	s.audit.Record(ctx, actor, "delete_user", "user", id.String(), ipAddress, userAgent, "")
	return nil
}
