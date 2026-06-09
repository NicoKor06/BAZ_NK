package oauth

import (
	oauth3 "BAZ/internal/utils"
	"context"
	"fmt"
	"strconv"

	"BAZ/internal/domain"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubProvider struct {
	config *oauth2.Config
}

type githubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (p *GitHubProvider) GetName() string {
	return "github"
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *GitHubProvider) ExchangeCode(ctx context.Context, code string) (*domain.OAuthUser, error) {
	token, err := p.config.Exchange(ctx, code, oauth2.SetAuthURLParam("redirect_uri", p.config.RedirectURL))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	userInfo, err := p.fetchUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	email := userInfo.Email
	if email == "" {
		email = userInfo.Login + "@users.noreply.github.com"
	}
	firstName, lastName := splitFullName(userInfo.Name)

	return &domain.OAuthUser{
		ProviderID:   strconv.Itoa(userInfo.ID),
		Provider:     p.GetName(),
		Email:        email,
		Name:         userInfo.Name,
		FirstName:    firstName,
		LastName:     lastName,
		AvatarURL:    userInfo.AvatarURL,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

func (p *GitHubProvider) fetchUserInfo(accessToken string) (*githubUserInfo, error) {
	url := "https://api.github.com/user"
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Accept":        "application/json",
	}
	var info githubUserInfo
	if err := oauth3.FetchJSON(url, headers, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func splitFullName(fullName string) (firstName, lastName string) {
	parts := splitName(fullName)
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return fullName, ""
}

func splitName(name string) []string {
	var result []string
	start := 0
	for i, c := range name {
		if c == ' ' {
			if start < i {
				result = append(result, name[start:i])
			}
			start = i + 1
		}
	}
	if start < len(name) {
		result = append(result, name[start:])
	}
	return result
}
