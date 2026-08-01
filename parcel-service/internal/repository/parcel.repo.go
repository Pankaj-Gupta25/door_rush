package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
	"github.com/sachinggsingh/PDTS/parcel-service/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ParcelRepository interface {
	CreateParcel(parcel *models.Parcel) (*models.Parcel, error)
	GetParcelsByUserId(userId string) ([]models.Parcel, error)
	UpdateParcelByUserId(userId, parcelID string, parcel *models.Parcel) (*models.Parcel, error)
	DeleteParcelByUserId(userId, parcelID string) error
	GetParcelByID(parcelID string) (*models.Parcel, error)
}

type parcelRepository struct {
	db     *db.Database
	logger *utils.Logger
}

func NewParcelRepository(db *db.Database, logger *utils.Logger) ParcelRepository {
	return &parcelRepository{
		db:     db,
		logger: logger,
	}
}

func (pr *parcelRepository) CreateParcel(parcel *models.Parcel) (*models.Parcel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parcel.CreatedAt = time.Now()
	parcel.UpdatedAt = time.Now()
	// Ensure ParcelID and _id are consistent
	if parcel.ID.IsZero() {
		parcel.ID = primitive.NewObjectID()
	}
	parcel.ParcelID = parcel.ID.Hex()

	result, err := pr.db.ParcelCollection.InsertOne(ctx, parcel)
	if err != nil {
		pr.logger.Error(fmt.Sprintf("Error creating parcel: %v", err))
		return nil, err
	}

	parcel.ID = result.InsertedID.(primitive.ObjectID)
	return parcel, nil
}

func (pr *parcelRepository) GetParcelsByUserId(userId string) ([]models.Parcel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"sender_id": userId}

	cursor, err := pr.db.ParcelCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var parcels []models.Parcel
	if err := cursor.All(ctx, &parcels); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (pr *parcelRepository) UpdateParcelByUserId(
	userId string,
	parcelID string,
	parcel *models.Parcel,
) (*models.Parcel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateFields := bson.M{}

	// Basic fields
	if parcel.ItemName != nil {
		updateFields["item_name"] = parcel.ItemName
	}
	if parcel.Quantity != nil {
		updateFields["quantity"] = parcel.Quantity
	}
	if parcel.Weight != nil {
		updateFields["weight"] = parcel.Weight
	}
	if parcel.Dimensions != nil {
		updateFields["dimensions"] = parcel.Dimensions
	}
	if parcel.Amount != nil {
		updateFields["amount"] = parcel.Amount
	}
	if parcel.Status != nil {
		updateFields["status"] = parcel.Status
	}
	if parcel.ReceiverID != "" {
		updateFields["receiver_id"] = parcel.ReceiverID
	}

	// Sender details
	if parcel.SenderDetails != nil {
		if parcel.SenderDetails.Name != nil {
			updateFields["sender_details.name"] = parcel.SenderDetails.Name
		}
		if parcel.SenderDetails.Phone != nil {
			updateFields["sender_details.phone"] = parcel.SenderDetails.Phone
		}

		if parcel.SenderDetails.Address != nil {
			addr := parcel.SenderDetails.Address

			if addr.AddressLine1 != nil {
				updateFields["sender_details.address.address_line_1"] = addr.AddressLine1
			}
			if addr.AddressLine2 != nil && *addr.AddressLine2 != "" {
				updateFields["sender_details.address.address_line_2"] = addr.AddressLine2
			}
			if addr.City != nil {
				updateFields["sender_details.address.city"] = addr.City
			}
			if addr.State != nil {
				updateFields["sender_details.address.state"] = addr.State
			}
			if addr.Country != nil {
				updateFields["sender_details.address.country"] = addr.Country
			}
			if addr.PostalCode != nil {
				updateFields["sender_details.address.postal_code"] = addr.PostalCode
			}
		}
	}

	// Receiver details
	if parcel.ReceiverDetails != nil {
		if parcel.ReceiverDetails.Name != nil {
			updateFields["receiver_details.name"] = parcel.ReceiverDetails.Name
		}
		if parcel.ReceiverDetails.Phone != nil {
			updateFields["receiver_details.phone"] = parcel.ReceiverDetails.Phone
		}

		if parcel.ReceiverDetails.Address != nil {
			addr := parcel.ReceiverDetails.Address

			if addr.AddressLine1 != nil {
				updateFields["receiver_details.address.address_line_1"] = addr.AddressLine1
			}
			if addr.AddressLine2 != nil && *addr.AddressLine2 != "" {
				updateFields["receiver_details.address.address_line_2"] = addr.AddressLine2
			}
			if addr.City != nil {
				updateFields["receiver_details.address.city"] = addr.City
			}
			if addr.State != nil {
				updateFields["receiver_details.address.state"] = addr.State
			}
			if addr.Country != nil {
				updateFields["receiver_details.address.country"] = addr.Country
			}
			if addr.PostalCode != nil {
				updateFields["receiver_details.address.postal_code"] = addr.PostalCode
			}
		}
	}

	// update timestamp
	updateFields["updated_at"] = time.Now()

	filter := bson.M{
		"parcel_id": parcelID,
		"sender_id": userId,
	}

	update := bson.M{
		"$set": updateFields,
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updatedParcel models.Parcel
	err := pr.db.ParcelCollection.
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updatedParcel)

	if err != nil {
		pr.logger.Error(fmt.Sprintf("Error updating parcel: %v", err))
		return nil, err
	}

	return &updatedParcel, nil
}
func (pr *parcelRepository) DeleteParcelByUserId(userId string, parcelID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"parcel_id": parcelID,
		"sender_id": userId,
	}

	_, err := pr.db.ParcelCollection.DeleteOne(ctx, filter)
	return err
}

func (pr *parcelRepository) GetParcelByID(parcelID string) (*models.Parcel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var parcel models.Parcel
	// Try searching by parcel_id (string)
	filter := bson.M{
		"$or": []bson.M{
			{"parcel_id": parcelID},
		},
	}

	// Also try searching by _id (ObjectID) if it's a valid hex
	if objID, err := primitive.ObjectIDFromHex(parcelID); err == nil {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"_id": objID})
	}

	err := pr.db.ParcelCollection.FindOne(ctx, filter).Decode(&parcel)
	if err != nil {
		return nil, err
	}

	// For legacy records, ensure ParcelID is populated
	if parcel.ParcelID == "" {
		parcel.ParcelID = parcel.ID.Hex()
	}

	return &parcel, nil
}
