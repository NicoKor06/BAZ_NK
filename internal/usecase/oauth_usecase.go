package usecase

import (
	"BAZ/internal/oauth"
	"context"
	"errors"
	"log"
	"strings"

	"BAZ/internal/domain"
	"BAZ/internal/repository"
	"BAZ/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type OAuthUsecase struct {
	userRepo  repository.UserRepository
	jwtUtil   *utils.JWTUtil
	providers map[string]oauth.OAuthProvider
}

func NewOAuthUsecase(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) *OAuthUsecase {
	return &OAuthUsecase{
		userRepo:  userRepo,
		jwtUtil:   jwtUtil,
		providers: make(map[string]oauth.OAuthProvider),
	}
}

func (o *OAuthUsecase) RegisterProvider(p oauth.OAuthProvider) {
	o.providers[p.GetName()] = p
}

func (o *OAuthUsecase) GetAuthURL(providerName, state string) (string, error) {
	p, ok := o.providers[providerName]
	if !ok {
		return "", errors.New("provider not found")
	}
	return p.GetAuthURL(state), nil
}

func (o *OAuthUsecase) HandleCallback(ctx context.Context, providerName, code string) (*domain.OAuthLoginResponse, error) {
	p, ok := o.providers[providerName]
	if !ok {
		return nil, errors.New("provider not found")
	}

	oauthUser, err := p.ExchangeCode(ctx, code)
	if err != nil {
		log.Printf("OAuth exchange failed: %v", err)
		return nil, errors.New("failed to authenticate with provider")
	}

	user, err := o.userRepo.FindByProviderID(ctx, oauthUser.Provider, oauthUser.ProviderID)
	if err != nil {
		return nil, err
	}

	isNewUser := false
	if user == nil {
		existingUser, _ := o.userRepo.FindByEmail(ctx, oauthUser.Email)
		if existingUser != nil {
			user = existingUser
		} else {
			user = &domain.User{
				Username:  generateUsername(oauthUser.Email, oauthUser.Provider),
				Firstname: oauthUser.FirstName,
				Lastname:  oauthUser.LastName,
				Email:     oauthUser.Email,
				Password:  randomPasswordHash(),
				Role:      "user",
			}
			if err := o.userRepo.Create(ctx, user); err != nil {
				log.Printf("Create user failed: %v", err)
				return nil, err
			}
			log.Printf("User created with ID: %d", user.UserID)
			isNewUser = true
		}
		if err := o.userRepo.LinkProvider(ctx, user.UserID, oauthUser.Provider, oauthUser.ProviderID); err != nil {
			log.Printf("Failed to link provider: %v", err)
		}
	}

	token, err := o.jwtUtil.GenerateToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.OAuthLoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		IsNewUser:   isNewUser,
	}, nil
}

func generateUsername(email, provider string) string {
	local := strings.Split(email, "@")[0]
	return local + "_" + provider
}

func randomPasswordHash() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(""), bcrypt.DefaultCost)
	return string(hash)
}
