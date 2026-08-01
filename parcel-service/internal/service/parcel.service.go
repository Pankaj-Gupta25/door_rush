package service

import (
	"context"
	"time"

	"github.com/sachinggsingh/PDTS/parcel-service/internal/messaging"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/repository"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
)

// Service Interface
type ParcelService interface {
	CreateParcel(parcel *models.Parcel) (*models.Parcel, error)
	GetParcelsByUserId(userId string) ([]models.Parcel, error)
	UpdateParcelByUserId(userId, parcelID string, parcel *models.Parcel) (*models.Parcel, error)
	DeleteParcelByUserId(userId, parcelID string) error
	GetParcelMeta(parcelID string) (models.ParcelMeta, bool)
	GetParcelStatusHistory(parcelID string) ([]models.StatusNode, error)
}

type parcelService struct {
	parcelRepo    repository.ParcelRepository
	parcelStore   repository.ParcelStore
	kafkaProducer messaging.KafkaProducer
	logger        *utils.Logger
}

// Constructor
func NewParcelService(
	parcelRepo repository.ParcelRepository,
	parcelStore repository.ParcelStore,
	kafkaProducer messaging.KafkaProducer,
	logger *utils.Logger,
) ParcelService {
	return &parcelService{
		parcelRepo:    parcelRepo,
		parcelStore:   parcelStore,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

/*
Service methods
*/
func (ps *parcelService) CreateParcel(parcel *models.Parcel) (*models.Parcel, error) {
	if err := utils.ValidateParcel(parcel); err != nil {
		return nil, err
	}
	status := models.ParcelCreated
	parcel.Status = &status

	createdParcel, err := ps.parcelRepo.CreateParcel(parcel)
	if err != nil {
		return nil, err
	}

	// Sync with in-memory store
	if createdParcel.Status != nil {
		status = *createdParcel.Status
	}

	ps.parcelStore.Add(models.ParcelMeta{
		ParcelID:   createdParcel.ParcelID,
		SenderID:   createdParcel.SenderID,
		ReceiverID: createdParcel.ReceiverID,
		Status:     status,
		UpdatedAt:  createdParcel.UpdatedAt,
	})

	// Publish to Kafka
	if ps.kafkaProducer != nil {
		go func() {
			err := ps.kafkaProducer.PublishParcelEvent(context.Background(), createdParcel.ParcelID, string(status))
			if err != nil {
				ps.logger.Error("Failed to publish parcel creation event: " + err.Error())
			}
		}()
	}

	return createdParcel, nil
}

func (ps *parcelService) GetParcelsByUserId(userId string) ([]models.Parcel, error) {
	if userId == "" {
		utils.Log.Warn("User ID not found")
		return nil, utils.ErrInvalidUserId()
	}

	return ps.parcelRepo.GetParcelsByUserId(userId)
}

func (ps *parcelService) UpdateParcelByUserId(userId, parcelID string, parcel *models.Parcel) (*models.Parcel, error) {
	if userId == "" {
		utils.Log.Warn("User ID not found")
		return nil, utils.ErrInvalidUserId()
	}

	// Get current parcel to check status change
	currentMeta, ok := ps.parcelStore.Get(parcelID)

	parcel.UpdatedAt = time.Now()
	updatedParcel, err := ps.parcelRepo.UpdateParcelByUserId(userId, parcelID, parcel)
	if err != nil {
		return nil, err
	}

	// Update in-memory store
	if !ok {
		// If not in store, add it
		status := models.ParcelCreated
		if updatedParcel.Status != nil {
			status = *updatedParcel.Status
		}
		ps.parcelStore.Add(models.ParcelMeta{
			ParcelID:   updatedParcel.ParcelID,
			SenderID:   updatedParcel.SenderID,
			ReceiverID: updatedParcel.ReceiverID,
			Status:     status,
			UpdatedAt:  updatedParcel.UpdatedAt,
		})
	} else if updatedParcel.Status != nil && currentMeta.Status != *updatedParcel.Status {
		// If status changed, use Update to append to history
		ps.parcelStore.Update(parcelID, *updatedParcel.Status)
	} else {
		// Just sync other fields, Add will preserve history
		status := models.ParcelCreated
		if updatedParcel.Status != nil {
			status = *updatedParcel.Status
		}
		ps.parcelStore.Add(models.ParcelMeta{
			ParcelID:   updatedParcel.ParcelID,
			SenderID:   updatedParcel.SenderID,
			ReceiverID: updatedParcel.ReceiverID,
			Status:     status,
			UpdatedAt:  updatedParcel.UpdatedAt,
		})
	}

	return updatedParcel, nil
}

func (ps *parcelService) DeleteParcelByUserId(userId, parcelID string) error {
	if userId == "" {
		utils.Log.Warn("User ID not found")
		return utils.ErrInvalidUserId()
	}
	err := ps.parcelRepo.DeleteParcelByUserId(userId, parcelID)
	if err != nil {
		return err
	}

	// Remove from in-memory store
	ps.parcelStore.Delete(parcelID)

	return nil
}

func (ps *parcelService) GetParcelMeta(parcelID string) (models.ParcelMeta, bool) {
	return ps.parcelStore.Get(parcelID)
}

func (ps *parcelService) GetParcelStatusHistory(parcelID string) ([]models.StatusNode, error) {
	history, err := ps.parcelStore.GetHistory(parcelID)
	if err == nil {
		return history, nil
	}

	ps.logger.Debug("Parcel history not in store, attempting lazy load for: " + parcelID)

	// Lazy load from DB if not in store
	parcel, err := ps.parcelRepo.GetParcelByID(parcelID)
	if err != nil {
		ps.logger.Error("Failed to lazy load parcel from DB for ID " + parcelID + ": " + err.Error())
		return nil, err
	}

	ps.logger.Debug("Found parcel in DB with canonical ID: " + parcel.ParcelID)

	// Add to store to initialize history
	status := models.ParcelCreated
	if parcel.Status != nil {
		status = *parcel.Status
	}

	// Add with its canonical ID
	ps.parcelStore.Add(models.ParcelMeta{
		ParcelID:   parcel.ParcelID,
		SenderID:   parcel.SenderID,
		ReceiverID: parcel.ReceiverID,
		Status:     status,
		UpdatedAt:  parcel.UpdatedAt,
	})

	// If the requested ID was different, also add it under that key to avoid repeated lookups
	if parcel.ParcelID != parcelID {
		ps.logger.Debug("Aliasing requested ID " + parcelID + " to canonical ID " + parcel.ParcelID)
		ps.parcelStore.Add(models.ParcelMeta{
			ParcelID:   parcelID,
			SenderID:   parcel.SenderID,
			ReceiverID: parcel.ReceiverID,
			Status:     status,
			UpdatedAt:  parcel.UpdatedAt,
		})
	}

	return ps.parcelStore.GetHistory(parcel.ParcelID)
}
