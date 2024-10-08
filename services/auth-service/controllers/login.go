package controllers

import (
	"auth-service/database"
	"auth-service/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds LoginCredentials
	json.NewDecoder(r.Body).Decode(&creds)

	var user models.User
	result := database.DB.Where("email = ?", creds.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := generateJWT(user)
	if err != nil {
		log.Println("Error generating JWT:", err)
		http.Error(w, "Could not login", http.StatusInternalServerError)
		return
	}

	// Return the token to the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func generateJWT(user models.User) (string, error) {

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not set in environment variables")
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,                               // User ID as the subject of the token
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
	})

	// Sign the token with our secret key
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
