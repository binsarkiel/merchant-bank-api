package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	AccountType    string  `json:"account_type"` // "merchant" or "customer"
	AccountBalance float64 `json:"account_balance"`
}

// Session represents a user session.
type Session struct {
	Activity  string    `json:"activity"` // "logged_in" or "signed_out"
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

// Transaction represents a money transfer.
type Transaction struct {
	Activity      string    `json:"activity"` // "transfer_money"
	TransactionID uuid.UUID `json:"transaction_id"`
	Sender        string    `json:"sender"`    // User's name
	Recipient     string    `json:"recipient"` // User's name
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

// LoginRequest is used for login endpoint.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PaymentRequest is used for payment endpoint.
type PaymentRequest struct {
	Recipient string  `json:"recipient" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

// Claims is used for JWT claims.
type Claims struct {
	Username    string `json:"username"`
	AccountType string `json:"account_type"`
	jwt.RegisteredClaims
}
