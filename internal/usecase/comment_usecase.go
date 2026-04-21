package usecase

import (
	"BAZ/internal/domain"
	"BAZ/internal/repository"
	"context"
	"errors"
	"time"
)

type CommentUsecase struct {
	commentRepo repository.CommentRepository
	blogRepo    repository.BlogRepository
}

func NewCommentUsecase(commentRepo repository.CommentRepository, blogRepo repository.BlogRepository) *CommentUsecase {
	return &CommentUsecase{
		commentRepo: commentRepo,
		blogRepo:    blogRepo,
	}
}

func (c *CommentUsecase) Create(ctx context.Context, blogID, userID int64, req *domain.CreateCommentRequest) (*domain.Comment, error) {
	// Prüfen ob Blog existiert
	_, err := c.blogRepo.FindByID(ctx, blogID)
	if err != nil {
		return nil, errors.New("blog not found")
	}

	now := time.Now()
	comment := &domain.Comment{
		Body:      req.Body,
		CreatedAt: now,
		UpdatedAt: now,
		BlogID:    blogID,
		UserID:    userID,
	}

	err = c.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *CommentUsecase) GetByID(ctx context.Context, commentID int64) (*domain.Comment, error) {
	return c.commentRepo.FindByID(ctx, commentID)
}

func (c *CommentUsecase) GetByBlogID(ctx context.Context, blogID int64, page, limit int) ([]domain.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return c.commentRepo.FindByBlogID(ctx, blogID, page, limit)
}

func (c *CommentUsecase) Update(ctx context.Context, commentID, userID int64, req *domain.UpdateCommentRequest) (*domain.Comment, error) {
	comment, err := c.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	// Prüfen ob User der Autor ist
	if comment.UserID != userID {
		return nil, errors.New("you are not the author of this comment")
	}

	comment.Body = req.Body
	comment.UpdatedAt = time.Now()

	err = c.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (c *CommentUsecase) Delete(ctx context.Context, commentID, userID int64) error {
	comment, err := c.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Prüfen ob User der Autor ist
	if comment.UserID != userID {
		return errors.New("you are not the author of this comment")
	}

	return c.commentRepo.Delete(ctx, commentID)
}
