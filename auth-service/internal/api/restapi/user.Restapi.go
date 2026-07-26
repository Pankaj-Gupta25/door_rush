package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/auth-service/internal/models"
	"github.com/sachinggsingh/PDTS/auth-service/internal/service"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
)

type RestAPI struct {
	Logger      *utils.Logger
	Router      *gin.Engine
	UserService *service.UserService
	jwtManager  *utils.JWTManager
}

func NewRestAPI(logger *utils.Logger, router *gin.Engine, userService *service.UserService, jwtManager *utils.JWTManager) *RestAPI {
	return &RestAPI{
		Logger:      logger,
		Router:      router,
		UserService: userService,
		jwtManager:  jwtManager,
	}
}

func (s *RestAPI) SetupRoutes() {
	s.Router.POST("/user/create", s.CreateUserHandler)
	s.Router.POST("/user/signin", s.SignInUserHandler)
	s.Router.GET("/user/profile", s.ProfileOfTheUserHandler)
}

// CreateUserHandler handles user creation requests
func (s *RestAPI) CreateUserHandler(c *gin.Context) {
	var user models.User

	// Bind JSON request body to user model
	if err := c.ShouldBindJSON(&user); err != nil {
		s.Logger.Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service method
	createdUser, err := s.UserService.CreateUser(&user)
	if err != nil {
		s.Logger.Error("Failed to create user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, createdUser)
}

// SignInUserHandler handles user sign-in requests
func (s *RestAPI) SignInUserHandler(c *gin.Context) {
	var user models.User

	// Bind JSON request body to user model
	if err := c.ShouldBindJSON(&user); err != nil {
		s.Logger.Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service method
	signedInUser, err := s.UserService.SignInUser(&user)
	if err != nil {
		s.Logger.Error("Failed to sign in user: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, signedInUser)
}

// ProfileOfTheUserHandler handles user profile retrieval requests
func (s *RestAPI) ProfileOfTheUserHandler(c *gin.Context) {
	// Extract email from JWT token in Authorization header
	email, err := s.jwtManager.GetEmailFromToken(c)
	if err != nil {
		s.Logger.Warn("Failed to extract email from token: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Call service method
	userProfile, err := s.UserService.ProfileOfTheUser(email)
	if err != nil {
		s.Logger.Error("Failed to get user profile: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, userProfile)
}
