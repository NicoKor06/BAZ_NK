package usecase

import (
	"BAZ/internal/domain"
	"BAZ/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo    repository.UserRepository
	blogRepo    repository.BlogRepository
	commentRepo repository.CommentRepository
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	blogRepo repository.BlogRepository,
	commentRepo repository.CommentRepository,
) *UserUsecase {
	return &UserUsecase{
		userRepo:    userRepo,
		blogRepo:    blogRepo,
		commentRepo: commentRepo,
	}
}

func (u *UserUsecase) GetOwnProfile(ctx context.Context, userID int64) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}

func (u *UserUsecase) GetPublicProfile(ctx context.Context, userID int64) (*domain.UserPublic, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserPublic{
		UserID:     user.UserID,
		Username:   user.Username,
		Firstname:  user.Firstname,
		Lastname:   user.Lastname,
		Role:       user.Role,
		LastOnline: user.LastOnline,
	}, nil
}

func (u *UserUsecase) DeleteOwnAccount(ctx context.Context, userID int64, password string) error {
	// User finden
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Passwort prüfen (zusätzliche Sicherheit)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	// Cascade Delete: erst Comments, dann Blogs, dann User
	_ = u.commentRepo.DeleteByUserID(ctx, userID)
	_ = u.blogRepo.DeleteByUserID(ctx, userID)

	return u.userRepo.Delete(ctx, userID)
}
