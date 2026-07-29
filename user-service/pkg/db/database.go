package db

import (
	"context"
	"time"

	"github.com/sachinggsingh/PDTS/user-service/config"
	"github.com/sachinggsingh/PDTS/user-service/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	DB                    *mongo.Database
	Client                *mongo.Client
	AddressCollection     *mongo.Collection
	UserProfileCollection *mongo.Collection
	Logger                *utils.Logger
	Ctx                   context.Context
}

func NewDatabase() *Database {
	return &Database{
		Logger: utils.Log,
		Ctx:    context.Background(),
	}
}

func (d *Database) Connect() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.GetEnv().MONGODB_URI).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	d.Logger.Info("You successfully connected to MongoDB!")

	// Assign the client and database to the struct
	d.Client = client
	d.DB = client.Database(config.GetEnv().MONGODB_DATABASE)
	d.AddressCollection = d.DB.Collection("addresses")
	d.UserProfileCollection = d.DB.Collection("user_profiles")
	return nil
}
func (d *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}
