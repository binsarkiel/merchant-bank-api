package services_test

import (
	"testing"

	"merchant-bank-api/services"
)

func TestHashPassword(t *testing.T) {
	// Arrange
	password := "password123"

	// Act
	hashedPassword, err := services.HashPassword(password)

	// Assert
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if !services.CheckPasswordHash(password, hashedPassword) {
		t.Fatalf("Hashed password does not match original password")
	}
}
