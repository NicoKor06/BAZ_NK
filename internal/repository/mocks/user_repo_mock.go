package mocks

import (
	"BAZ/internal/domain"
	"context"
	"errors"
	"time"
)

type MockUserRepository struct {
	Users  map[int64]*domain.User
	NextID int64

	CreateError           error
	FindByIDError         error
	FindByUsernameError   error
	FindByEmailError      error
	UpdateError           error
	DeleteError           error
	UpdateLastOnlineError error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:  make(map[int64]*domain.User),
		NextID: 1,
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	user.UserID = m.NextID
	m.NextID++
	userCopy := *user
	m.Users[user.UserID] = &userCopy
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	if m.FindByIDError != nil {
		return nil, m.FindByIDError
	}

	user, exists := m.Users[id]
	if !exists {
		return nil, nil
	}

	userCopy := *user
	return &userCopy, nil
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.FindByUsernameError != nil {
		return nil, m.FindByUsernameError
	}

	for _, user := range m.Users {
		if user.Username == username {
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailError != nil {
		return nil, m.FindByEmailError
	}

	for _, user := range m.Users {
		if email == user.Email {
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	existing, exists := m.Users[user.UserID]
	if !exists {
		return errors.New("user not found")
	}

	existing.Firstname = user.Firstname
	existing.Lastname = user.Lastname
	existing.Email = user.Email
	existing.Birthday = user.Birthday
	existing.UpdatedAt = time.Now()
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	_, exists := m.Users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(m.Users, id)
	return nil
}

func (m *MockUserRepository) UpdateLastOnline(ctx context.Context, id int64) error {
	if m.UpdateLastOnlineError != nil {
		return m.UpdateLastOnlineError
	}

	user, exists := m.Users[id]
	if !exists {
		return errors.New("user not found")
	}

	user.LastOnline = time.Now()
	return nil
}
