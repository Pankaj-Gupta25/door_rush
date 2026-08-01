package utils

import (
	"errors"

	"github.com/sachinggsingh/PDTS/parcel-service/internal/models"
)

func TotalAmount(itemAmount *float64, quantity *int) float64 {
	if itemAmount == nil || quantity == nil {
		return 0
	}
	if *itemAmount < 0 || *quantity < 0 {
		return 0
	}
	return *itemAmount * float64(*quantity)
}

func validateAddress(address *models.Address, owner string) error {
	if address == nil {
		return errors.New(owner + " address is required")
	}

	if address.AddressLine1 == nil || *address.AddressLine1 == "" {
		return errors.New(owner + " address line 1 is required")
	}
	if address.City == nil || *address.City == "" {
		return errors.New(owner + " city is required")
	}
	if address.State == nil || *address.State == "" {
		return errors.New(owner + " state is required")
	}
	if address.Country == nil || *address.Country == "" {
		return errors.New(owner + " country is required")
	}
	if address.PostalCode == nil || *address.PostalCode == "" {
		return errors.New(owner + " postal code is required")
	}

	return nil
}

func ValidateSenderDetails(sender *models.SenderDetails) error {
	if sender == nil {
		return errors.New("sender details are required")
	}
	if sender.Name == nil || *sender.Name == "" {
		return errors.New("sender name is required")
	}
	if sender.Phone == nil || *sender.Phone == "" {
		return errors.New("sender phone is required")
	}

	return validateAddress(sender.Address, "sender")
}

func ValidateReceiverDetails(receiver *models.ReceiverDetails) error {
	if receiver == nil {
		return errors.New("receiver details are required")
	}
	if receiver.Name == nil || *receiver.Name == "" {
		return errors.New("receiver name is required")
	}
	if receiver.Phone == nil || *receiver.Phone == "" {
		return errors.New("receiver phone is required")
	}

	return validateAddress(receiver.Address, "receiver")
}

func ValidateParcel(parcel *models.Parcel) error {
	if parcel == nil {
		return errors.New("parcel is nil")
	}

	if parcel.ItemName == nil || *parcel.ItemName == "" {
		return errors.New("item name is required")
	}

	if parcel.Quantity == nil || *parcel.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if parcel.Weight == nil || *parcel.Weight <= 0 {
		return errors.New("weight must be greater than zero")
	}

	// if parcel.Destination == nil || *parcel.Destination == "" {
	// 	return errors.New("destination is required")
	// }

	if parcel.SenderID == "" {
		return errors.New("sender id is required")
	}
	if parcel.ReceiverID == "" {
		return errors.New("receiver id is required")
	}

	if err := ValidateSenderDetails(parcel.SenderDetails); err != nil {
		return err
	}
	if err := ValidateReceiverDetails(parcel.ReceiverDetails); err != nil {
		return err
	}

	return nil
}

func ErrInvalidUserId() error {
	return errors.New("invalid user id")
}
