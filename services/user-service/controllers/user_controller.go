package controllers

import (
    "encoding/json"
    "net/http"
    "user-service/database"
    "user-service/models"
    "gorm.io/gorm"
    "log"
    "golang.org/x/crypto/bcrypt"
)

type UserInput struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Hash password
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword), err
}

// Register a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var input UserInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Hash the password
    hashedPassword, err := HashPassword(input.Password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    user := models.User{
        Username: input.Username,
        Email:    input.Email,
        Password: hashedPassword,
    }

    // Insert into the database
    result := database.DB.Create(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrDuplicatedKey {
            http.Error(w, "User already exists", http.StatusBadRequest)
        } else {
            log.Println(result.Error)
            http.Error(w, "Could not create user", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
