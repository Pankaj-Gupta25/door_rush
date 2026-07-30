package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/user-service/internal/middleware"
	"github.com/sachinggsingh/PDTS/user-service/internal/models"
	"github.com/sachinggsingh/PDTS/user-service/internal/service"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
)

type ProfileHandler struct {
	profileService service.UserProfileService
	logger         *utils.Logger
	Router         *gin.Engine
}

func NewProfileHandler(profileHandler *service.UserProfileService, loggger *utils.Logger, Router *gin.Engine) *ProfileHandler {
	return &ProfileHandler{
		profileService: *profileHandler,
		logger:         loggger,
		Router:         Router,
	}
}

func (ph *ProfileHandler) SetupRoutes(secretKey string) {
	ph.Router.POST("/profile/create", ph.CreateUserProfileHandler)
	ph.Router.GET("/profile", ph.GetUserProfileHandler)
	ph.Router.PUT("/profile", ph.UpdateUserProfileHandler)
	ph.Router.DELETE("/profile", ph.DeleteUserProfileHandler)
	ph.Router.GET("/profile/:user_id", ph.GetProfileByTheUserId)
}

func (ph *ProfileHandler) CreateUserProfileHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		ph.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userProfile models.UserProfile
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		ph.logger.Error("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userProfile.UserID = userId
	createdUserProfile, err := ph.profileService.CreateUserProfile(userProfile)
	if err != nil {
		ph.logger.Error("Failed to create user profile: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user profile"})
		return
	}
	c.JSON(http.StatusCreated, createdUserProfile)
}

func (ph *ProfileHandler) GetUserProfileHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		ph.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userProfile, err := ph.profileService.GetUserProfileByUserID(userId)
	if err != nil {
		ph.logger.Error("Failed to get user profile: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user profile"})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}
func (ph *ProfileHandler) UpdateUserProfileHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		ph.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var userProfile models.UserProfile
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		ph.logger.Error("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	userProfile.UserID = userId
	updatedUserProfile, err := ph.profileService.UpdateUserProfile(userProfile)
	if err != nil {
		ph.logger.Error("Failed to update user profile: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user profile"})
		return
	}
	c.JSON(http.StatusOK, updatedUserProfile)
}
func (ph *ProfileHandler) DeleteUserProfileHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		ph.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userProfile, err := ph.profileService.DeleteUserProfile(userId)
	if err != nil {
		ph.logger.Error("Failed to delete user profile: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user profile"})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}

func (ph *ProfileHandler) GetProfileByTheUserId(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		ph.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userProfile, err := ph.profileService.GetUserProfileByUserID(userId)
	if err != nil {
		ph.logger.Error("Failed to get user profile: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user profile"})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}
