package testutil

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once     sync.Once
	client   *mongo.Client
	mongoUri string
)

func SetupTestMongoDb(t *testing.T) *mongo.Client {
	t.Helper()
	var _client *mongo.Client
	once.Do(func() {
		uri, terminate := setupMongoDbContainer(context.Background(), t)
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			t.Fatalf("failed to connect to mongo: %v", err)
		}
		mongoUri = uri
		_client = client
		t.Cleanup(func() {
			_ = _client.Disconnect(context.Background())
			terminate()
			log.Default().Println("Cleaned up resources")
		})
	})
	client = _client
	return _client
}

func MongoURI() string {
	return mongoUri
}

func MongoDBClient() *mongo.Client {
	return client
}

func setupMongoDbContainer(ctx context.Context, t *testing.T) (string, func()) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:7",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	terminate := func() {
		_ = container.Terminate(ctx)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "27017/tcp")
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	return mongoURI, terminate
}
