package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMIDdleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Catch errors
		for _, err := range c.Errors {
			log.Printf("ERROR: %v", err)

			errType := int(err.Type)
			if errType != http.StatusInternalServerError {
				c.JSON(int(err.Type), err)
			} else {
				c.JSON(int(err.Type), "Internal Server Error")
			}
		}
	}
}
