package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
)

type ENV struct {
	PORT             string
	MONGODB_URI      string
	MONGODB_DATABASE string
	JWT_SECRET       string
}

func GetEnv() *ENV {

	_ = godotenv.Load()
	// Ignore error as .env file is optional (e.g., in Docker)
	port := os.Getenv("PORT")
	if port == "" {
		utils.Log.Error("PORT is not set")
		log.Fatalf("PORT is not set")
	}
	mongo_uri := os.Getenv("MONGODB_URI")
	if mongo_uri == "" {
		utils.Log.Error("MONGODB_URI is not set")
		log.Fatalf("MONGODB_URI is not set")
	}
	mongo_database := os.Getenv("MONGODB_DATABASE")
	if mongo_database == "" {
		utils.Log.Error("MONGODB_DATABASE is not set")
		log.Fatalf("MONGODB_DATABASE is not set")
	}
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		utils.Log.Error("JWT_SECRET is not set")
		log.Fatalf("JWT_SECRET is not set")
	}
	return &ENV{
		PORT:             port,
		MONGODB_URI:      mongo_uri,
		MONGODB_DATABASE: mongo_database,
		JWT_SECRET:       jwt_secret,
	}
}
