package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	connStr := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the database!")
	DB = database
	return DB, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func RegisterUser(username, email, password string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	db, err := ConnectDB()
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	// Check if email already exists
	log.Printf("Attempting to register email: '%s'", email)

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return err
	}

	if exists {
		log.Printf("Email '%s' already exists in database", email)
		return fmt.Errorf("email already registered")
	}

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
func AuthenticateUser(email, password string) (bool, error) {
	db, err := ConnectDB()
	if err != nil {
		return false, fmt.Errorf("database connection error: %v", err)
	}

	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return false, fmt.Errorf("user not found")
	} else if err != nil {
		return false, fmt.Errorf("database error: %v", err)
	}

	if !CheckPasswordHash(password, hashedPassword) {
		return false, fmt.Errorf("invalid credentials")
	}

	return true, nil
}
