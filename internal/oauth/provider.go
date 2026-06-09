package oauth

import (
	"BAZ/internal/domain"
	"context"
)

type OAuthProvider interface {
	GetName() string
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*domain.OAuthUser, error)
}
