package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"url_shortner/internal/database"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("JSON Decode Error: %v", err)
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	err = database.RegisterUser(creds.Username, creds.Email, creds.Password)
	if err != nil {
		log.Printf("Database Error: %v", err) // Log the exact error
		if err.Error() == "email already registered" {
			http.Error(w, `{"error": "Email already registered"}`, http.StatusConflict) // 409 Conflict
		} else {
			http.Error(w, `{"error": "Failed to register user"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})

}

var jwtSecret = []byte(getEnv("JWT_SECRET", "default_secret_key"))

// getEnv fetches the environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
func GenerateJWT(email string) (string, error) {
	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,                                 // Payload data
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expiry time (24 hours)
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("JSON Decode Error: %v", err)
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Retrieve user from the database
	var storedPassword string
	err = database.DB.QueryRow("SELECT password FROM users WHERE email=$1", creds.Email).Scan(&storedPassword)
	if err != nil {
		log.Printf("Database Error: %v", err)
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token (assuming you have a `GenerateJWT` function)
	token, err := GenerateJWT(creds.Email)
	if err != nil {
		log.Printf("Token Generation Error: %v", err)
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
