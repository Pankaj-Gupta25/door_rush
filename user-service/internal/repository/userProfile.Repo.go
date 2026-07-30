package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sachinggsingh/PDTS/user-service/internal/models"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
	"github.com/sachinggsingh/PDTS/user-service/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProfile interface {
	CreateUserProfile(models.UserProfile) (models.UserProfile, error)
	UpdateUserProfile(models.UserProfile) (models.UserProfile, error)
	DeleteUserProfile(string, models.UserProfile) (models.UserProfile, error)
	GetUserProfileByUserID(string) (models.UserProfile, error)
}

type userProfileRepo struct {
	logger *utils.Logger
	db     *db.Database
}

func NewUserProfileRepo(logger *utils.Logger, db *db.Database) UserProfile {
	return &userProfileRepo{
		logger: logger,
		db:     db,
	}
}

func (up *userProfileRepo) CreateUserProfile(userProfile models.UserProfile) (models.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := up.db.UserProfileCollection.InsertOne(ctx, userProfile)

	userProfile.ID = primitive.NewObjectID()
	userProfile.CreatedAt = time.Now()
	userProfile.UpdatedAt = time.Now()
	userProfile.UserProfileID = primitive.NewObjectID().Hex()
	if err != nil {
		up.logger.Error(fmt.Sprintf("Error creating user profile: %v", err))
		return models.UserProfile{}, err
	}
	return userProfile, nil
}

func (up *userProfileRepo) UpdateUserProfile(userProfile models.UserProfile) (models.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_profile_id": userProfile.UserProfileID}
	update := bson.M{"$set": userProfile}

	userProfile.UpdatedAt = time.Now()

	_, err := up.db.UserProfileCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		up.logger.Error(fmt.Sprintf("Error updating user profile: %v", err))
		return models.UserProfile{}, err
	}
	return userProfile, nil
}

func (up *userProfileRepo) DeleteUserProfile(userID string, userProfile models.UserProfile) (models.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := up.db.UserProfileCollection.DeleteOne(ctx, bson.M{"user_id": userID, "user_profile_id": userProfile.UserProfileID})
	if err != nil {
		up.logger.Error(fmt.Sprintf("Error deleting user profile: %v", err))
		return models.UserProfile{}, err
	}
	return userProfile, nil
}

func (up *userProfileRepo) GetUserProfileByUserID(userID string) (models.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result models.UserProfile
	err := up.db.UserProfileCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		up.logger.Error(fmt.Sprintf("Error getting user profile by user ID: %v", err))
		return models.UserProfile{}, err
	}
	return result, nil
}
