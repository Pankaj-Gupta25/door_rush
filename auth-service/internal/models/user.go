package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Password *string            `json:"password,omitempty" bson:"password,omitempty" validate:"omitempty,min=6"`
	Email    *string            `json:"email" bson:"email" validate:"email,required"`

	// OAuth fields (optional - only for OAuth users)
	OAuthProvider  *string `json:"oauth_provider,omitempty" bson:"oauth_provider,omitempty"`
	OAuthID        *string `json:"oauth_id,omitempty" bson:"oauth_id,omitempty"`
	Avatar         *string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Name           *string `json:"name,omitempty" bson:"name,omitempty"`
	GithubUsername *string `json:"github_username,omitempty" bson:"github_username,omitempty"`

	Token         *string   `json:"token,omitempty" bson:"token,omitempty"`
	Refresh_Token *string   `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	Created_at    time.Time `json:"created_at" bson:"created_at"`
	Updated_at    time.Time `json:"updated_at" bson:"updated_at"`
	User_id       string    `json:"user_id" bson:"user_id"`
}

// GitHubUser represents user information from GitHub API
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubEmail represents email information from GitHub API
type GitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}
