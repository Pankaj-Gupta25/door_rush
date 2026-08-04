package db

import (
	"context"
	"time"

	"github.com/sachinggsingh/PDTS-Go/tracking-service/config"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DB                 *mongo.Database
	TrackingCollection *mongo.Collection
	Logger             *utils.Logger
	Ctx                context.Context
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

	env := config.LoadENV()
	opts := options.Client().ApplyURI(env.MONGODB_URI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		d.Logger.Error("Failed to connect to database: " + err.Error())
		return err
	}

	d.Logger.Info("Connected to Tracking-service Database")
	d.DB = client.Database(env.MONGODB_DATABASE)
	d.TrackingCollection = d.DB.Collection("tracking")
	return nil
}

func (d *Database) DisconnectFromDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.Logger.Info("Disconnected from Tracking-service Database")
	return d.DB.Client().Disconnect(ctx)
}
