package config

import (
	"errors"
	"fmt"
	"time"
)

type StorageEngine string

var (
	mongodb    StorageEngine = "mongodb"
	postgreSQL StorageEngine = "postgresql"
	// inMemory   StorageEngine = "inmemory" // Not fully implemented yet
	// Add more storage engines here
)

var StorageEngines = struct {
	Mongodb    StorageEngine
	PostgreSQL StorageEngine
}{
	Mongodb:    mongodb,
	PostgreSQL: postgreSQL,
}

var (
	ErrUnsupportedEngine = errors.New("unsupported storage engine")
)

type StorageEngineConfig interface {
	Engine() StorageEngine
}

type MongoConfig struct {
	Uri          string
	DatabaseName string
	Timeout      time.Duration
	ReplicaSet   string
	MaxPoolSize  int32
	MinPoolSize  int32
	ConnIdleTime time.Duration
	RetryWrites  bool
	RetryReads   bool
}

func (m *MongoConfig) Engine() StorageEngine {
	return mongodb
}

type PostgresConfig struct {
	Uri          string
	MaxPoolSize  int32
	MinPoolSize  int32
	ConnIdleTime time.Duration
	Timeout      time.Duration
}

func (p *PostgresConfig) Engine() StorageEngine {
	return postgreSQL
}

type InMemoryConfig struct {
}

// func (i *InMemoryConfig) Engine() StorageEngine {
// 	return inMemory
// }

func NewStorageEngineConfig(engine StorageEngine) (StorageEngineConfig, error) {
	config := NewEnv()
	switch engine {
	case mongodb:
		if config.MongoDbURI() == "" {
			return nil, fmt.Errorf("MongoDbURI is required for the selected storage engine: %s", engine)
		}
		if config.MongoDbName() == "" {
			return nil, fmt.Errorf("MongoDbName is required for the selected storage engine: %s", engine)
		}
		return &MongoConfig{
			Uri:          config.MongoDbURI(),
			DatabaseName: config.MongoDbName(),
			Timeout:      config.Timeout(),
			ReplicaSet:   config.MongoDbReplicaSet(),
			MaxPoolSize:  config.MaxPoolSize(),
			MinPoolSize:  config.MinPoolSize(),
			ConnIdleTime: config.ConnIdleTime(),
			RetryWrites:  config.MongoDbRetryWrites(),
			RetryReads:   config.MongoDbRetryReads(),
		}, nil
	case postgreSQL:
		if config.PostgresURI() == "" {
			return nil, fmt.Errorf("PostgresURI is required for the selected storage engine: %s", engine)
		}
		return &PostgresConfig{
			Uri:          config.PostgresURI(),
			MaxPoolSize:  config.MaxPoolSize(),
			MinPoolSize:  config.MinPoolSize(),
			ConnIdleTime: config.ConnIdleTime(),
			Timeout:      config.Timeout(),
		}, nil
	default:
		return nil, ErrUnsupportedEngine
	}
}
