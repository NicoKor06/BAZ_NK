package usecase

import (
	"BAZ/internal/domain"
	"BAZ/internal/repository"
	"BAZ/internal/utils"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (a *AuthUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Prüfen ob Username existiert
	existingUser, _ := a.userRepo.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Prüfen ob Email existiert
	existingEmail, _ := a.userRepo.FindByEmail(ctx, req.Email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Passwort hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		Username:   req.Username,
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Birthday:   req.Birthday,
		Role:       "user",
		CreatedAt:  now,
		UpdatedAt:  now,
		LastOnline: now,
	}

	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// User finden
	user, err := a.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Passwort prüfen
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// LastOnline updaten
	_ = a.userRepo.UpdateLastOnline(ctx, user.UserID)

	// Token generieren
	token, err := a.jwtUtil.GenerateToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}
