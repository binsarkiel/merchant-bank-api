package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"merchant-bank-api/models"
	"merchant-bank-api/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte
var invalidTokens = sync.Map{}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Cost factor 14
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login authenticates a user and generates a JWT token.
func Login(username, password string) (string, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return "", err
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	// Create JWT token
	expirationTime := time.Now().Add(1 * time.Minute)
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
	InvalidateToken(tokenString)

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

	remainingBalance, err := repository.GetUserBalance(claims.Username)
	if err != nil {
		fmt.Println("Failed to get user balance:", err)
		return 0, err
	}

	return remainingBalance, nil
}

// ValidateToken validates a JWT token.
func ValidateToken(tokenString string) (*models.Claims, error) {
	// Cek apakah token ada di daftar token yang tidak valid
	_, exists := invalidTokens.Load(tokenString)
	if exists {
		return nil, errors.New("invalid token")
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// InvalidateToken adds a token to the invalid tokens list.
func InvalidateToken(tokenString string) {
	invalidTokens.Store(tokenString, true)
}

func ProcessPayment(senderUsername string, recipientUsername string, amount float64) (*models.Transaction, error) {
	// **Periksa apakah pengirim dan penerima sama**
	if senderUsername == recipientUsername {
		return nil, errors.New("sender and recipient cannot be the same")
	}

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
		Sender:        sender.Username,
		Recipient:     recipientUser.Username,
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
