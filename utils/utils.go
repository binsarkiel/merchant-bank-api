package utils

import (
	"github.com/gin-gonic/gin"
)

// HandleError handles errors and sends appropriate responses.
func HandleError(c *gin.Context, statusCode int, err interface{}) {
	c.JSON(statusCode, gin.H{"error": err})
}
