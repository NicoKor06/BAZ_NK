package test

import (
	"BAZ/internal/usecase"
	"context"
	"testing"
	"time"

	"BAZ/internal/domain"
	"BAZ/internal/repository/mock"
)

func TestCommentUsecase_Create(t *testing.T) {
	mockCommentRepo := mock.NewMockCommentRepository()
	mockBlogRepo := mock.NewMockBlogRepository()

	commentUsecase := usecase.NewCommentUsecase(mockCommentRepo, mockBlogRepo)
	ctx := context.Background()

	// Testdaten: Blog muss existieren
	existingBlog := &domain.Blog{
		BlogID:    1,
		Headline:  "Test Blog",
		Body:      "Test Inhalt",
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockBlogRepo.Blogs[1] = existingBlog

	t.Run("Create comment successfully", func(t *testing.T) {
		req := &domain.CreateCommentRequest{
			Body: "Das ist ein Test-Kommentar",
		}

		comment, err := commentUsecase.Create(ctx, 1, 1, req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if comment == nil {
			t.Errorf("Expected comment, got nil")
		}
		if comment.Body != "Das ist ein Test-Kommentar" {
			t.Errorf("Expected body 'Das ist ein Test-Kommentar', got '%s'", comment.Body)
		}
		if comment.BlogID != 1 {
			t.Errorf("Expected BlogID 1, got %d", comment.BlogID)
		}
		if comment.UserID != 1 {
			t.Errorf("Expected UserID 1, got %d", comment.UserID)
		}
		if comment.CommentID == 0 {
			t.Errorf("Expected CommentID to be set, got 0")
		}
	})

	t.Run("Create comment with empty body fails", func(t *testing.T) {
		req := &domain.CreateCommentRequest{
			Body: "",
		}

		comment, err := commentUsecase.Create(ctx, 1, 1, req)

		if err == nil {
			t.Errorf("Expected error for empty body, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})

	t.Run("Create comment for non-existing blog fails", func(t *testing.T) {
		req := &domain.CreateCommentRequest{
			Body: "Test",
		}

		comment, err := commentUsecase.Create(ctx, 999, 1, req)

		if err == nil {
			t.Errorf("Expected error for non-existing blog, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})

	t.Run("Create comment with invalid blogID fails", func(t *testing.T) {
		req := &domain.CreateCommentRequest{
			Body: "Test",
		}

		comment, err := commentUsecase.Create(ctx, 0, 1, req)

		if err == nil {
			t.Errorf("Expected error for invalid blogID, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})

	t.Run("Create comment with invalid userID fails", func(t *testing.T) {
		req := &domain.CreateCommentRequest{
			Body: "Test",
		}

		comment, err := commentUsecase.Create(ctx, 1, 0, req)

		if err == nil {
			t.Errorf("Expected error for invalid userID, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})
}

func TestCommentUsecase_GetByID(t *testing.T) {
	mockCommentRepo := mock.NewMockCommentRepository()
	mockBlogRepo := mock.NewMockBlogRepository()

	commentUsecase := usecase.NewCommentUsecase(mockCommentRepo, mockBlogRepo)
	ctx := context.Background()

	// Testdaten: Bestehenden Kommentar anlegen
	existingComment := &domain.Comment{
		CommentID: 1,
		Body:      "Vorhandener Kommentar",
		BlogID:    1,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockCommentRepo.Comments[1] = existingComment

	t.Run("Get existing comment by ID", func(t *testing.T) {
		comment, err := commentUsecase.GetByID(ctx, 1)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if comment == nil {
			t.Errorf("Expected comment, got nil")
		}
		if comment.CommentID != 1 {
			t.Errorf("Expected CommentID 1, got %d", comment.CommentID)
		}
		if comment.Body != "Vorhandener Kommentar" {
			t.Errorf("Expected body 'Vorhandener Kommentar', got '%s'", comment.Body)
		}
	})

	t.Run("Get non-existing comment returns error", func(t *testing.T) {
		comment, err := commentUsecase.GetByID(ctx, 999)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})

	t.Run("Get comment with invalid ID returns error", func(t *testing.T) {
		comment, err := commentUsecase.GetByID(ctx, 0)

		if err == nil {
			t.Errorf("Expected error for invalid ID, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})
}

func TestCommentUsecase_GetByBlogID(t *testing.T) {
	mockCommentRepo := mock.NewMockCommentRepository()
	mockBlogRepo := mock.NewMockBlogRepository()

	commentUsecase := usecase.NewCommentUsecase(mockCommentRepo, mockBlogRepo)
	ctx := context.Background()

	// Testdaten: Mehrere Kommentare für Blog 1
	mockCommentRepo.Comments[1] = &domain.Comment{CommentID: 1, Body: "Kommentar 1", BlogID: 1, UserID: 1}
	mockCommentRepo.Comments[2] = &domain.Comment{CommentID: 2, Body: "Kommentar 2", BlogID: 1, UserID: 2}
	mockCommentRepo.Comments[3] = &domain.Comment{CommentID: 3, Body: "Kommentar 3", BlogID: 2, UserID: 1}

	t.Run("Get all comments for blog 1", func(t *testing.T) {
		comments, total, err := commentUsecase.GetByBlogID(ctx, 1, 1, 10)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments, got %d", len(comments))
		}
		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}
	})

	t.Run("Get comments for blog with no comments", func(t *testing.T) {
		comments, total, err := commentUsecase.GetByBlogID(ctx, 3, 1, 10)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(comments) != 0 {
			t.Errorf("Expected 0 comments, got %d", len(comments))
		}
		if total != 0 {
			t.Errorf("Expected total 0, got %d", total)
		}
	})

	t.Run("Get comments with invalid blogID returns error", func(t *testing.T) {
		comments, total, err := commentUsecase.GetByBlogID(ctx, 0, 1, 10)

		if err == nil {
			t.Errorf("Expected error for invalid blogID, got nil")
		}
		if comments != nil {
			t.Errorf("Expected nil comments, got %v", comments)
		}
		if total != 0 {
			t.Errorf("Expected total 0, got %d", total)
		}
	})

	t.Run("Get comments with default pagination values", func(t *testing.T) {
		// page < 1 sollte auf 1 gesetzt werden
		comments, total, err := commentUsecase.GetByBlogID(ctx, 1, 0, 0)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments (default limit), got %d", len(comments))
		}
		if total != 2 {
			t.Errorf("Expected total 2, got %d", total)
		}
	})
}

func TestCommentUsecase_Update(t *testing.T) {
	mockCommentRepo := mock.NewMockCommentRepository()
	mockBlogRepo := mock.NewMockBlogRepository()

	commentUsecase := usecase.NewCommentUsecase(mockCommentRepo, mockBlogRepo)
	ctx := context.Background()

	// Testdaten: Bestehenden Kommentar (Autor: User 1)
	existingComment := &domain.Comment{
		CommentID: 1,
		Body:      "Alter Text",
		BlogID:    1,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockCommentRepo.Comments[1] = existingComment

	t.Run("Update own comment successfully", func(t *testing.T) {
		req := &domain.UpdateCommentRequest{
			Body: "Neuer Text",
		}

		comment, err := commentUsecase.Update(ctx, 1, 1, req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if comment == nil {
			t.Errorf("Expected comment, got nil")
		}
		if comment.Body != "Neuer Text" {
			t.Errorf("Expected body 'Neuer Text', got '%s'", comment.Body)
		}
	})

	t.Run("Update someone else's comment fails", func(t *testing.T) {
		req := &domain.UpdateCommentRequest{
			Body: "Fremder Text",
		}

		comment, err := commentUsecase.Update(ctx, 1, 999, req)

		if err == nil {
			t.Errorf("Expected error for updating foreign comment, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})

	t.Run("Update non-existing comment fails", func(t *testing.T) {
		req := &domain.UpdateCommentRequest{
			Body: "Text",
		}

		comment, err := commentUsecase.Update(ctx, 999, 1, req)

		if err == nil {
			t.Errorf("Expected error for non-existing comment, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})

	t.Run("Update with empty body fails", func(t *testing.T) {
		req := &domain.UpdateCommentRequest{
			Body: "",
		}

		comment, err := commentUsecase.Update(ctx, 1, 1, req)

		if err == nil {
			t.Errorf("Expected error for empty body, got nil")
		}
		if comment != nil {
			t.Errorf("Expected nil comment, got %v", comment)
		}
	})
}

func TestCommentUsecase_Delete(t *testing.T) {
	mockCommentRepo := mock.NewMockCommentRepository()
	mockBlogRepo := mock.NewMockBlogRepository()

	commentUsecase := usecase.NewCommentUsecase(mockCommentRepo, mockBlogRepo)
	ctx := context.Background()

	// Testdaten: Bestehenden Kommentar (Autor: User 1)
	existingComment := &domain.Comment{
		CommentID: 1,
		Body:      "Zu löschender Kommentar",
		BlogID:    1,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockCommentRepo.Comments[1] = existingComment

	t.Run("Delete own comment successfully", func(t *testing.T) {
		err := commentUsecase.Delete(ctx, 1, 1)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Prüfen ob Kommentar gelöscht wurde
		if _, exists := mockCommentRepo.Comments[1]; exists {
			t.Errorf("Comment should be deleted")
		}
	})

	t.Run("Delete someone else's comment fails", func(t *testing.T) {
		// Neuen Kommentar für diesen Test anlegen
		mockCommentRepo.Comments[2] = &domain.Comment{
			CommentID: 2,
			Body:      "Fremder Kommentar",
			BlogID:    1,
			UserID:    999,
		}

		err := commentUsecase.Delete(ctx, 2, 1)

		if err == nil {
			t.Errorf("Expected error for deleting foreign comment, got nil")
		}

		// Kommentar sollte noch existieren
		if _, exists := mockCommentRepo.Comments[2]; !exists {
			t.Errorf("Comment should still exist after failed delete")
		}
	})

	t.Run("Delete non-existing comment fails", func(t *testing.T) {
		err := commentUsecase.Delete(ctx, 999, 1)

		if err == nil {
			t.Errorf("Expected error for non-existing comment, got nil")
		}
	})

	t.Run("Delete with invalid commentID fails", func(t *testing.T) {
		err := commentUsecase.Delete(ctx, 0, 1)

		if err == nil {
			t.Errorf("Expected error for invalid commentID, got nil")
		}
	})

	t.Run("Delete with invalid userID fails", func(t *testing.T) {
		err := commentUsecase.Delete(ctx, 1, 0)

		if err == nil {
			t.Errorf("Expected error for invalid userID, got nil")
		}
	})
}
