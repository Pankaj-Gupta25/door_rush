package db

import (
	"context"
	"time"

	"github.com/sachinggsingh/PDTS/auth-service/config"
	"github.com/sachinggsingh/PDTS/auth-service/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	Client         *mongo.Client
	Db             *mongo.Database
	Ctx            context.Context
	logger         *utils.Logger
	UserCollection *mongo.Collection
}

func NewDb() *MongoClient {
	return &MongoClient{logger: utils.New()}
}

func (m *MongoClient) ConnectDb() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Fix: Apply ServerAPI options
	opts := options.Client().ApplyURI(config.GetEnv().MONGODB_URI).SetServerAPIOptions(serverAPI)

	// Best Practice: Use a context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	m.logger.Info(" You successfully connected to MongoDB!")

	// Assign the client and database to the struct
	m.Client = client
	m.Db = client.Database(config.GetEnv().MONGODB_DATABASE)
	m.UserCollection = m.Db.Collection("users")
	return nil
}

func (d *MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}
