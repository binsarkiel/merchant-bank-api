package api

import (
	"net/http"

	"merchant-bank-api/models"
	"merchant-bank-api/services"

	"github.com/gin-gonic/gin"
)

// LoginHandler handles the login endpoint.
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleware is a middleware to authenticate requests using JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("account_type", claims.AccountType)
		c.Next()
	}
}

// PaymentHandler handles the payment endpoint.
func PaymentHandler(c *gin.Context) {
	username, _ := c.Get("username")
	senderUsername := username.(string) // Assert the type of the retrieved value

	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := services.ProcessPayment(senderUsername, req.Recipient, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// LogoutHandler handles the logout endpoint.
func LogoutHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	remainingBalance, err := services.Logout(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful", "remaining_balance": remainingBalance})
}
