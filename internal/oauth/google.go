package oauth

import (
	oauth3 "BAZ/internal/utils"
	"context"
	"fmt"

	"BAZ/internal/domain"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	config *oauth2.Config
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) *GoogleProvider {
	return &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (p *GoogleProvider) GetName() string { return "google" }

func (p *GoogleProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*domain.OAuthUser, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	userInfo, err := p.fetchUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	return &domain.OAuthUser{
		ProviderID:   userInfo.ID,
		Provider:     p.GetName(),
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		FirstName:    userInfo.GivenName,
		LastName:     userInfo.FamilyName,
		AvatarURL:    userInfo.Picture,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

func (p *GoogleProvider) fetchUserInfo(accessToken string) (*googleUserInfo, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo"
	headers := map[string]string{"Authorization": "Bearer " + accessToken}
	var info googleUserInfo
	if err := oauth3.FetchJSON(url, headers, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
