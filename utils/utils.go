package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError handles errors and sends appropriate responses.
func HandleError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
