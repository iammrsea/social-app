package storage

import (
	"context"
	"fmt"

	"github.com/iammrsea/social-app/internal/shared/config"
	"github.com/iammrsea/social-app/internal/shared/storage/mongodb"
	"github.com/iammrsea/social-app/internal/shared/storage/postgres"
	"github.com/iammrsea/social-app/internal/user/domain"
	mongoUserRepo "github.com/iammrsea/social-app/internal/user/infra/repos/mongoimpl"
	pgUserRepo "github.com/iammrsea/social-app/internal/user/infra/repos/postgresimpl"
)

type Storage struct {
	Repos Repos
}

type Repos struct {
	UserRepo          domain.UserRepository
	UserReadModelRepo domain.UserReadModelRepository
}

func NewStorage(ctx context.Context, storageEngine config.StorageEngine) (*Storage, func() error, error) {
	storageConfig, err := config.NewStorageEngineConfig(storageEngine)
	if err != nil {
		return nil, nil, err
	}
	return buildStorage(ctx, storageConfig)
}

func buildStorage(ctx context.Context, storageConfig config.StorageEngineConfig) (*Storage, func() error, error) {
	switch conf := storageConfig.(type) {
	case *config.MongoConfig:
		return buildMongoRepos(ctx, conf)
	case *config.PostgresConfig:
		return buildPostgresRepos(ctx, conf)
	default:
		return nil, nil, fmt.Errorf("unsupported storage engine: %s", storageConfig.Engine())
	}
}

func buildMongoRepos(ctx context.Context, conf *config.MongoConfig) (*Storage, func() error, error) {
	db, closeStorage := mongodb.SetupMongoDB(ctx, conf)
	// Repositories
	storage := &Storage{
		Repos: Repos{
			UserRepo:          mongoUserRepo.NewUserRepository(db),
			UserReadModelRepo: mongoUserRepo.NewUserReadModelRepository(db),
		},
	}
	return storage, closeStorage, nil
}

func buildPostgresRepos(ctx context.Context, cf *config.PostgresConfig) (*Storage, func() error, error) {
	pool, closeStorage := postgres.SetupPostgreSQL(ctx, cf)
	// Repositories
	storage := &Storage{
		Repos: Repos{
			UserRepo:          pgUserRepo.NewUserRepository(pool),
			UserReadModelRepo: pgUserRepo.NewUserReadModelRepository(pool),
		},
	}
	return storage, closeStorage, nil
}
