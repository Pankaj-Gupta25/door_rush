package service

import (
	"errors"
	"time"

	"github.com/sachinggsingh/PDTS/auth-service/internal/models"
	"github.com/sachinggsingh/PDTS/auth-service/internal/repository"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	repo       repository.UserRepository
	logger     *utils.Logger
	jwtManager *utils.JWTManager
}

func NewUserService(repo repository.UserRepository, log *utils.Logger, jwtManager *utils.JWTManager) *UserService {
	return &UserService{repo: repo, logger: log, jwtManager: jwtManager}
}

func (u *UserService) CreateUser(user *models.User) (*models.User, error) {
	if user.Email == nil {
		u.logger.Info("Email is nil")
		return nil, errors.New("email is empty")
	}
	if user.Password == nil {
		u.logger.Info("Password is nil")
		return nil, errors.New("password is empty")
	}

	//check if the user.email exist of not
	existingEmail, err := u.repo.FindUserByEmailIfItExistOrNot(*user.Email)
	if err != nil {
		u.logger.Error("Error finding user")
		return nil, err
	}
	if existingEmail {
		u.logger.Info("Email already exists")
		return nil, errors.New("email already exists")
	}
	//check if the password is hashed
	hashedPassword, err := utils.HashPassword(*user.Password)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Info("Password is empty")
			return nil, errors.New("password is empty")
		}
		u.logger.Error("Error hashing password")
		return nil, err
	}
	user.Password = &hashedPassword

	now := time.Now()
	user.Created_at = now
	user.Updated_at = now
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	tokenString, refreshTokenString, err := u.jwtManager.GenerateToken(user.User_id, *user.Email, *user.Password, "")
	if err != nil {
		u.logger.Error("Error generating token")
		return nil, err
	}
	user.Token = &tokenString
	user.Refresh_Token = &refreshTokenString

	// create the user
	createUser, err := u.repo.CreateUser(user)
	if err != nil {
		u.logger.Error("Error creating user")
		return nil, err
	}
	return createUser, nil
}

func (u *UserService) SignInUser(user *models.User) (*models.User, error) {

	if user.Email == nil {
		u.logger.Warn("Email is nil")
		return nil, errors.New("email is empty")
	}
	if user.Password == nil {
		u.logger.Warn("Password is nil")
		return nil, errors.New("password is empty")
	}

	// Find the user by email
	storedUser, err := u.repo.FindUserByEmail(*user.Email)
	if err != nil {
		u.logger.Error("Error finding user")
		return nil, err
	}
	if storedUser == nil {
		u.logger.Info("User not found")
		return nil, errors.New("user not found")
	}

	// Check if the password is correct
	checkPassword, err := utils.CheckIfThePasswordIsCorrect(*user.Password, *storedUser.Password)
	if err != nil {
		u.logger.Error("Error checking password")
		return nil, err
	}
	if !checkPassword {
		u.logger.Info("Password is incorrect")
		return nil, errors.New("password is incorrect")
	}

	// Generate new tokens
	tokenString, refreshTokenString, err := u.jwtManager.GenerateToken(storedUser.User_id, *storedUser.Email, *storedUser.Password, "")
	if err != nil {
		u.logger.Error("Error generating token")
		return nil, err
	}
	storedUser.Token = &tokenString
	storedUser.Refresh_Token = &refreshTokenString

	// Update the user with new tokens
	updatedUser, err := u.repo.SignInUser(storedUser)
	if err != nil {
		u.logger.Error("Error updating user")
		return nil, err
	}
	return updatedUser, nil
}

func (u *UserService) ProfileOfTheUser(email string) (*models.User, error) {
	if email == "" {
		u.logger.Warn("Email is nil")
		return nil, errors.New("email is empty")
	}
	return u.repo.FindUserByEmail(email)
}
