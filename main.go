package main

import (
	"log"
	"os"

	"merchant-bank-api/api"

	"github.com/gin-gonic/gin"
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

	// Initialize the Gin router
	router := gin.Default()

	// Define API routes
	router.POST("/login", api.LoginHandler)

	// Protected routes (require authentication)
	authorized := router.Group("/")
	authorized.Use(api.AuthMiddleware())
	{
		authorized.POST("/payment", api.PaymentHandler)
		authorized.POST("/logout", api.LogoutHandler)
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
