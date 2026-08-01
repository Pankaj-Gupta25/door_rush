package models

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParcelStatus string

const (
	ParcelCreated   ParcelStatus = "CREATED"
	ParcelPickedUp  ParcelStatus = "PICKED_UP"
	ParcelInTransit ParcelStatus = "IN_TRANSIT"
	ParcelDelivered ParcelStatus = "DELIVERED"
	ParcelCancelled ParcelStatus = "CANCELLED"
)

type Parcel struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ParcelID string             `json:"parcel_id" bson:"parcel_id"`

	ItemName   *string  `json:"item_name,omitempty" bson:"item_name,omitempty"`
	Quantity   *int     `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Weight     *float64 `json:"weight,omitempty" bson:"weight,omitempty"`
	Dimensions *string  `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Amount     *float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	// Destination *string  `json:"destination,omitempty" bson:"destination,omitempty"`

	SenderID   string `json:"sender_id" bson:"sender_id"`
	ReceiverID string `json:"receiver_id" bson:"receiver_id"`

	SenderDetails   *SenderDetails   `json:"sender_details,omitempty" bson:"sender_details,omitempty"`
	ReceiverDetails *ReceiverDetails `json:"receiver_details,omitempty" bson:"receiver_details,omitempty"`

	Status    *ParcelStatus `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type ParcelMeta struct {
	ParcelID   string         `json:"parcel_id"`
	SenderID   string         `json:"sender_id"`
	ReceiverID string         `json:"receiver_id"`
	Status     ParcelStatus   `json:"status"`
	History    *StatusHistory `json:"history"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type SenderDetails struct {
	Name    *string  `json:"name,omitempty" bson:"name,omitempty"`
	Phone   *string  `json:"phone,omitempty" bson:"phone,omitempty"`
	Address *Address `json:"address,omitempty" bson:"address,omitempty"`
}

type ReceiverDetails struct {
	Name    *string  `json:"name,omitempty" bson:"name,omitempty"`
	Phone   *string  `json:"phone,omitempty" bson:"phone,omitempty"`
	Address *Address `json:"address,omitempty" bson:"address,omitempty"`
}

type Address struct {
	AddressLine1 *string `json:"address_line_1,omitempty" bson:"address_line_1,omitempty"`
	AddressLine2 *string `json:"address_line_2,omitempty" bson:"address_line_2,omitempty"`
	City         *string `json:"city,omitempty" bson:"city,omitempty"`
	State        *string `json:"state,omitempty" bson:"state,omitempty"`
	Country      *string `json:"country,omitempty" bson:"country,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
}

type StatusNode struct {
	Status    ParcelStatus
	TimeStamp time.Time
	Next      *StatusNode
}

type StatusHistory struct {
	mu   sync.RWMutex
	Head *StatusNode
	Tail *StatusNode
}

func (sh *StatusHistory) Append(status ParcelStatus) {
	// Lock the mutex to ensure thread safety
	sh.mu.Lock()
	// Unlock the mutex when done
	defer sh.mu.Unlock()

	// Create a new status node
	newNode := &StatusNode{
		Status:    status,
		TimeStamp: time.Now(),
	}

	// If the list is empty, set the head and tail to the new node
	if sh.Head == nil {
		sh.Head = newNode
		sh.Tail = newNode
		return
	}

	// Add the new node to the end of the list
	sh.Tail.Next = newNode
	sh.Tail = newNode
}

func (sh *StatusHistory) GetAll() []StatusNode {
	// Lock the mutex to ensure thread safety
	sh.mu.RLock()
	// Unlock the mutex when done
	defer sh.mu.RUnlock()

	// Create a slice to store the history
	var history []StatusNode
	current := sh.Head
	// Iterate through the list and append each node to the history slice
	for current != nil {
		history = append(history, *current)
		current = current.Next
	}
	// Return the history slice
	return history
}

// Note : working with the Mutex because we are using the multiple go routines to update the status of the parcel and we want to ensure that the status is updated correctly taking the scenerio of the real world where multiple users are updating the status of the parcel so it won't lead to the data corruption and data loss

// Mutex : provide a way to synchronize access to shared resources in a concurrent environment

// RWMutex : provide a way to synchronize access to shared resources in a concurrent environment
