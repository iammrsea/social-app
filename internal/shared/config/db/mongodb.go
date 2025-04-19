package db

import (
	"context"
	"time"

	"github.com/iammrsea/social-app/internal/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig holds configuration for MongoDB connection
type MongoConfig struct {
	URI          string
	DatabaseName string
	Timeout      time.Duration
}

// NewDefaultMongoConfig creates a default MongoDB configuration
func NewMongoConfig() MongoConfig {
	uri := config.Env().MongoDbURI()
	dbName := config.Env().MongoDbName()

	return MongoConfig{
		URI:          uri,
		DatabaseName: dbName,
		Timeout:      10 * time.Second,
	}
}

// Connect establishes a connection to MongoDB and returns the client and database
func Connect(ctx context.Context, config MongoConfig) (*mongo.Client, *mongo.Database, error) {
	// Create a context with timeout for the connection
	connectionCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	// Set client options
	clientOpts := options.Client().
		ApplyURI(config.URI).
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
	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, nil, err
	}

	// Get the database
	db := client.Database(config.DatabaseName)
	return client, db, nil
}

// Disconnect closes the MongoDB connection
func Disconnect(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		return nil
	}
	disconnectionCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return client.Disconnect(disconnectionCtx)
}
