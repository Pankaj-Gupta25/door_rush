package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	trackingPb "github.com/sachinggsingh/PDTS-Go/pb/tracking"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/clients"
)

type TrackingHandler struct {
	client clients.ValidationClient
}

func NewTrackingHandler(client clients.ValidationClient) *TrackingHandler {
	return &TrackingHandler{client: client}
}

func (h *TrackingHandler) GetTrackingDetails(c *gin.Context) {
	parcelID := c.Param("id")
	if parcelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.TrackingClient.GetTrackingDetails(ctx, &trackingPb.GetTrackingDetailsRequest{
		ParcelId: parcelID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
