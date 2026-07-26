package repository

import (
	"context"
	"time"

	"github.com/sachinggsingh/PDTS/auth-service/internal/models"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
	"github.com/sachinggsingh/PDTS/auth-service/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	SignInUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByEmailIfItExistOrNot(email string) (bool, error)

	FindUserByOAuthId(oauthId string) (*models.User, error)
	CreateAndUpdateOAuthUser(user *models.User) (*models.User, error)
}

func NewUserRepository(database *db.MongoClient) UserRepository {
	return &userRepository{
		logger: utils.New(),
		db:     database,
	}
}

type userRepository struct {
	logger *utils.Logger
	db     *db.MongoClient
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.db.UserCollection.InsertOne(ctx, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Error creating user")
		return nil, err
	}
	return user, nil
}

func (r *userRepository) SignInUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": user.User_id,
	}
	update := bson.M{
		"$set": bson.M{
			"token":         user.Token,
			"refresh_token": user.Refresh_Token,
		},
	}
	_, err := r.db.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Error updating user")
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindUserByEmailIfItExistOrNot(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{
		"email": email,
	}
	var user models.User
	err := r.db.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		r.logger.Error("Error finding user")
		return false, err
	}
	return true, nil
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var findUserByEmail models.User

	err := r.db.UserCollection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&findUserByEmail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Error finding user")
		return nil, err
	}
	return &findUserByEmail, nil
}

func (r *userRepository) FindUserByOAuthId(oauthId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{
		"oauth_id": oauthId,
	}
	var OAuthUser models.User

	err := r.db.UserCollection.FindOne(ctx, filter).Decode(&OAuthUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Error finding user")
		return nil, err
	}
	return &OAuthUser, nil
}

func (r *userRepository) CreateAndUpdateOAuthUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.db.UserCollection.InsertOne(ctx, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		r.logger.Error("Error creating user")
		return nil, err
	}
	return user, nil
}
