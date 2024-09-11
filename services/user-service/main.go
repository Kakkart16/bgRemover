package main

import (
	"log"
	"net/http"
	"user-service/controllers"
	"user-service/database"
	"user-service/middleware"
	"user-service/models"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize the database connection
	database.ConnectDatabase()

	// Auto-migrate the User model to create the users table
	database.DB.AutoMigrate(&models.User{})

	// Initialize the router
	router := mux.NewRouter()

	// Define the user registration route
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Protected routes (require JWT authentication)
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.JWTAuthentication)
	// protectedRouter.HandleFunc("/history", controllers.GetUserHistory).Methods("GET")

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
