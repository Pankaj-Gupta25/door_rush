package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sachinggsingh/PDTS/auth-service/config"
	"github.com/sachinggsingh/PDTS/auth-service/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// GenerateStateToken generates a secure random state token for CSRF protection
func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetGitHubOAuthConfig returns GitHub OAuth configuration
func GetGitHubOAuthConfig() *oauth2.Config {
	env := config.GetEnv()
	return &oauth2.Config{
		ClientID:     env.GITHUB_CLIENT_ID,
		ClientSecret: env.GITHUB_CLIENT_SECRET,
		RedirectURL:  env.GITHUB_REDIRECT_URL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

// FetchGitHubUserInfo fetches user information from GitHub API
func FetchGitHubUserInfo(accessToken string) (*models.GitHubUser, error) {
	// Create HTTP client
	client := &http.Client{}

	// Fetch user info
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info from GitHub")
	}

	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// If email is not public, fetch from emails endpoint
	if user.Email == "" {
		email, err := fetchPrimaryEmail(accessToken, client)
		if err == nil {
			user.Email = email
		}
	}

	return &user, nil
}

// fetchPrimaryEmail fetches the primary verified email from GitHub
func fetchPrimaryEmail(accessToken string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("failed to fetch emails from GitHub: " + string(body))
	}

	var emails []models.GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// If no primary email, return first verified email
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", errors.New("no verified email found")
}
