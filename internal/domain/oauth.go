package domain

import "time"

type OAuthUser struct {
	ProviderID   string    `json:"providerId"`
	Provider     string    `json:"provider"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	AvatarURL    string    `json:"avatarUrl"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type OAuthLoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IsNewUser   bool   `json:"is_new_user"`
}
