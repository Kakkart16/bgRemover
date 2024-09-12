// main.go
package main

import (
    "log"
    "net/http"
    "history-service/database"
    "history-service/controllers"
	"history-service/models"
    "github.com/joho/godotenv"
    "os"
)

func main() {

    // Load .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize the database connection
	database.ConnectDatabase()

    // Automatically migrate the History model
    database.DB.AutoMigrate(&models.History{})

    // Set up the HTTP routes
    http.HandleFunc("/history", controllers.GetHistory)

    // Start the server
    port := os.Getenv("PORT")
    log.Println("Server is running on port: ", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
