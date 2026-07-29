package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Phone        string             `json:"phone" bson:"phone" binding:"required,min=10,max=15"`
	AddressLine1 string             `json:"address_line_1" bson:"address_line_1" binding:"required"`
	AddressLine2 string             `json:"address_line_2,omitempty" bson:"address_line_2,omitempty"`
	City         string             `json:"city" bson:"city" binding:"required"`
	State        string             `json:"state" bson:"state" binding:"required"`
	Country      string             `json:"country" bson:"country" binding:"required"`
	Pincode      string             `json:"pincode" bson:"pincode" binding:"required,len=6"`
	AddressID    string             `json:"address_id" bson:"address_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserProfile struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email         string             `json:"email" bson:"email" binding:"required,email"`
	DateOfBirth   *string            `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty" binding:"required"`
	Gender        *string            `json:"gender,omitempty" bson:"gender,omitempty" binding:"required"`
	Name          string             `json:"name" bson:"name" binding:"required"`
	Phone         string             `json:"phone" bson:"phone" binding:"required,min=10,max=15"`
	TotalOrders   int                `json:"total_orders" bson:"total_orders"`
	TotalSpent    float64            `json:"total_spent" bson:"total_spent"`
	UserID        string             `json:"user_id" bson:"user_id"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
	UserProfileID string             `json:"user_profile_id" bson:"user_profile_id"`
}
