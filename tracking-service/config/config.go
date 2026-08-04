package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
)

type ENV struct {
	PORT             string
	MONGODB_URI      string
	MONGODB_DATABASE string
	JWT_SECRET       string
	KAFKA_BROKERS    string
}

func LoadENV() ENV {
	_ = godotenv.Load()
	// Ignore error as .env file is optional (e.g., in Docker)
	port := os.Getenv("PORT")
	if port == "" {
		utils.Log.Warn("PORT is not set")
	}
	mongodbUrl := os.Getenv("MONGODB_URI")
	if mongodbUrl == "" {
		utils.Log.Warn("MONGODB_URL is not set")
	}
	mongodbDatabase := os.Getenv("MONGODB_DATABASE")
	if mongodbDatabase == "" {
		utils.Log.Warn("MONGODB_DATABASE is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		utils.Log.Warn("JWT_SECRET is not set")
	}
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		utils.Log.Warn("KAFKA_BROKERS is not set")
	}
	return ENV{
		PORT:             port,
		MONGODB_URI:      mongodbUrl,
		MONGODB_DATABASE: mongodbDatabase,
		JWT_SECRET:       jwtSecret,
		KAFKA_BROKERS:    kafkaBrokers,
	}
}
