package main

import (
	"log"
	"os"

	"merchant-bank-api/api"
	"merchant-bank-api/models"
	"merchant-bank-api/repository"
	"merchant-bank-api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Create data directory and files if they don't exist
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	createFileIfNotExist("data/users.json")
	createFileIfNotExist("data/sessions.json")
	createFileIfNotExist("data/transactions.json")

	// Setup initial users for testing purpose only
	setupInitialUsers()

	// Initialize the Gin router
	router := gin.Default()

	// Define API routes
	router.POST("/login", api.LoginHandler)

	// Protected routes (require authentication)
	authorized := router.Group("/")
	authorized.Use(api.AuthMiddleware()) // Pastikan AuthMiddleware terdaftar di sini
	{
		authorized.POST("/payment", api.PaymentHandler)
		authorized.DELETE("/logout", api.LogoutHandler)
	}

	// Run the server
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func createFileIfNotExist(filepath string) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		emptyData := []byte("[]") // Empty JSON array
		err = os.WriteFile(filepath, emptyData, 0644)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", filepath, err)
		}
	}
}

func setupInitialUsers() {
	users, err := repository.LoadUsers()
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	// Check if users are already initialized
	if len(users) > 0 {
		return // Skip initialization if users already exist
	}

	// Initialize with some default users
	initialUsers := []struct {
		Name           string
		Username       string
		Password       string
		AccountType    string
		AccountBalance float64
	}{
		{"John Doe", "johndoe", "password123", "customer", 1000},
		{"Jane Smith", "janesmith", "password456", "merchant", 100},
		{"Peter Jones", "peterjones", "password789", "customer", 0},
	}

	for _, u := range initialUsers {
		hashedPassword, err := services.HashPassword(u.Password)
		if err != nil {
			log.Fatalf("Failed to hash password for %s: %v", u.Username, err)
		}

		user := models.User{
			ID:             uuid.New().String(), // Generate UUID
			Name:           u.Name,
			Username:       u.Username,
			Password:       hashedPassword,
			AccountType:    u.AccountType,
			AccountBalance: u.AccountBalance,
		}
		users = append(users, user)
	}

	err = repository.SaveUsers(users)
	if err != nil {
		log.Fatalf("Failed to save initial users: %v", err)
	}
}
