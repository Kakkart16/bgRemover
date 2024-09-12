package main

import (
	"log"
	"net/http"
	"os"
	"user-service/controllers"
	"user-service/database"
	"user-service/middleware"
	"user-service/models"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	port := os.Getenv("PORT")
	log.Println("Server is running on port: ",port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
