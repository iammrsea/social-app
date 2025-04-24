package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iammrsea/social-app/cmd/server/graphql"
	"github.com/iammrsea/social-app/internal"
	"github.com/iammrsea/social-app/internal/shared/auth"
	"github.com/iammrsea/social-app/internal/shared/config"
	"github.com/iammrsea/social-app/internal/shared/guards"
	"github.com/iammrsea/social-app/internal/shared/storage"
	userService "github.com/iammrsea/social-app/internal/user/app"
)

func main() {
	port := config.NewEnv().Port()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(auth.AuthMiddleware)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	ctx := context.Background()

	// // Use any database of your choice
	// mongoDB, closeConnection := db.SetupMongoDB(ctx)

	storage, closeStorage, err := storage.NewStorage(ctx, config.StorageEngines.PostgreSQL)

	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	defer closeStorage()

	// Repositories
	userRepo := storage.Repos.UserRepo
	userReadModelRepo := storage.Repos.UserReadModelRepo

	// Guards
	guard := guards.New()

	services := &internal.Services{
		UserService: userService.New(userRepo, userReadModelRepo, guard),
	}

	graphql.SetupHttGraphQLServer(router, services)

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
