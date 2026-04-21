package repository

import (
	"BAZ/internal/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
	UpdateLastOnline(ctx context.Context, id int64) error
}
