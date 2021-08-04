package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("key") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "secret key missing from request",
			})

			return
		}

		if c.Query("key") != app.key {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "invalid secret key",
			})

			return
		}

		c.Next()
	}
}
