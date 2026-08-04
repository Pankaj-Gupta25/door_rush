package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrentStatus string

const (
	ParcelPickedUp  CurrentStatus = "PICKED_UP"
	ParcelInTransit CurrentStatus = "IN_TRANSIT"
	ParcelDelivered CurrentStatus = "DELIVERED"
	ParcelCancelled CurrentStatus = "CANCELLED"
)

type Tracking struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TrackingID      string             `json:"tracking_id" bson:"tracking_id"`
	ParcelID        string             `json:"parcel_id" bson:"parcel_id"`
	Status          CurrentStatus      `json:"status" bson:"status"`
	CurrentLocation string             `json:"current_location" bson:"current_location"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	History         []StatusUpdate     `json:"history" bson:"history"`
}

type StatusUpdate struct {
	Status    CurrentStatus `json:"status"`
	Location  string        `json:"location"`
	Timestamp time.Time     `json:"timestamp"`
}

type ParcelLocation struct {
	ParcelID  string        `json:"parcel_id" bson:"parcel_id"`
	Latitude  float64       `json:"latitude" bson:"latitude"`
	Longitude float64       `json:"longitude" bson:"longitude"`
	Status    CurrentStatus `json:"status" bson:"status"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}
