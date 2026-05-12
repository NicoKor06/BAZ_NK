package usecase

import (
	"context"
	"errors"
	"time"

	"BAZ/internal/cache"
	"BAZ/internal/domain"
	"BAZ/internal/repository"
)

type BlogUsecase struct {
	blogRepo    repository.BlogRepository
	commentRepo repository.CommentRepository
	cache       cache.Cache
}

func NewBlogUsecase(
	blogRepo repository.BlogRepository,
	commentRepo repository.CommentRepository,
	cache cache.Cache,
) *BlogUsecase {
	return &BlogUsecase{
		blogRepo:    blogRepo,
		commentRepo: commentRepo,
		cache:       cache,
	}
}

// GetByID – Holt einen Blog aus dem Repository (nicht sich selbst!)
func (b *BlogUsecase) GetByID(ctx context.Context, blogID int64) (*domain.Blog, error) {
	return b.blogRepo.FindByID(ctx, blogID)
}

// GetAll – Holt alle Blogs (mit Paginierung)
func (b *BlogUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Blog, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return b.blogRepo.FindAll(ctx, page, limit)
}

// Create – Erstellt einen neuen Blog
func (b *BlogUsecase) Create(ctx context.Context, userID int64, req *domain.CreateBlogRequest) (*domain.Blog, error) {
	blog := &domain.Blog{
		Headline:  req.Headline,
		Body:      req.Body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}

	err := b.blogRepo.Create(ctx, blog)
	if err != nil {
		return nil, err
	}

	// Cache invalidieren
	b.cache.Delete(ctx, "/blogs")
	return blog, nil
}

// Update – Aktualisiert einen Blog
func (b *BlogUsecase) Update(ctx context.Context, blogID, userID int64, req *domain.UpdateBlogRequest) (*domain.Blog, error) {
	blog, err := b.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return nil, err
	}
	if blog.UserID != userID {
		return nil, errors.New("you are not the author")
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

	// Cache invalidieren
	b.cache.Delete(ctx, "/blogs")
	return blog, nil
}

// Delete – Löscht einen Blog
func (b *BlogUsecase) Delete(ctx context.Context, blogID, userID int64) error {
	blog, err := b.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return err
	}
	if blog.UserID != userID {
		return errors.New("you are not the author")
	}

	if err := b.blogRepo.Delete(ctx, blogID); err != nil {
		return err
	}

	// Cache invalidieren
	b.cache.Delete(ctx, "/blogs")
	return nil
}
