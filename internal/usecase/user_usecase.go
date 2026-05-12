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
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *UserUsecase) GetPublicProfile(ctx context.Context, userID int64) (*domain.UserPublic, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
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
	if userID == 0 {
		return errors.New("invalid user id")
	}

	if password == "" {
		return errors.New("required password")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	// Cascade Delete: erst Comments, dann Blogs, dann User
	if err := u.commentRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}

	if err := u.blogRepo.DeleteByUserID(ctx, userID); err != nil {
		return err
	}

	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) UpdateLastOnline(ctx context.Context, userID int64) error {
	if userID == 0 {
		return errors.New("invalid user id")
	}
	return u.userRepo.UpdateLastOnline(ctx, userID)
}
