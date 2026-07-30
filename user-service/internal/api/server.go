package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/user-service/config"
	"github.com/sachinggsingh/PDTS/user-service/internal/api/restapi"
	"github.com/sachinggsingh/PDTS/user-service/internal/repository"
	"github.com/sachinggsingh/PDTS/user-service/internal/service"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
	"github.com/sachinggsingh/PDTS/user-service/pkg/db"
)

type Server struct {
	env    *config.ENV
	logger *utils.Logger
	DB     *db.Database
}

func NewServer(env *config.ENV, logger *utils.Logger, DB *db.Database) *Server {
	return &Server{
		env:    env,
		logger: logger,
		DB:     DB,
	}
}

func (s *Server) APIServer() *gin.Engine {
	r := gin.New()

	// Add recovery middleware (handles panics)
	r.Use(gin.Recovery())

	// trusting only local host
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Add a simple health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok from user-service"})
	})

	// Initialize repository layer
	addressRepo := repository.NewAddressRepo(s.DB, s.logger)
	userProfileRepo := repository.NewUserProfileRepo(s.logger, s.DB)

	// Initialize service layer
	addressService := service.NewAddressService(addressRepo, s.logger)
	userProfileService := service.NewUserProfileService(s.logger, userProfileRepo)

	// Initialize REST API handlers
	addressHandler := restapi.NewAddressHandler(addressService, r, s.logger)
	profileHandler := restapi.NewProfileHandler(userProfileService, s.logger, r)

	// Setup routes
	addressHandler.SetupRoutes(s.env.JWT_SECRET)
	profileHandler.SetupRoutes(s.env.JWT_SECRET)

	return r
}
