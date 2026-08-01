package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
)

type ENV struct {
	PORT             string
	MONGODB_URI      string
	MONGODB_DATABASE string
	JWT_SECRET       string
	KAFKA_BROKERS    string
}

func LoadENV() *ENV {
	_ = godotenv.Load()
	// Ignore error as .env file is optional (e.g., in Docker)
	port := os.Getenv("PORT")
	if port == "" {
		utils.Log.Error("PORT is not set in the .env file")
		log.Fatal("PORT is not set in the .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		utils.Log.Error("MONGODB_URI is not set")
		log.Fatal("MONGODB_URI is not set")
	}

	mongoDB := os.Getenv("MONGODB_DATABASE")
	if mongoDB == "" {
		utils.Log.Error("MONGODB_DATABASE is not set")
		log.Fatal("MONGODB_DATABASE is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		utils.Log.Error("JWT_SECRET is not set")
		log.Fatal("JWT_SECRET is not set")
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		utils.Log.Warn("KAFKA_BROKERS is not set")
	}

	return &ENV{
		PORT:             port,
		MONGODB_URI:      mongoURI,
		MONGODB_DATABASE: mongoDB,
		JWT_SECRET:       jwtSecret,
		KAFKA_BROKERS:    kafkaBrokers,
	}
}
