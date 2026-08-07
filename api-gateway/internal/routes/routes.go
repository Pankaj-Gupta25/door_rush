package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/api-gateway/config"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/caching"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/clients"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/handlers"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/middleware"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, grpcClients *clients.ValidationClient, rateLimiter *caching.RateLimiter) {
	parcelHandler := handlers.NewParcelHandler(*grpcClients)
	trackingHandler := handlers.NewTrackingHandler(*grpcClients)

	// Apply Rate Limiting Middleware Globally
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP AND RUNNING"})
	})

	// Auth Service Routes (Public)
	// Example: POST /auth/user/create -> Auth Service /user/create
	authGroup := router.Group("/auth")
	{
		authGroup.Any("/*proxyPath", handlers.SimpleProxyHandler(cfg.AuthServiceURL))
	}

	// User Service Routes (Protected)
	userGroup := router.Group("/user")
	userGroup.Use(middleware.AuthMiddleware(cfg))
	{
		userGroup.Any("/*proxyPath", handlers.SimpleProxyHandler(cfg.UserServiceURL))
	}

	// Parcel Service Routes (Protected)
	parcelGroup := router.Group("/parcel")
	parcelGroup.Use(middleware.AuthMiddleware(cfg))
	{
		parcelGroup.Any("/*proxyPath", func(c *gin.Context) {
			proxyPath := c.Param("proxyPath")
			// Remove leading slash for cleaner matching if needed, or matched directly

			// POST /parcel/create
			if c.Request.Method == http.MethodPost && proxyPath == "/create" {
				parcelHandler.CreateParcel(c)
				return
			}

			// PUT /parcel/:id
			if c.Request.Method == http.MethodPut && proxyPath != "" && proxyPath != "/" {
				// Assume everything after / is the ID
				// Inject "id" param for the handler
				id := proxyPath[1:] // simple trim /
				c.Params = append(c.Params, gin.Param{Key: "id", Value: id})
				parcelHandler.UpdateParcel(c)
				return
			}

			// Fallback to proxy
			handlers.SimpleProxyHandler(cfg.ParcelServiceURL)(c)
		})
	}

	// Tracking Service Routes (Protected)
	trackingGroup := router.Group("/tracking")
	trackingGroup.Use(middleware.AuthMiddleware(cfg))
	{
		trackingGroup.Any("/*proxyPath", func(c *gin.Context) {
			proxyPath := c.Param("proxyPath")

			// GET /tracking/:id
			if c.Request.Method == http.MethodGet && proxyPath != "" && proxyPath != "/" {
				// Assume everything after / is the ID
				id := proxyPath[1:]
				c.Params = append(c.Params, gin.Param{Key: "id", Value: id})
				trackingHandler.GetTrackingDetails(c)
				return
			}

			// Fallback
			handlers.SimpleProxyHandler(cfg.TrackingServiceURL)(c)
		})
	}

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}
