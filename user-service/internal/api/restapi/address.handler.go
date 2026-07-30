package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS/user-service/internal/middleware"
	"github.com/sachinggsingh/PDTS/user-service/internal/models"
	"github.com/sachinggsingh/PDTS/user-service/internal/service"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
)

type AddressHandler struct {
	addressService *service.AddressService
	logger         *utils.Logger
	Router         *gin.Engine
}

func NewAddressHandler(addressService *service.AddressService, router *gin.Engine, logger *utils.Logger) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		logger:         logger,
		Router:         router,
	}
}

func (h *AddressHandler) SetupRoutes(secretKey string) {
	h.Router.POST("/address/create", h.CreateAddressHandler)
	h.Router.GET("/address", h.GetAddressHandler)
	h.Router.PUT("/address", h.UpdateAddressHandler)
	h.Router.DELETE("/address", h.DeleteAddressHandler)
}

// CreateAddressHandler handles address creation requests
func (h *AddressHandler) CreateAddressHandler(c *gin.Context) {
	// Get userID from middleware
	userID, ok := middleware.HasUserId(c)
	if !ok {
		// HasUserId already sends the error response and aborts
		h.logger.Error("Failed to extract user ID from token")
		return
	}

	var address models.Address

	// Bind JSON request body to address model
	if err := c.ShouldBindJSON(&address); err != nil {
		h.logger.Warn("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service method with userID from context
	createdAddress, err := h.addressService.CreateAddress(userID, &address)
	if err != nil {
		h.logger.Error("Failed to create address: " + err.Error())
		// Determine appropriate status code based on error
		statusCode := http.StatusInternalServerError
		if err.Error() == "address already exists for this user" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, createdAddress)
}

// GetAddressHandler handles address retrieval requests
func (h *AddressHandler) GetAddressHandler(c *gin.Context) {
	// Get userID from middleware
	userID, ok := middleware.HasUserId(c)
	if !ok {
		// HasUserId already sends the error response and aborts
		h.logger.Error("Failed to extract user ID from token")
		return
	}

	// Call service method
	address, err := h.addressService.GetAddressByUserId(userID)
	if err != nil {
		h.logger.Error("Failed to get address: " + err.Error())
		statusCode := http.StatusInternalServerError
		if err.Error() == "address not found for user" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, address)
}

// UpdateAddressHandler handles address update requests
func (h *AddressHandler) UpdateAddressHandler(c *gin.Context) {
	// Get userID from middleware
	userID, ok := middleware.HasUserId(c)
	if !ok {
		// HasUserId already sends the error response and aborts
		h.logger.Error("Failed to extract user ID from token")
		return
	}

	var address models.Address

	// Bind JSON request body to address model
	if err := c.ShouldBindJSON(&address); err != nil {
		h.logger.Warn("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service method with userID from context
	updatedAddress, err := h.addressService.UpdateAddressByUserId(userID, &address)
	if err != nil {
		h.logger.Error("Failed to update address: " + err.Error())
		statusCode := http.StatusInternalServerError
		if err.Error() == "address not found for user" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, updatedAddress)
}

func (h *AddressHandler) DeleteAddressHandler(c *gin.Context) {
	// Get userID from middleware
	userID, ok := middleware.HasUserId(c)
	if !ok {
		// HasUserId already sends the error response and aborts
		h.logger.Error("Failed to extract user ID from token")
		return
	}

	// Call service method
	err := h.addressService.DeleteAddressByUserId(userID)
	if err != nil {
		h.logger.Error("Failed to delete address: " + err.Error())
		statusCode := http.StatusInternalServerError
		if err.Error() == "address not found for user" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "address deleted successfully"})
}
