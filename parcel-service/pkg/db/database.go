package db

import (
	"context"
	"time"

	"github.com/sachinggsingh/PDTS/parcel-service/config"
	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DB               *mongo.Database
	ParcelCollection *mongo.Collection
	Logger           *utils.Logger
	Ctx              context.Context
}

func NewDatabase() *Database {
	return &Database{
		Logger: utils.Log,
		Ctx:    context.Background(),
	}
}

func (d *Database) ConnectToDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.LoadENV().MONGODB_URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		d.Logger.Error("Failed to connect to database: " + err.Error())
		return err
	}
	d.Logger.Info("Connected to Parcel-service Database")
	d.DB = client.Database(config.LoadENV().MONGODB_DATABASE)
	d.ParcelCollection = d.DB.Collection("parcels")
	return nil
}

func (d *Database) DisconnectFromDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.Logger.Info("Disconnected from Parcel-service Database")
	return d.DB.Client().Disconnect(ctx)
}
