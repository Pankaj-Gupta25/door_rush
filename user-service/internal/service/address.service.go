package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/sachinggsingh/PDTS/user-service/internal/models"
	"github.com/sachinggsingh/PDTS/user-service/internal/repository"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressService struct {
	addressRepo repository.AddressRepo
	logger      *utils.Logger
}

func NewAddressService(addressRepo repository.AddressRepo, logger *utils.Logger) *AddressService {
	return &AddressService{addressRepo: addressRepo, logger: logger}
}

func (ad *AddressService) CreateAddress(userID string, address *models.Address) (*models.Address, error) {
	// Check if address already exists for this user
	addressExists, err := ad.addressRepo.ChecksIfTheAddressExistsForTheUser(userID)
	if err != nil {
		ad.logger.Error(fmt.Sprintf("Error checking if address exists: %v", err))
		return nil, err
	}
	if addressExists {
		ad.logger.Error(fmt.Sprintf("Address already exists for user: %s", userID))
		return nil, errors.New("address already exists for this user")
	}

	// Set userID from context (not from request body)
	address.UserID = userID
	now := time.Now()
	address.CreatedAt = now
	address.UpdatedAt = now
	address.AddressID = primitive.NewObjectID().Hex()

	// Save the address to the database
	createdAddress, err := ad.addressRepo.CreateAddress(*address)
	if err != nil {
		ad.logger.Error(fmt.Sprintf("Failed to create address in database: %v", err))
		return nil, err
	}

	return &createdAddress, nil
}

func (ad *AddressService) GetAddressByUserId(userId string) (*models.Address, error) {
	address, err := ad.addressRepo.GetAddressByUserID(userId)
	if err != nil {
		ad.logger.Error(fmt.Sprintf("Error getting address by user ID: %v", err))
		return nil, err
	}
	return &address, nil
}

func (ad *AddressService) UpdateAddressByUserId(userId string, address *models.Address) (*models.Address, error) {

	// Set userID from context (not from request body)
	address.UserID = userId
	address.UpdatedAt = time.Now()

	updatedAddress, err := ad.addressRepo.UpdateAddressByUserId(userId, *address)
	if err != nil {
		ad.logger.Error(fmt.Sprintf("Error updating address by user ID: %v", err))
		return nil, err
	}
	return &updatedAddress, nil
}

func (ad *AddressService) DeleteAddressByUserId(userId string) error {
	_, err := ad.addressRepo.DeleteAddressByUserId(userId, models.Address{})
	if err != nil {
		ad.logger.Error(fmt.Sprintf("Error deleting address by user ID: %v", err))
		return err
	}
	return nil
}
