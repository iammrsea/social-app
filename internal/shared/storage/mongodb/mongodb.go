package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/iammrsea/social-app/internal/shared/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// SetupMongoDB sets up and connects to MongoDB
func SetupMongoDB(ctx context.Context, cf *config.MongoConfig) (*mongo.Database, func() error) {
	log.Default().Println("Connecting to MongoDB...")
	client, db, err := connect(ctx, cf)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Default().Println("✅ Connected to MongoDB")
	disconnectClient := func() error {
		return disconnect(ctx, client)
	}

	return db, disconnectClient
}

// connect establishes a connection to MongoDB and returns the client and database
func connect(ctx context.Context, cf *config.MongoConfig) (*mongo.Client, *mongo.Database, error) {
	// Create a context with timeout for the connection
	connectionCtx, cancel := context.WithTimeout(ctx, cf.Timeout)
	defer cancel()

	// Set client options
	clientOpts := options.Client().
		ApplyURI(cf.Uri).
		SetMaxPoolSize(uint64(cf.MaxPoolSize)).
		SetMinPoolSize(uint64(cf.MinPoolSize)).
		SetMaxConnIdleTime(cf.ConnIdleTime).
		SetRetryWrites(cf.RetryWrites).
		SetRetryReads(cf.RetryReads).
		SetReplicaSet(cf.ReplicaSet)

	//Connect to MongoDB
	client, err := mongo.Connect(connectionCtx, clientOpts)
	if err != nil {
		return nil, nil, err
	}
	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, cf.Timeout)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return nil, nil, err
	}
	// Get the database
	db := client.Database(cf.DatabaseName)
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
