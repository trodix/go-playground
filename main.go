package main

import (
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

	// Setup routes
	r := mux.NewRouter()
	r.Use(middleware.JSONMiddleware)

	r.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users", userHandler.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	// Start the server
	log.Printf("Server is running on port %d", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), r))
}
