package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	parcelPb "github.com/sachinggsingh/PDTS-Go/pb/parcel"
	"github.com/sachinggsingh/PDTS/api-gateway/internal/clients"
)

type ParcelHandler struct {
	client clients.ValidationClient
}

func NewParcelHandler(client clients.ValidationClient) *ParcelHandler {
	return &ParcelHandler{client: client}
}

func (h *ParcelHandler) CreateParcel(c *gin.Context) {
	var req struct {
		UserId string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.ParcelClient.CreateParcel(ctx, &parcelPb.CreateParcelRequest{
		UserId: req.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ParcelHandler) UpdateParcel(c *gin.Context) {
	parcelID := c.Param("id")

	var req struct {
		Status  string `json:"status"`
		Address *struct {
			AddressLine1 string `json:"address_line_1"`
			AddressLine2 string `json:"address_line_2"`
			City         string `json:"city"`
			State        string `json:"state"`
			Country      string `json:"country"`
			PostalCode   string `json:"postal_code"`
		} `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map request to Proto
	protoReq := &parcelPb.UpdateParcelRequest{
		ParcelId: parcelID,
	}

	if req.Status != "" {
		// Simple mapping, might need more robust validation
		switch req.Status {
		case "CREATED":
			protoReq.Status = parcelPb.ParcelStatus_PARCEL_STATUS_CREATED
		case "PICKED_UP":
			protoReq.Status = parcelPb.ParcelStatus_PARCEL_STATUS_PICKED_UP
		case "IN_TRANSIT":
			protoReq.Status = parcelPb.ParcelStatus_PARCEL_STATUS_IN_TRANSIT
		case "DELIVERED":
			protoReq.Status = parcelPb.ParcelStatus_PARCEL_STATUS_DELIVERED
		case "CANCELLED":
			protoReq.Status = parcelPb.ParcelStatus_PARCEL_STATUS_CANCELLED
		}
	}

	if req.Address != nil {
		protoReq.Address = &parcelPb.Address{
			AddressLine_1: req.Address.AddressLine1,
			AddressLine_2: req.Address.AddressLine2,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			PostalCode:    req.Address.PostalCode,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.ParcelClient.UpdateParcel(ctx, protoReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
