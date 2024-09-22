package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trodix/go-rest-api/api/handlers"
	"github.com/trodix/go-rest-api/api/middleware"
	"github.com/trodix/go-rest-api/config"
	"github.com/trodix/go-rest-api/database"
	"github.com/trodix/go-rest-api/repository"
)

func main() {
	// Load configuration and connect to the database
	cfg := config.LoadConfig() // Assume this function loads your configuration
	dbPool, db := database.ConnectPostgres(cfg.Database)

	// Migrate the database
    database.MigrateDatabase(db)

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)

	// Create handlers
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup routes
	r := mux.NewRouter()
	r.Use(middleware.JSONMiddleware)
	
	r.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users", userHandler.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/v1/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
