package service

import (
	"errors"

	"github.com/sachinggsingh/PDTS/user-service/internal/models"
	"github.com/sachinggsingh/PDTS/user-service/internal/repository"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
)

type UserProfileService struct {
	logger   *utils.Logger
	userRepo repository.UserProfile
}

func NewUserProfileService(logger *utils.Logger, userRepo repository.UserProfile) *UserProfileService {
	return &UserProfileService{
		logger:   logger,
		userRepo: userRepo,
	}
}

func (us *UserProfileService) CreateUserProfile(userProfile models.UserProfile) (models.UserProfile, error) {
	if userProfile.UserID == "" {
		us.logger.Error("User ID is missing")
		return models.UserProfile{}, errors.New("user ID is missing")
	}

	return us.userRepo.CreateUserProfile(userProfile)
}

func (us *UserProfileService) GetUserProfileByUserID(userID string) (models.UserProfile, error) {
	if userID == "" {
		us.logger.Error("User ID is missing")
		return models.UserProfile{}, errors.New("user ID is missing")
	}
	return us.userRepo.GetUserProfileByUserID(userID)
}

func (us *UserProfileService) UpdateUserProfile(userProfile models.UserProfile) (models.UserProfile, error) {
	if userProfile.UserID == "" {
		us.logger.Error("User ID is missing")
		return models.UserProfile{}, errors.New("user ID is missing")
	}

	return us.userRepo.UpdateUserProfile(userProfile)
}

func (us *UserProfileService) DeleteUserProfile(userID string) (models.UserProfile, error) {
	if userID == "" {
		us.logger.Error("User ID is missing")
		return models.UserProfile{}, errors.New("user ID is missing")
	}
	return us.userRepo.DeleteUserProfile(userID, models.UserProfile{})
}
