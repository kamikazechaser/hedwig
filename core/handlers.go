package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok":              true,
		"enabledServices": app.enabledServices,
		"commitVersion":   version,
		"serverMode":      gin.Mode(),
	})
}

func pushMessage(c *gin.Context) {
	if service, ok := app.services[c.Param("service")]; ok {
		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"service": service.ServiceName(),
		})

		return
		// TODO: Implement push in a go routine
	}

	c.JSON(http.StatusNotFound, gin.H{
		"ok":      false,
		"message": "service not found",
	})
}

func pushAll(c *gin.Context) {
	// TODO: Fanout to all services
	c.JSON(http.StatusNotFound, gin.H{
		"ok":      false,
		"message": "not implemented",
	})
}
