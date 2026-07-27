package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sachinggsingh/PDTS/auth-service/internal/models"
	"github.com/sachinggsingh/PDTS/auth-service/internal/repository"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuthService struct {
	repo       repository.UserRepository
	logger     *utils.Logger
	jwtManager *utils.JWTManager
}

func NewOAuthService(repo repository.UserRepository, log *utils.Logger, jwtManager *utils.JWTManager) *OAuthService {
	return &OAuthService{repo: repo, logger: log, jwtManager: jwtManager}
}

// GetGithubAuthUrl generates the GitHub OAuth authorization URL
func (o *OAuthService) GetGithubAuthUrl(state string) (string, error) {
	if state == "" {
		return "", errors.New("state is empty")
	}

	config := utils.GetGitHubOAuthConfig()
	authURL := config.AuthCodeURL(state)

	o.logger.Info("Generated GitHub OAuth URL")
	return authURL, nil
}

// HandleGithubCallBack processes the GitHub OAuth
func (o *OAuthService) HandleGithubCallBack(code string, state string) (*models.User, error) {
	if code == "" {
		return nil, errors.New("authorization code is empty")
	}
	if state == "" {
		return nil, errors.New("state is empty")
	}

	// Exchange code for token
	config := utils.GetGitHubOAuthConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		o.logger.Error("Failed to exchange code for token: " + err.Error())
		return nil, errors.New("failed to exchange authorization code")
	}

	// Fetch user info from GitHub
	githubUser, err := utils.FetchGitHubUserInfo(token.AccessToken)
	if err != nil {
		o.logger.Error("Failed to fetch GitHub user info: " + err.Error())
		return nil, errors.New("failed to fetch user information from GitHub")
	}

	// Check if user already exists by OAuth ID
	oauthIDStr := strconv.FormatInt(githubUser.ID, 10)
	existingUser, err := o.repo.FindUserByOAuthId(oauthIDStr)
	if err != nil {
		o.logger.Error("Error checking for existing OAuth user: " + err.Error())
		return nil, err
	}

	// If user exists, update tokens and return
	if existingUser != nil {
		o.logger.Info("Existing OAuth user found, updating tokens")
		return o.updateUserTokens(existingUser)
	}

	// Check if user exists by email (link accounts)
	if githubUser.Email != "" {
		existingEmailUser, err := o.repo.FindUserByEmail(githubUser.Email)
		if err != nil {
			o.logger.Error("Error checking for existing email user: " + err.Error())
			return nil, err
		}

		// If user with this email exists, link OAuth to existing account
		if existingEmailUser != nil {
			o.logger.Info("Linking GitHub OAuth to existing email account")
			return o.linkOAuthToExistingUser(existingEmailUser, githubUser)
		}
	}

	// Create new user
	o.logger.Info("Creating new OAuth user")
	return o.createNewOAuthUser(githubUser)
}

// updateUserTokens generates new JWT tokens for existing user
func (o *OAuthService) updateUserTokens(user *models.User) (*models.User, error) {
	// Generate new JWT tokens
	tokenString, refreshTokenString, err := o.jwtManager.GenerateToken(
		user.User_id,
		*user.Email,
		"", // No password for OAuth users
		"",
	)
	if err != nil {
		o.logger.Error("Error generating tokens: " + err.Error())
		return nil, err
	}

	user.Token = &tokenString
	user.Refresh_Token = &refreshTokenString
	user.Updated_at = time.Now()

	// Update user in database
	updatedUser, err := o.repo.CreateAndUpdateOAuthUser(user)
	if err != nil {
		o.logger.Error("Error updating user: " + err.Error())
		return nil, err
	}

	return updatedUser, nil
}

// linkOAuthToExistingUser links GitHub OAuth to an existing email-based account
func (o *OAuthService) linkOAuthToExistingUser(user *models.User, githubUser *models.GitHubUser) (*models.User, error) {
	// Update user with OAuth information
	oauthProvider := "github"
	oauthID := strconv.FormatInt(githubUser.ID, 10)

	user.OAuthProvider = &oauthProvider
	user.OAuthID = &oauthID
	user.Avatar = &githubUser.AvatarURL
	user.GithubUsername = &githubUser.Login

	if githubUser.Name != "" {
		user.Name = &githubUser.Name
	}

	return o.updateUserTokens(user)
}

// createNewOAuthUser creates a new user from GitHub OAuth data
func (o *OAuthService) createNewOAuthUser(githubUser *models.GitHubUser) (*models.User, error) {
	if githubUser.Email == "" {
		return nil, errors.New("email is required but not provided by GitHub")
	}

	now := time.Now()
	oauthProvider := "github"
	oauthID := strconv.FormatInt(githubUser.ID, 10)

	user := &models.User{
		ID:             primitive.NewObjectID(),
		Email:          &githubUser.Email,
		OAuthProvider:  &oauthProvider,
		OAuthID:        &oauthID,
		Avatar:         &githubUser.AvatarURL,
		GithubUsername: &githubUser.Login,
		Created_at:     now,
		Updated_at:     now,
	}

	if githubUser.Name != "" {
		user.Name = &githubUser.Name
	}

	user.User_id = user.ID.Hex()

	// Generate JWT tokens
	tokenString, refreshTokenString, err := o.jwtManager.GenerateToken(
		user.User_id,
		*user.Email,
		"", // No password for OAuth users
		"",
	)
	if err != nil {
		o.logger.Error("Error generating tokens: " + err.Error())
		return nil, err
	}

	user.Token = &tokenString
	user.Refresh_Token = &refreshTokenString

	// Create user in database
	createdUser, err := o.repo.CreateAndUpdateOAuthUser(user)
	if err != nil {
		o.logger.Error("Error creating OAuth user: " + err.Error())
		return nil, err
	}

	o.logger.Info(fmt.Sprintf("Created new OAuth user: %s", *createdUser.Email))
	return createdUser, nil
}
