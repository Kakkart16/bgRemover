package controllers

import (
	"encoding/json"
	"history-service/database"
	"history-service/models"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func GetHistory(w http.ResponseWriter, r *http.Request) {
    // Parse the JWT token from the request
    tokenString := r.Header.Get("Authorization")

    if tokenString == "" {
        http.Error(w, "Missing token", http.StatusUnauthorized)
        return
    }

    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }

    // Parse the token
    claims := &jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET_KEY")), nil
    })
    
    if err != nil || !token.Valid {
        log.Println("Invalid token:", err)
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Get user ID from token claims
    userID := uint((*claims)["sub"].(float64))

    // Fetch the history for the user
    var history []models.History
    result := database.DB.Where("user_id = ?", userID).Find(&history)
    if result.Error != nil {
        http.Error(w, "Error fetching history", http.StatusInternalServerError)
        return
    }

    // Return the history as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(history)
}
