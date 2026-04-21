package repository

import (
	"BAZ/internal/domain"
	"context"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	FindByID(ctx context.Context, id int64) (*domain.Comment, error)
	FindByBlogID(ctx context.Context, blogID int64, page, limit int) ([]domain.Comment, int64, error)
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id int64) error
	DeleteByBlogID(ctx context.Context, blogID int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
}
