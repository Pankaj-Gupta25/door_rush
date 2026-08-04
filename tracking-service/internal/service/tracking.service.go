package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/models"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/repository"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/websockets"
)

type TrackingService interface {
	CreateTracking(tracking *models.Tracking) (*models.Tracking, error)
	UpsertTracking(tracking *models.Tracking) (*models.Tracking, error)
	GetTracking(parcelID string) (*models.Tracking, error)
	UpdateStatus(parcelID string, status models.CurrentStatus, location string) (*models.Tracking, error)
	SaveLocation(ctx context.Context, loc *models.ParcelLocation) error
}

type trackingService struct {
	repo   repository.TrackingRepository
	hub    *websockets.Hub
	logger *utils.Logger
}

func NewTrackingService(repo repository.TrackingRepository, hub *websockets.Hub, logger *utils.Logger) TrackingService {
	return &trackingService{
		repo:   repo,
		hub:    hub,
		logger: logger,
	}
}

func (s *trackingService) CreateTracking(tracking *models.Tracking) (*models.Tracking, error) {
	if tracking.Status == "" {
		tracking.Status = models.ParcelPickedUp
	}
	tracking.History = append(tracking.History, models.StatusUpdate{
		Status:    tracking.Status,
		Location:  tracking.CurrentLocation,
		Timestamp: time.Now(),
	})

	return s.repo.CreateTracking(tracking)
}

func (s *trackingService) UpsertTracking(tracking *models.Tracking) (*models.Tracking, error) {
	if tracking.Status == "" {
		tracking.Status = models.ParcelPickedUp
	}
	// For upsert, we might want to be careful about appending history if it's already there.
	// But since this is called on "Create Parcel" event, we assume we want to ensure the initial state is recorded.
	// If the record exists, our Repo logic uses $setOnInsert for history, so it won't duplicate.
	if len(tracking.History) == 0 {
		tracking.History = append(tracking.History, models.StatusUpdate{
			Status:    tracking.Status,
			Location:  tracking.CurrentLocation,
			Timestamp: time.Now(),
		})
	}

	return s.repo.UpsertTracking(tracking)
}

func (s *trackingService) GetTracking(parcelID string) (*models.Tracking, error) {
	return s.repo.GetTrackingByParcelId(parcelID)
}

func (s *trackingService) UpdateStatus(
	parcelID string,
	status models.CurrentStatus,
	location string,
) (*models.Tracking, error) {
	updatedTracking, err := s.repo.UpdateTrackingStatus(parcelID, status, location)
	if err != nil {
		return nil, err
	}

	// Broadcast update
	msg, _ := json.Marshal(updatedTracking)
	s.hub.Broadcast <- msg

	return updatedTracking, nil
}

func (s *trackingService) SaveLocation(ctx context.Context, loc *models.ParcelLocation) error {
	err := s.repo.SaveLocation(ctx, loc)
	if err != nil {
		return err
	}

	// Broadcast update
	msg, _ := json.Marshal(loc)
	s.hub.Broadcast <- msg

	return nil
}
