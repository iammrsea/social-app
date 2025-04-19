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
	"github.com/iammrsea/social-app/internal/shared/config/db"
	"github.com/iammrsea/social-app/internal/user/infra/db/mongodb"
	"github.com/iammrsea/social-app/internal/user/service"
)

func main() {
	port := config.Env().Port()

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

	// Use any database of choice
	mongoDbConfig := db.NewMongoConfig()
	client, mongoDatabase, err := db.Connect(ctx, mongoDbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect(ctx, client)

	// Create repositories
	userRepo := mongodb.NewUserRepository(mongoDatabase)
	userReadModelRepo := mongodb.NewUserReadModelRepository(mongoDatabase)
	services := &internal.Services{
		UserService: service.NewUserService(userRepo, userReadModelRepo),
	}

	graphql.SetupHttGraphQLServer(router, services)

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
