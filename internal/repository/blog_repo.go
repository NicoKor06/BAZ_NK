package repository

import (
	"BAZ/internal/domain"
	"context"
)

type BlogRepository interface {
	Create(ctx context.Context, blog *domain.Blog) error
	FindByID(ctx context.Context, id int64) (*domain.Blog, error)
	FindAll(ctx context.Context, page, limit int) ([]domain.Blog, int64, error)
	Update(ctx context.Context, blog *domain.Blog) error
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
}
