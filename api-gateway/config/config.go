package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	AuthServiceURL     string
	UserServiceURL     string
	ParcelServiceURL   string
	TrackingServiceURL string
	JWTSecret          string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		log.Fatal("AUTH_SERVICE_URL environment variable is not set")
	}
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL environment variable is not set")
	}
	parcelServiceURL := os.Getenv("PARCEL_SERVICE_URL")
	if parcelServiceURL == "" {
		log.Fatal("PARCEL_SERVICE_URL environment variable is not set")
	}
	trackingServiceURL := os.Getenv("TRACKING_SERVICE_URL")
	if trackingServiceURL == "" {
		log.Fatal("TRACKING_SERVICE_URL environment variable is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Default to localhost if not set
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	// Password can be empty

	redisDB := 0
	// You could parse REDIS_DB from env if needed, keeping simple for now

	return &Config{
		Port:               port,
		AuthServiceURL:     authServiceURL,
		UserServiceURL:     userServiceURL,
		ParcelServiceURL:   parcelServiceURL,
		TrackingServiceURL: trackingServiceURL,
		JWTSecret:          jwtSecret,
		RedisAddr:          redisAddr,
		RedisPassword:      redisPassword,
		RedisDB:            redisDB,
	}
}
