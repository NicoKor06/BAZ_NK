package mocks

import (
	"BAZ/internal/domain"
	"context"
	"errors"
)

type MockBlogRepository struct {
	Blogs  map[int64]*domain.Blog
	NextID int64

	CreateError         error
	FindByIDError       error
	FindAllError        error
	UpdateError         error
	DeleteError         error
	DeleteByUserIDError error
}

func NewMockBlogRepository() *MockBlogRepository {
	return &MockBlogRepository{
		Blogs:  make(map[int64]*domain.Blog),
		NextID: 1,
	}
}

func (m *MockBlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	blog.BlogID = m.NextID
	m.NextID++
	m.Blogs[blog.BlogID] = blog
	return nil
}

func (m *MockBlogRepository) FindByID(ctx context.Context, id int64) (*domain.Blog, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}

	blog, exists := m.Blogs[id]
	if !exists {
		return nil, nil
	}

	return blog, nil
}

func (m *MockBlogRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Blog, int64, error) {
	if m.FindAllError != nil {
		return nil, 0, m.FindAllError
	}

	allBlogs := make([]domain.Blog, 0, len(m.Blogs))
	for _, blog := range m.Blogs {
		allBlogs = append(allBlogs, *blog)
	}

	total := int64(len(allBlogs))

	offset := (page - 1) * limit
	if offset >= len(allBlogs) {
		return []domain.Blog{}, total, nil
	}

	end := offset + limit
	if end > len(allBlogs) {
		end = len(allBlogs)
	}

	return allBlogs[offset:end], total, nil
}

func (m *MockBlogRepository) Update(ctx context.Context, blog *domain.Blog) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	existing, exists := m.Blogs[blog.BlogID]
	if !exists {
		return errors.New("blog not found")
	}

	existing.Headline = blog.Headline
	existing.Body = blog.Body
	existing.UpdatedAt = blog.UpdatedAt
	return nil
}

func (m *MockBlogRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	delete(m.Blogs, id)
	return nil
}

func (m *MockBlogRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	if m.DeleteByUserIDError != nil {
		return m.DeleteByUserIDError
	}

	for id, blog := range m.Blogs {
		if blog.UserID == userID {
			delete(m.Blogs, id)
		}
	}
	return nil
}
