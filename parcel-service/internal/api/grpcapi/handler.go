package grpcapi

import (
	"context"
	"fmt"

	pb "github.com/sachinggsingh/PDTS-Go/pb/parcel"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/service"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
)

type ParcelGrpcHandler struct {
	pb.UnimplementedParcelServiceServer
	service service.ParcelService
	logger  *utils.Logger
}

func NewParcelGrpcHandler(service service.ParcelService, logger *utils.Logger) *ParcelGrpcHandler {
	return &ParcelGrpcHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ParcelGrpcHandler) CreateParcel(ctx context.Context, req *pb.CreateParcelRequest) (*pb.CreateParcelResponse, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	h.logger.Info("Received CreateParcel gRPC request for user: " + req.UserId)

	// Create new parcel struct
	newParcel := &models.Parcel{
		SenderID: req.UserId,
	}

	createdParcel, err := h.service.CreateParcel(newParcel)
	if err != nil {
		h.logger.Error("Failed to create parcel via gRPC: " + err.Error())
		return nil, err
	}

	return &pb.CreateParcelResponse{
		ParcelId: createdParcel.ParcelID,
	}, nil
}

func (h *ParcelGrpcHandler) UpdateParcel(ctx context.Context, req *pb.UpdateParcelRequest) (*pb.UpdateParcelResponse, error) {
	h.logger.Info("Received UpdateParcel gRPC request for parcel: " + req.ParcelId)

	// Try to find the parcel meta to identify the user (Owner)
	meta, ok := h.service.GetParcelMeta(req.ParcelId)
	var userId string
	if ok {
		userId = meta.SenderID
	} else {
		// Use a system/unknown user ID if not found in cache.
		// Note: This relies on service allowing updates if we pass a valid userId,
		// but since we don't know it, we might get an error if the service checks strict ownership.
		// However, for this implementation, we will try to proceed or return error.
		h.logger.Warn("Parcel not found in cache for UpdateParcel gRPC. ID: " + req.ParcelId)
		// Assuming we can't proceed without UserID for now as UpdateParcelByUserId requires it.
		return nil, fmt.Errorf("parcel not found or not active")
	}

	updateModel := &models.Parcel{}

	if req.Status != pb.ParcelStatus_PARCEL_STATUS_UNSPECIFIED {
		statusStr := mapStatus(req.Status)
		updateModel.Status = &statusStr
	}

	if req.Address != nil {
		updateModel.ReceiverDetails = &models.ReceiverDetails{
			Address: &models.Address{
				AddressLine1: &req.Address.AddressLine_1,
				AddressLine2: &req.Address.AddressLine_2,
				City:         &req.Address.City,
				State:        &req.Address.State,
				Country:      &req.Address.Country,
				PostalCode:   &req.Address.PostalCode,
			},
		}
	}

	updated, err := h.service.UpdateParcelByUserId(userId, req.ParcelId, updateModel)
	if err != nil {
		h.logger.Error("Failed to update parcel via gRPC: " + err.Error())
		return nil, err
	}

	// Prepare response
	resp := &pb.UpdateParcelResponse{
		ParcelId: updated.ParcelID,
	}

	if updated.Status != nil {
		resp.Status = mapModelStatusToProto(*updated.Status)
	}

	if updated.ReceiverDetails != nil && updated.ReceiverDetails.Address != nil {
		addr := updated.ReceiverDetails.Address
		resp.Address = &pb.Address{}
		if addr.AddressLine1 != nil {
			resp.Address.AddressLine_1 = *addr.AddressLine1
		}
		if addr.AddressLine2 != nil {
			resp.Address.AddressLine_2 = *addr.AddressLine2
		}
		if addr.City != nil {
			resp.Address.City = *addr.City
		}
		if addr.State != nil {
			resp.Address.State = *addr.State
		}
		if addr.Country != nil {
			resp.Address.Country = *addr.Country
		}
		if addr.PostalCode != nil {
			resp.Address.PostalCode = *addr.PostalCode
		}
	}

	return resp, nil
}

func mapStatus(s pb.ParcelStatus) models.ParcelStatus {
	switch s {
	case pb.ParcelStatus_PARCEL_STATUS_CREATED:
		return models.ParcelCreated
	case pb.ParcelStatus_PARCEL_STATUS_PICKED_UP:
		return models.ParcelPickedUp
	case pb.ParcelStatus_PARCEL_STATUS_IN_TRANSIT:
		return models.ParcelInTransit
	case pb.ParcelStatus_PARCEL_STATUS_DELIVERED:
		return models.ParcelDelivered
	case pb.ParcelStatus_PARCEL_STATUS_CANCELLED:
		return models.ParcelCancelled
	default:
		return models.ParcelCreated
	}
}

func mapModelStatusToProto(s models.ParcelStatus) pb.ParcelStatus {
	switch s {
	case models.ParcelCreated:
		return pb.ParcelStatus_PARCEL_STATUS_CREATED
	case models.ParcelPickedUp:
		return pb.ParcelStatus_PARCEL_STATUS_PICKED_UP
	case models.ParcelInTransit:
		return pb.ParcelStatus_PARCEL_STATUS_IN_TRANSIT
	case models.ParcelDelivered:
		return pb.ParcelStatus_PARCEL_STATUS_DELIVERED
	case models.ParcelCancelled:
		return pb.ParcelStatus_PARCEL_STATUS_CANCELLED
	default:
		return pb.ParcelStatus_PARCEL_STATUS_UNSPECIFIED
	}
}
