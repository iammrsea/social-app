package db

import (
	"context"
	"log"
	"time"

	"github.com/iammrsea/social-app/internal/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// mongoConfig holds configuration for MongoDB connection
type mongoConfig struct {
	uri          string
	databaseName string
	timeout      time.Duration
}

// SetupMongoDB sets up and connects to MongoDB
func SetupMongoDB(ctx context.Context) (*mongo.Database, func() error) {
	uri := config.Env().MongoDbURI()
	dbName := config.Env().MongoDbName()
	mongoConfig := mongoConfig{
		uri:          uri,
		databaseName: dbName,
		timeout:      10 * time.Second,
	}
	client, db, err := connect(ctx, mongoConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	disconnectClient := func() error {
		return disconnect(ctx, client)
	}

	return db, disconnectClient
}

// connect establishes a connection to MongoDB and returns the client and database
func connect(ctx context.Context, config mongoConfig) (*mongo.Client, *mongo.Database, error) {
	// Create a context with timeout for the connection
	connectionCtx, cancel := context.WithTimeout(ctx, config.timeout)
	defer cancel()

	// Set client options
	clientOpts := options.Client().
		ApplyURI(config.uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Minute).
		SetRetryWrites(true).
		SetRetryReads(true)

	//Connect to MongoDB
	client, err := mongo.Connect(connectionCtx, clientOpts)
	if err != nil {
		return nil, nil, err
	}
	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, config.timeout)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, nil, err
	}

	// Get the database
	db := client.Database(config.databaseName)
	return client, db, nil
}

// disconnect closes the MongoDB connection
func disconnect(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		return nil
	}
	disconnectionCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return client.Disconnect(disconnectionCtx)
}
