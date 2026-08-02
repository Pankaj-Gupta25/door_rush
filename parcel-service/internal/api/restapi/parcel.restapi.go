package restapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/middleware"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/service"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
)

type ParcelRestApi struct {
	parcelService service.ParcelService
	logger        *utils.Logger
	Router        *gin.Engine
}

func NewParcelHandler(parcelService service.ParcelService, logger *utils.Logger, Router *gin.Engine) *ParcelRestApi {
	return &ParcelRestApi{
		parcelService: parcelService,
		logger:        logger,
		Router:        Router,
	}
}

func (pr *ParcelRestApi) SetupRoutes(secretKey string) {
	parcelGroup := pr.Router.Group("/parcel")
	parcelGroup.Use(middleware.Authorize(secretKey))
	{
		parcelGroup.POST("/create", pr.CreateParcelHandler)
		parcelGroup.GET("/get", pr.GetParcelHandler)
		parcelGroup.PUT("/update/:id", pr.UpdateParcelHandler)
		parcelGroup.DELETE("/delete/:id", pr.DeleteParcelHandler)
	}
	pr.Router.GET("/parcel/history/:id", pr.GetParcelHistoryHandler)
}

func (pr *ParcelRestApi) CreateParcelHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		pr.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var parcel models.Parcel
	if err := c.ShouldBindJSON(&parcel); err != nil {
		pr.logger.Error("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	parcel.SenderID = userId

	createdParcel, err := pr.parcelService.CreateParcel(&parcel)
	if err != nil {
		pr.logger.Error("Failed to create parcel: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create parcel"})
		return
	}

	c.JSON(http.StatusCreated, createdParcel)
	pr.logger.Info("Parcel created successfully")
}

func (pr *ParcelRestApi) GetParcelHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		pr.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	parcels, err := pr.parcelService.GetParcelsByUserId(userId)
	if err != nil {
		pr.logger.Error("Failed to get parcels: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get parcels"})
		return
	}
	c.JSON(http.StatusOK, parcels)
}

func (pr *ParcelRestApi) UpdateParcelHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		pr.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	parcelID := c.Param("id")
	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parcel ID is required"})
		return
	}

	var parcel models.Parcel
	if err := c.ShouldBindJSON(&parcel); err != nil {
		pr.logger.Error("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	updatedParcel, err := pr.parcelService.UpdateParcelByUserId(userId, parcelID, &parcel)
	if err != nil {
		pr.logger.Error("Failed to update parcel: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update parcel"})
		return
	}
	pr.logger.Info("Parcel updated successfully for ID: " + parcelID)
	c.JSON(http.StatusOK, updatedParcel)
}

func (pr *ParcelRestApi) DeleteParcelHandler(c *gin.Context) {
	userId, ok := middleware.HasUserId(c)
	if !ok {
		pr.logger.Error("Failed to extract user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	parcelID := c.Param("id")
	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parcel ID is required"})
		return
	}

	var parcel models.Parcel
	if err := c.ShouldBindJSON(&parcel); err != nil {
		// If body is empty or invalid, still try to delete using URL ID
		pr.logger.Warn("Invalid request body for delete, but proceeding with URL ID: " + err.Error())
	}

	err := pr.parcelService.DeleteParcelByUserId(userId, parcelID)
	if err != nil {
		pr.logger.Error("Failed to delete parcel: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete parcel"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "parcel deleted successfully"})
}

func (pr *ParcelRestApi) GetParcelHistoryHandler(c *gin.Context) {
	parcelID := c.Param("id")
	pr.logger.Info("Received history request for ID: " + parcelID)

	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parcel ID is required"})
		return
	}

	history, err := pr.parcelService.GetParcelStatusHistory(parcelID)
	if err != nil {
		pr.logger.Error("Failed to get parcel history for ID " + parcelID + ": " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "parcel history not found"})
		return
	}

	pr.logger.Info(fmt.Sprintf("Returning %d history nodes for ID %s", len(history), parcelID))
	c.JSON(http.StatusOK, history)
}
