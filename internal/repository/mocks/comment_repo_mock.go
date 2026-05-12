package mocks

import (
	"BAZ/internal/domain"
	"context"
	"errors"
)

type MockCommentRepository struct {
	Comments map[int64]*domain.Comment
	NextID   int64

	CreateError         error
	FindByIDError       error
	FindByBlogIDError   error
	UpdateError         error
	DeleteError         error
	DeleteByBlogIDError error
	DeleteByUserIDError error
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		Comments: make(map[int64]*domain.Comment),
		NextID:   1,
	}
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	comment.CommentID = m.NextID
	m.NextID++
	m.Comments[comment.CommentID] = comment
	return nil
}

func (m *MockCommentRepository) FindByID(ctx context.Context, id int64) (*domain.Comment, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}

	comment, exists := m.Comments[id]
	if !exists {
		return nil, nil
	}

	return comment, nil
}

func (m *MockCommentRepository) FindByBlogID(ctx context.Context, blogID int64, page, limit int) ([]domain.Comment, int64, error) {
	if m.FindByBlogIDError != nil {
		return nil, 0, m.FindByBlogIDError
	}

	var comments []domain.Comment
	for _, c := range m.Comments {
		if c.BlogID == blogID {
			comments = append(comments, *c)
		}
	}
	return comments, int64(len(comments)), nil
}

func (m *MockCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	existing, exists := m.Comments[comment.CommentID]
	if !exists {
		return errors.New("comment not found")
	}

	existing.Body = comment.Body
	existing.UpdatedAt = comment.UpdatedAt
	return nil
}

func (m *MockCommentRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	delete(m.Comments, id)
	return nil
}

func (m *MockCommentRepository) DeleteByBlogID(ctx context.Context, blogID int64) error {
	if m.DeleteByBlogIDError != nil {
		return m.DeleteByBlogIDError
	}

	for id, c := range m.Comments {
		if c.BlogID == blogID {
			delete(m.Comments, id)
		}
	}

	return nil
}

func (m *MockCommentRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	if m.DeleteByUserIDError != nil {
		return m.DeleteByUserIDError
	}

	for id, c := range m.Comments {
		if c.UserID == userID {
			delete(m.Comments, id)
		}
	}
	return nil
}
