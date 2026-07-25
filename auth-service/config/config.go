package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	PORT             string
	MONGODB_URI      string
	MONGODB_DATABASE string
	JWT_SECRET       string

	// GitHub OAuth
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
	GITHUB_REDIRECT_URL  string
	OAUTH_STATE_SECRET   string
}

func GetEnv() *ENV {
	// Try to load .env from current directory first
	_ = godotenv.Load()
	// Ignore error as .env file is optional (e.g., in Docker)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("Error the Port is not set")
	}
	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		log.Fatalf("Error the MONGODB_URI is not set")
	}
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatalf("Error the JWT_SECRET is not set")
	}
	mongodb_database := os.Getenv("MONGODB_DATABASE")
	if mongodb_database == "" {
		log.Fatalf("Error the MONGODB_DATABASE is not set")
	}

	// GitHub OAuth (optional for development)
	github_client_id := os.Getenv("GITHUB_CLIENT_ID")
	github_client_secret := os.Getenv("GITHUB_CLIENT_SECRET")
	github_redirect_url := os.Getenv("GITHUB_REDIRECT_URL")
	oauth_state_secret := os.Getenv("OAUTH_STATE_SECRET")

	return &ENV{
		PORT:                 port,
		MONGODB_URI:          mongodb_uri,
		MONGODB_DATABASE:     mongodb_database,
		JWT_SECRET:           jwt_secret,
		GITHUB_CLIENT_ID:     github_client_id,
		GITHUB_CLIENT_SECRET: github_client_secret,
		GITHUB_REDIRECT_URL:  github_redirect_url,
		OAUTH_STATE_SECRET:   oauth_state_secret,
	}
}
