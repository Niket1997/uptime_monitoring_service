package handler

import (
	"github.com/gin-gonic/gin"
)

// ReturnError function to return error
func ReturnError(errText string, err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(500, gin.H{
			"status": errText,
			"error":  err.Error(),
		})
	}
}
