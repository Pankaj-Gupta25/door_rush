package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/models"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrackingRepository interface {
	CreateTracking(tracking *models.Tracking) (*models.Tracking, error)
	UpsertTracking(tracking *models.Tracking) (*models.Tracking, error)
	GetTrackingByParcelId(parcelID string) (*models.Tracking, error)
	UpdateTrackingStatus(parcelID string, status models.CurrentStatus, location string) (*models.Tracking, error)
	SaveLocation(ctx context.Context, loc *models.ParcelLocation) error
}

type trackingRepository struct {
	db     *db.Database
	logger *utils.Logger
}

func NewTrackingRepository(db *db.Database, logger *utils.Logger) TrackingRepository {
	return &trackingRepository{
		db:     db,
		logger: logger,
	}
}

func (r *trackingRepository) CreateTracking(tracking *models.Tracking) (*models.Tracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tracking.TrackingID = primitive.NewObjectID().Hex()
	tracking.UpdatedAt = time.Now()

	_, err := r.db.TrackingCollection.InsertOne(ctx, tracking)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error creating tracking: %v", err))
		return nil, err
	}

	return tracking, nil
}

func (r *trackingRepository) GetTrackingByParcelId(parcelID string) (*models.Tracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tracking models.Tracking
	filter := bson.M{"parcel_id": parcelID}

	err := r.db.TrackingCollection.FindOne(ctx, filter).Decode(&tracking)
	if err != nil {
		return nil, err
	}

	return &tracking, nil
}

func (r *trackingRepository) UpdateTrackingStatus(
	parcelID string,
	status models.CurrentStatus,
	location string,
) (*models.Tracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":           status,
			"current_location": location,
			"updated_at":       time.Now(),
		},
		"$push": bson.M{
			"history": models.StatusUpdate{
				Status:    status,
				Location:  location,
				Timestamp: time.Now(),
			},
		},
	}

	filter := bson.M{"parcel_id": parcelID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedTracking models.Tracking
	err := r.db.TrackingCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedTracking)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error updating tracking status: %v", err))
		return nil, err
	}

	return &updatedTracking, nil
}

func (r *trackingRepository) UpsertTracking(tracking *models.Tracking) (*models.Tracking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If no ID is set, generate one for potential new insertion
	if tracking.TrackingID == "" {
		tracking.TrackingID = primitive.NewObjectID().Hex()
	}
	tracking.UpdatedAt = time.Now()

	filter := bson.M{"parcel_id": tracking.ParcelID}

	// Create update document
	update := bson.M{
		"$setOnInsert": bson.M{
			"tracking_id":      tracking.TrackingID,
			"parcel_id":        tracking.ParcelID,
			"status":           tracking.Status,
			"current_location": tracking.CurrentLocation,
			"history":          tracking.History,
			"created_at":       time.Now(),
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// If the tracking object has fields we want to ensure are updated even if exists, add them to $set.
	// However, for "Create" event, we usually just want to ensure it exists.
	// If we want to restart status on re-create, we would put status in $set.
	// Assuming "Create" event means "Start Tracking", so we preserve existing history if it exists.

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result models.Tracking
	err := r.db.TrackingCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error upserting tracking: %v", err))
		return nil, err
	}

	return &result, nil
}

func (r *trackingRepository) SaveLocation(ctx context.Context, loc *models.ParcelLocation) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	loc.UpdatedAt = time.Now()

	filter := bson.M{"parcel_id": loc.ParcelID}

	update := bson.M{
		"$set": bson.M{
			"latitude":   loc.Latitude,
			"longitude":  loc.Longitude,
			"status":     loc.Status,
			"updated_at": loc.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result models.ParcelLocation
	err := r.db.TrackingCollection.FindOneAndUpdate(ctx, filter, update, opts.SetUpsert(true)).Decode(&result)
	if err != nil {
		r.logger.Error(fmt.Sprintf("Error saving location: %v", err))
		return err
	}

	return nil
}
