package grpcapi

import (
	"context"
	"fmt"
	"time"

	pb "github.com/sachinggsingh/PDTS-Go/pb/tracking"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/service"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
)

type TrackingGrpcHandler struct {
	pb.UnimplementedTrackingServiceServer
	service service.TrackingService
	logger  *utils.Logger
}

func NewTrackingGrpcHandler(service service.TrackingService, logger *utils.Logger) *TrackingGrpcHandler {
	return &TrackingGrpcHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TrackingGrpcHandler) GetTrackingDetails(ctx context.Context, req *pb.GetTrackingDetailsRequest) (*pb.GetTrackingDetailsResponse, error) {
	if req.ParcelId == "" {
		return nil, fmt.Errorf("parcel_id is required")
	}

	h.logger.Info("Received GetTrackingDetails gRPC request for parcel: " + req.ParcelId)

	// Call service
	tracking, err := h.service.GetTracking(req.ParcelId)
	if err != nil {
		h.logger.Error("Failed to get tracking details via gRPC: " + err.Error())
		return nil, err
	}

	// Format history
	var historyStrs []string
	for _, update := range tracking.History {
		entry := fmt.Sprintf("%s at %s (%s)", update.Status, update.Location, update.Timestamp.Format(time.RFC3339))
		historyStrs = append(historyStrs, entry)
	}

	return &pb.GetTrackingDetailsResponse{
		ParcelId: tracking.ParcelID,
		Status:   string(tracking.Status),
		History:  historyStrs,
	}, nil
}
