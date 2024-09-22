package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trodix/go-rest-api/api/handlers"
	"github.com/trodix/go-rest-api/api/middleware"
	"github.com/trodix/go-rest-api/config"
	"github.com/trodix/go-rest-api/database"
	"github.com/trodix/go-rest-api/repository"
	"github.com/trodix/go-rest-api/service"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
)

func main() {
	// Load configuration and connect to the database
	cfg := config.LoadConfig() // Assume this function loads your configuration
	dbPool, db := database.ConnectPostgres(cfg.Database)

	// Migrate the database
	database.MigrateDatabase(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo)

	// Create handlers
	userHandler := handlers.NewUserHandler(userService)
	publicHandler := handlers.NewPublicHandler()

	// Create the ZITADEL resource server provider
	provider, err := rs.NewResourceServerClientCredentials(context.TODO(), cfg.OIDC.IssuerUrl, cfg.OIDC.ClientId, cfg.OIDC.ClientSecret)
	if err != nil {
		log.Fatalf("error creating provider: %s", err.Error())
	}
	log.Printf("Connected to OIDC issuer %s with clientId %s", cfg.OIDC.IssuerUrl, cfg.OIDC.ClientId)

	// Setup routes with Gorilla Mux
	r := mux.NewRouter()

	// Add global middleware for JSON handling and error handling
	r.Use(middleware.JSONMiddleware)
	r.Use(middleware.ErrorHandlingMiddleware)

	// Apply Keycloak authentication middleware to protected routes
	authenticatedRoutes := r.PathPrefix("/api/v1").Subrouter()
	authenticatedRoutes.Use(middleware.AuthMiddleware(provider)) // Applying the auth middleware

	// Public routes (no authentication needed)
	r.HandleFunc("/public/hello", publicHandler.Hello).Methods("GET")

	// Protected routes (authentication required)
	authenticatedRoutes.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	authenticatedRoutes.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	authenticatedRoutes.HandleFunc("/users/me", userHandler.GetMe).Methods("GET")
	authenticatedRoutes.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	authenticatedRoutes.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	authenticatedRoutes.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	// Start the server
	log.Printf("Server is running on port %d", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), r))
}
