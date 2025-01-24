package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"merchant-bank-api/models"
	"merchant-bank-api/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// JWT secret key (keep this secure in a real application)
var jwtKey []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		// Handle error appropriately, e.g., exit if .env is crucial
	}

	jwtKeyString := os.Getenv("JWT_SECRET_KEY")
	if jwtKeyString == "" {
		panic("JWT_SECRET_KEY not set in .env file")
	}
	jwtKey = []byte(jwtKeyString) // Konversi string ke []byte
}

// Login authenticates a user and generates a JWT token.
func Login(username, password string) (string, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return "", err
	}

	// In a real application, use a secure method like bcrypt to verify passwords
	if user.Password != password {
		return "", errors.New("invalid password")
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &models.Claims{
		Username:    user.Username,
		AccountType: user.AccountType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	// Record login activity
	session := models.Session{
		Activity:  "logged_in",
		Username:  username,
		Timestamp: time.Now(),
	}
	err = repository.AddSession(session)
	if err != nil {
		fmt.Println("Failed to record login activity:", err)
	}

	return tokenString, nil
}

// Logout handles user logout.
func Logout(tokenString string) (float64, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	// Invalidate token
	err = InvalidateToken(tokenString)
	if err != nil {
		return 0, fmt.Errorf("failed to invalidate token: %v", err)
	}

	// Record logout activity
	session := models.Session{
		Activity:  "signed_out",
		Username:  claims.Username,
		Timestamp: time.Now(),
	}
	err = repository.AddSession(session)
	if err != nil {
		fmt.Println("Failed to record logout activity:", err)
	}

	// Retrieve the user's balance after recording the logout
	remainingBalance, err := repository.GetUserBalance(claims.Username)
	if err != nil {
		fmt.Println("Failed to get user balance:", err)
		return 0, err
	}

	return remainingBalance, nil
}

// ValidateToken validates a JWT token.
func ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// InvalidateToken invalidates a JWT token by setting its expiration time to the past.
func InvalidateToken(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &models.Claims{})
	if err != nil {
		return fmt.Errorf("invalid token format: %v", err)
	}

	if claims, ok := token.Claims.(*models.Claims); ok {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now()) // Set expiration to now

		// (Optional) Store the invalidated token, e.g., in a database or cache.
		// For this example, we'll just print the invalidated token.
		invalidatedTokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			return fmt.Errorf("failed to sign invalidated token: %v", err)
		}
		fmt.Printf("Invalidated token: %s\n", invalidatedTokenString)

		return nil
	}

	return errors.New("invalid token claims")
}

// ProcessPayment handles the payment process.
func ProcessPayment(senderUsername string, recipientUsername string, amount float64) (*models.Transaction, error) {
	// Retrieve sender
	sender, err := repository.FindUserByUsername(senderUsername)
	if err != nil {
		return nil, fmt.Errorf("sender not found: %v", err)
	}

	// Retrieve recipient by USERNAME
	recipientUser, err := repository.FindUserByUsername(recipientUsername)
	if err != nil {
		return nil, fmt.Errorf("recipient not found: %v", err)
	}

	// Check if sender has sufficient balance
	if sender.AccountBalance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Update balances
	sender.AccountBalance -= amount
	recipientUser.AccountBalance += amount

	// Update sender's balance in repository
	err = repository.UpdateUserBalance(sender.Username, sender.AccountBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update sender's balance: %v", err)
	}

	// Update recipient's balance in repository
	err = repository.UpdateUserBalance(recipientUser.Username, recipientUser.AccountBalance)
	if err != nil {
		// If updating recipient's balance fails, revert sender's balance
		_ = repository.UpdateUserBalance(sender.Username, sender.AccountBalance+amount)
		return nil, fmt.Errorf("failed to update recipient's balance: %v", err)
	}

	// Record transaction
	transaction := models.Transaction{
		Activity:      "transfer_money",
		TransactionID: uuid.New(),
		Sender:        sender.Name,
		Recipient:     recipientUser.Name, // Simpan Nama Penerima di dalam transaksi
		Amount:        amount,
		CreatedAt:     time.Now(),
	}

	err = repository.AddTransaction(transaction)
	if err != nil {
		// Revert balances if transaction recording fails
		_ = repository.UpdateUserBalance(sender.Username, sender.AccountBalance+amount)
		_ = repository.UpdateUserBalance(recipientUser.Username, recipientUser.AccountBalance-amount)
		return nil, fmt.Errorf("failed to record transaction: %v", err)
	}

	return &transaction, nil
}
