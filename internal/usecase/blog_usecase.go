package usecase

import (
	"BAZ/internal/domain"
	"BAZ/internal/repository"
	"context"
	"errors"
	"time"
)

type BlogUsecase struct {
	blogRepo    repository.BlogRepository
	commentRepo repository.CommentRepository
}

func NewBlogUsecase(blogRepo repository.BlogRepository, commentRepo repository.CommentRepository) *BlogUsecase {
	return &BlogUsecase{
		blogRepo:    blogRepo,
		commentRepo: commentRepo,
	}
}

func (b *BlogUsecase) Create(ctx context.Context, userID int64, req *domain.CreateBlogRequest) (*domain.Blog, error) {
	now := time.Now()
	blog := &domain.Blog{
		Headline:  req.Headline,
		Body:      req.Body,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
	}

	err := b.blogRepo.Create(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (b *BlogUsecase) GetByID(ctx context.Context, blogID int64) (*domain.Blog, error) {
	return b.blogRepo.FindByID(ctx, blogID)
}

func (b *BlogUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Blog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return b.blogRepo.FindAll(ctx, page, limit)
}

func (b *BlogUsecase) Update(ctx context.Context, blogID, userID int64, req *domain.UpdateBlogRequest) (*domain.Blog, error) {
	blog, err := b.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return nil, err
	}

	// Prüfen ob User der Autor ist
	if blog.UserID != userID {
		return nil, errors.New("you are not the author of this blog")
	}

	if req.Headline != "" {
		blog.Headline = req.Headline
	}
	if req.Body != "" {
		blog.Body = req.Body
	}
	blog.UpdatedAt = time.Now()

	err = b.blogRepo.Update(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (b *BlogUsecase) Delete(ctx context.Context, blogID, userID int64) error {
	blog, err := b.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return err
	}

	// Prüfen ob User der Autor ist
	if blog.UserID != userID {
		return errors.New("you are not the author of this blog")
	}

	// Erst Comments löschen, dann Blog
	_ = b.commentRepo.DeleteByBlogID(ctx, blogID)
	return b.blogRepo.Delete(ctx, blogID)
}
