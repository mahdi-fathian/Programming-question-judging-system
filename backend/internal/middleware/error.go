package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle different types of errors
			switch err {
			case gorm.ErrRecordNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}

			c.Abort()
		}
	}
} 