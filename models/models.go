package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	AccountType    string  `json:"account_type"`
	AccountBalance float64 `json:"account_balance"`
}

type Session struct {
	Activity  string    `json:"activity"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

type Transaction struct {
	Activity      string    `json:"activity"`
	TransactionID uuid.UUID `json:"transaction_id"`
	Sender        string    `json:"sender"`
	Recipient     string    `json:"recipient"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PaymentRequest struct {
	Recipient string  `json:"recipient" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

type Claims struct {
	Username    string `json:"username"`
	AccountType string `json:"account_type"`
	jwt.RegisteredClaims
}
