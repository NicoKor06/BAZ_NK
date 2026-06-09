package test

import (
	"BAZ/internal/cache/mocks"
	"BAZ/internal/usecase"
	"context"
	"testing"
	"time"

	"BAZ/internal/domain"
	"BAZ/internal/repository/mock"
)

func createTestBlog(id int64, userID int64, headline string) *domain.Blog {
	return &domain.Blog{
		BlogID:    id,
		Headline:  headline,
		Body:      "Test Inhalt",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}
}

func createTestRequest(headline, body string) *domain.CreateBlogRequest {
	return &domain.CreateBlogRequest{
		Headline: headline,
		Body:     body,
	}
}

func TestBlogUsecase_Create(t *testing.T) {
	mockBlogRepo := mock.NewMockBlogRepository()
	mockCommentRepo := mock.NewMockCommentRepository()
	mockCache := mocks.NewMockCache()

	blogUsecase := usecase.NewBlogUsecase(mockBlogRepo, mockCommentRepo, mockCache)
	ctx := context.Background()

	t.Run("Create blog successfully", func(t *testing.T) {
		req := createTestRequest("Mein erster Blog", "Das ist der Inhalt")
		userID := int64(1)

		blog, err := blogUsecase.Create(ctx, userID, req)

		// Überprüfungen
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if blog == nil {
			t.Errorf("Expected blog, got nil")
		}

		if blog.Headline != "Mein erster Blog" {
			t.Errorf("Expected headline 'Mein erster Blog', got '%s'", blog.Headline)
		}

		if blog.Body != "Das ist der Inhalt" {
			t.Errorf("Expected body 'Das ist der Inhalt', got '%s'", blog.Body)
		}

		if blog.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, blog.UserID)
		}

		if blog.BlogID == 0 {
			t.Errorf("Expected BlogID to be set, got 0")
		}

		// Prüfen, ob der Blog wirklich im Mock gespeichert wurde
		savedBlog, _ := mockBlogRepo.FindByID(ctx, blog.BlogID)
		if savedBlog == nil {
			t.Errorf("Blog was not saved in repository")
		}
	})

	t.Run("Create fails when repository returns error", func(t *testing.T) {
		// Mock so konfigurieren, dass es einen Fehler zurückgibt
		mockBlogRepo.CreateError = nil // Reset
		mockBlogRepo.CreateError = &repositoryError{msg: "database error"}
		defer func() { mockBlogRepo.CreateError = nil }()

		req := createTestRequest("Test", "Test")
		userID := int64(1)

		blog, err := blogUsecase.Create(ctx, userID, req)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if blog != nil {
			t.Errorf("Expected nil blog, got %v", blog)
		}
	})
}

func TestBlogUsecase_GetByID(t *testing.T) {
	mockBlogRepo := mock.NewMockBlogRepository()
	mockCommentRepo := mock.NewMockCommentRepository()
	mockCache := mocks.NewMockCache()

	blogUsecase := usecase.NewBlogUsecase(mockBlogRepo, mockCommentRepo, mockCache)
	ctx := context.Background()

	existingBlog := createTestBlog(1, 1, "Vorhandener Blog")
	mockBlogRepo.Blogs[1] = existingBlog

	t.Run("Get existing blog by ID", func(t *testing.T) {
		blog, err := blogUsecase.GetByID(ctx, 1)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if blog == nil {
			t.Errorf("Expected blog, got nil")
		}

		if blog.BlogID != 1 {
			t.Errorf("Expected BlogID 1, got %d", blog.BlogID)
		}

		if blog.Headline != "Vorhandener Blog" {
			t.Errorf("Expected headline 'Vorhandener Blog', got '%s'", blog.Headline)
		}
	})

	t.Run("Get non-existing blog returns nil", func(t *testing.T) {
		blog, err := blogUsecase.GetByID(ctx, 999)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if blog != nil {
			t.Errorf("Expected nil blog, got %v", blog)
		}
	})
}

func TestBlogUsecase_GetAll(t *testing.T) {
	mockBlogRepo := mock.NewMockBlogRepository()
	mockCommentRepo := mock.NewMockCommentRepository()
	mockCache := mocks.NewMockCache()

	blogUsecase := usecase.NewBlogUsecase(mockBlogRepo, mockCommentRepo, mockCache)
	ctx := context.Background()

	for i := 1; i <= 25; i++ {
		blog := createTestBlog(int64(i), int64(i%3+1), "Blog "+string(rune(i)))
		mockBlogRepo.Blogs[int64(i)] = blog
	}

	t.Run("Get page 1 with limit 10", func(t *testing.T) {
		blogs, total, err := blogUsecase.GetAll(ctx, 1, 10)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(blogs) != 10 {
			t.Errorf("Expected 10 blogs, got %d", len(blogs))
		}

		if total != 25 {
			t.Errorf("Expected total 25, got %d", total)
		}
	})

	t.Run("Get page 2 with limit 10", func(t *testing.T) {
		blogs, total, err := blogUsecase.GetAll(ctx, 2, 10)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(blogs) != 10 {
			t.Errorf("Expected 10 blogs, got %d", len(blogs))
		}

		if total != 25 {
			t.Errorf("Expected total 25, got %d", total)
		}
	})

	t.Run("Get page 10 (beyond range) returns empty slice", func(t *testing.T) {
		blogs, total, err := blogUsecase.GetAll(ctx, 10, 10)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(blogs) != 0 {
			t.Errorf("Expected 0 blogs, got %d", len(blogs))
		}

		if total != 25 {
			t.Errorf("Expected total 25, got %d", total)
		}
	})

	t.Run("Negative page and limit use defaults", func(t *testing.T) {
		blogs, total, err := blogUsecase.GetAll(ctx, 0, 0)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(blogs) != 10 {
			t.Errorf("Expected 10 blogs (default limit), got %d", len(blogs))
		}

		if total != 25 {
			t.Errorf("Expected total 25, got %d", total)
		}
	})
}

type repositoryError struct {
	msg string
}

func (e *repositoryError) Error() string {
	return e.msg
}
