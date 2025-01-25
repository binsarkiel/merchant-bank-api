package api

import (
	"net/http"
	"strings"

	"merchant-bank-api/models"
	"merchant-bank-api/services"
	"merchant-bank-api/utils"

	"github.com/gin-gonic/gin"
)

// LoginHandler handles the login endpoint.
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err)
		return
	}

	token, err := services.Login(req.Username, req.Password)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleware is a middleware to authenticate requests using JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.HandleError(c, http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Memeriksa format "Bearer <token>"
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.HandleError(c, http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := headerParts[1] // Mengambil bagian token saja

		// Prioritaskan validasi token
		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			// Handle token expiration error separately
			if strings.Contains(err.Error(), "token has expired") {
				utils.HandleError(c, http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				utils.HandleError(c, http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		// Set claims ke context jika token valid
		c.Set("username", claims.Username)
		c.Set("account_type", claims.AccountType)

		c.Next()
	}
}

// PaymentHandler handles the payment endpoint.
func PaymentHandler(c *gin.Context) {
	// Ambil username dari context, yang sudah diset oleh AuthMiddleware
	username, _ := c.Get("username")
	senderUsername := username.(string)

	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pindahkan pengecekan error validasi data ke sini
	transaction, err := services.ProcessPayment(senderUsername, req.Recipient, req.Amount)
	if err != nil {
		if strings.Contains(err.Error(), "sender and recipient cannot be the same") ||
			strings.Contains(err.Error(), "insufficient balance") ||
			strings.Contains(err.Error(), "sender not found") ||
			strings.Contains(err.Error(), "recipient not found") {
			utils.HandleError(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			utils.HandleError(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer success!", "transaction": transaction})
}

// LogoutHandler handles the logout endpoint.
func LogoutHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	headerParts := strings.Split(authHeader, " ")
	tokenString := headerParts[1]

	remainingBalance, err := services.Logout(tokenString)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout success!", "remaining_balance": remainingBalance})
}
