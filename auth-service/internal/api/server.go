package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/auth-service/config"
	"github.com/sachinggsingh/PDTS/auth-service/internal/api/restapi"
	"github.com/sachinggsingh/PDTS/auth-service/internal/middleware"
	"github.com/sachinggsingh/PDTS/auth-service/internal/repository"
	"github.com/sachinggsingh/PDTS/auth-service/internal/service"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
	"github.com/sachinggsingh/PDTS/auth-service/pkg/db"
)

type Server struct {
	env    *config.ENV
	logger *utils.Logger
	DB     *db.MongoClient
}

func NewServer(env *config.ENV, logger *utils.Logger, DB *db.MongoClient) *Server {
	return &Server{
		env:    env,
		logger: logger,
		DB:     DB,
	}
}

func (s *Server) APIServer() *gin.Engine {

	// Create a new Gin engine without default middleware
	r := gin.New()

	// Add custom logger middleware
	r.Use(middleware.CustomLogger(s.logger))

	// Add recovery middleware (handles panics)
	r.Use(gin.Recovery())

	// trusting only local host
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Add a simple health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok from auth-service"})
	})

	// Initialize repository layer
	userRepo := repository.NewUserRepository(s.DB)

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(s.env.JWT_SECRET)

	// Initialize service layer
	userService := service.NewUserService(userRepo, s.logger, jwtManager)
	oauthService := service.NewOAuthService(userRepo, s.logger, jwtManager)

	// Initialize REST API handlers
	userAPI := restapi.NewRestAPI(s.logger, r, userService, jwtManager)
	oauthAPI := restapi.NewOAuthRestAPI(s.logger, r, oauthService)

	// Setup routes
	userAPI.SetupRoutes()
	oauthAPI.SetupOAuthRoutes()

	return r
}
