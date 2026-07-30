package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sachinggsingh/PDTS/user-service/internal/models"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
	"github.com/sachinggsingh/PDTS/user-service/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddressRepo interface {
	CreateAddress(models.Address) (models.Address, error)
	GetAddresses(models.Address) (models.Address, error)
	UpdateAddressByUserId(string, models.Address) (models.Address, error)
	DeleteAddressByUserId(string, models.Address) (models.Address, error)
	GetAddressByUserID(string) (models.Address, error)
	ChecksIfTheAddressExistsForTheUser(string) (bool, error)
}

type addressRepo struct {
	logger *utils.Logger
	db     *db.Database
}

func NewAddressRepo(database *db.Database, logger *utils.Logger) AddressRepo {
	return &addressRepo{
		logger: logger,
		db:     database,
	}
}

func (ar *addressRepo) CreateAddress(address models.Address) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ar.db.AddressCollection.InsertOne(ctx, address)
	if err != nil {
		ar.logger.Error(fmt.Sprintf("Error creating address: %v", err))
		return models.Address{}, err
	}
	return address, nil
}

func (ar *addressRepo) UpdateAddressByUserId(userId string, address models.Address) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}
	// Update the address fields
	setFields := bson.M{
		"address_line_1": address.AddressLine1,
		"city":           address.City,
		"state":          address.State,
		"country":        address.Country,
		"pincode":        address.Pincode,
		"phone":          address.Phone,
		"updated_at":     time.Now(),
	}

	// Only include address_line_2 if it's not empty
	if address.AddressLine2 != "" {
		setFields["address_line_2"] = address.AddressLine2
	}

	update := bson.M{
		"$set": setFields,
	}

	result, err := ar.db.AddressCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		ar.logger.Error(fmt.Sprintf("Error updating address: %v", err))
		return models.Address{}, err
	}
	if result.MatchedCount == 0 {
		return models.Address{}, fmt.Errorf("address not found for user")
	}

	// Fetch the updated address
	updatedAddress, err := ar.GetAddressByUserID(userId)
	if err != nil {
		return models.Address{}, err
	}
	return updatedAddress, nil
}

func (ar *addressRepo) DeleteAddressByUserId(userId string, address models.Address) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := ar.db.AddressCollection.DeleteOne(ctx, bson.M{"user_id": userId})
	if err != nil {
		ar.logger.Error(fmt.Sprintf("Error deleting address: %v", err))
		return models.Address{}, err
	}
	return address, nil
}

func (ar *addressRepo) GetAddresses(address models.Address) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result models.Address
	err := ar.db.AddressCollection.FindOne(ctx, bson.M{"address_id": address.AddressID}).Decode(&result)
	if err != nil {
		ar.logger.Error(fmt.Sprintf("Error getting address: %v", err))
		return models.Address{}, err
	}
	return result, nil
}

func (ar *addressRepo) GetAddressByUserID(userID string) (models.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	var result models.Address
	err := ar.db.AddressCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Address{}, fmt.Errorf("address not found for user")
		}
		ar.logger.Error(fmt.Sprintf("Error getting address by user ID: %v", err))
		return models.Address{}, err
	}
	return result, nil
}

func (ar *addressRepo) ChecksIfTheAddressExistsForTheUser(userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the address exists for the user
	filter := bson.M{"user_id": userID}
	var result models.Address
	err := ar.db.AddressCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		ar.logger.Error(fmt.Sprintf("Error checking if address exists for user: %v", err))
		return false, err
	}
	return true, nil
}
