package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/api-gateway/config"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/caching"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/clients"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/routes"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Redis Client
	redisClient := config.NewRedisClient(cfg)
	rateLimiter := caching.NewRateLimiter(redisClient)

	// Initialize gRPC Clients
	// Note: We need to define the ports/URLs for gRPC which might be different from HTTP URLs in config.
	// For now assuming existing URLs are HTTP, we might need to append different ports or use new config vars.
	// We'll trust the user to provide the correct host:port in the environment variables or we append default ports.
	grpcClients, err := clients.InitGrpcClients(cfg.ParcelServiceURL, cfg.TrackingServiceURL)
	if err != nil {
		log.Printf("Failed to initialize gRPC clients: %v", err)
		// Proceeding without gRPC might not be desired, but for now we log.
	}

	// Initialize GinRouter
	r := gin.Default()

	// Setup Routes
	routes.SetupRoutes(r, cfg, grpcClients, rateLimiter)

	// Start Server
	log.Printf("Starting API Gateway on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
