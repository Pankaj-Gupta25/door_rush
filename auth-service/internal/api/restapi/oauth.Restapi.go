package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/auth-service/internal/service"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
)

type OAuthRestAPI struct {
	Logger       *utils.Logger
	Router       *gin.Engine
	OAuthService *service.OAuthService
}

func NewOAuthRestAPI(logger *utils.Logger, router *gin.Engine, oauthService *service.OAuthService) *OAuthRestAPI {
	return &OAuthRestAPI{
		Logger:       logger,
		Router:       router,
		OAuthService: oauthService,
	}
}

func (o *OAuthRestAPI) SetupOAuthRoutes() {
	o.Router.GET("/auth/github", o.GitHubAuthHandler)
	o.Router.GET("/auth/github/callback", o.GitHubCallbackHandler)
}

// GitHubAuthHandler initiates the GitHub OAuth flow
func (o *OAuthRestAPI) GitHubAuthHandler(c *gin.Context) {
	// Generate state token for CSRF protection
	state, err := utils.GenerateStateToken()
	if err != nil {
		o.Logger.Error("Failed to generate state token: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate OAuth flow"})
		return
	}

	// Store state in session/cookie for validation (in production, use secure session storage)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true) // 5 minutes expiry

	// Get GitHub OAuth URL
	authURL, err := o.OAuthService.GetGithubAuthUrl(state)
	if err != nil {
		o.Logger.Error("Failed to generate GitHub auth URL: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate OAuth flow"})
		return
	}

	// Redirect to GitHub
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallbackHandler handles the GitHub OAuth callback
func (o *OAuthRestAPI) GitHubCallbackHandler(c *gin.Context) {
	// Get code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	// Check for OAuth errors
	if errorParam := c.Query("error"); errorParam != "" {
		errorDesc := c.Query("error_description")
		o.Logger.Warn("GitHub OAuth error: " + errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             errorParam,
			"error_description": errorDesc,
		})
		return
	}

	// Validate state token (CSRF protection)
	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		o.Logger.Warn("Invalid state token - possible CSRF attack")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid state token"})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Handle OAuth callback
	user, err := o.OAuthService.HandleGithubCallBack(code, state)
	if err != nil {
		o.Logger.Error("Failed to handle GitHub callback: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return user data and tokens
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully authenticated with GitHub",
		"user": gin.H{
			"user_id":         user.User_id,
			"email":           user.Email,
			"name":            user.Name,
			"github_username": user.GithubUsername,
			"avatar":          user.Avatar,
			"oauth_provider":  user.OAuthProvider,
		},
		"token":         user.Token,
		"refresh_token": user.Refresh_Token,
	})
}
