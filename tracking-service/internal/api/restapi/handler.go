package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/middleware"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/models"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/service"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/websockets"
)

type TrackingHandler struct {
	service service.TrackingService
	hub     *websockets.Hub
	logger  *utils.Logger
}

func NewTrackingHandler(service service.TrackingService, hub *websockets.Hub, logger *utils.Logger) *TrackingHandler {
	return &TrackingHandler{
		service: service,
		hub:     hub,
		logger:  logger,
	}
}

func (h *TrackingHandler) SetupRoutes(r *gin.Engine, secretKey string) {
	tracking := r.Group("/tracking")
	tracking.Use(middleware.Authorize(secretKey))
	{
		tracking.POST("/create", h.CreateTracking)
		tracking.GET("/:parcel_id", h.GetTracking)
		tracking.PUT("/update/:parcel_id", h.UpdateStatus)
		tracking.GET("/ws/:parcel_id", h.ServeWs)
	}
}

func (h *TrackingHandler) CreateTracking(c *gin.Context) {
	var tracking models.Tracking
	if err := c.ShouldBindJSON(&tracking); err != nil {
		h.logger.Error("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	created, err := h.service.CreateTracking(&tracking)
	if err != nil {
		h.logger.Error("Failed to create tracking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tracking"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *TrackingHandler) GetTracking(c *gin.Context) {
	parcelID := c.Param("parcel_id")
	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parcel_id is required"})
		return
	}

	tracking, err := h.service.GetTracking(parcelID)
	if err != nil {
		h.logger.Error("Failed to get tracking: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "tracking not found"})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

func (h *TrackingHandler) UpdateStatus(c *gin.Context) {
	parcelID := c.Param("parcel_id")
	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parcel_id is required"})
		return
	}

	var req struct {
		Status   models.CurrentStatus `json:"status" binding:"required"`
		Location string               `json:"location" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate Status
	switch req.Status {
	case models.ParcelPickedUp, models.ParcelInTransit, models.ParcelDelivered, models.ParcelCancelled:
		// Valid
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Allowed values: PICKED_UP, IN_TRANSIT, DELIVERED, CANCELLED"})
		return
	}

	updated, err := h.service.UpdateStatus(parcelID, req.Status, req.Location)
	if err != nil {
		h.logger.Error("Failed to update status: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *TrackingHandler) ServeWs(c *gin.Context) {
	parcelID := c.Param("parcel_id")
	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parcel_id is required"})
		return
	}

	conn, err := websockets.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to websocket: " + err.Error())
		return
	}

	client := &websockets.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		ParcelID: parcelID,
	}
	h.hub.Register <- client

	// Start pump goroutines
	go client.WritePump()
	go client.ReadPump()
}
