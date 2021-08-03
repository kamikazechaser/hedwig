package main

import "github.com/gin-gonic/gin"

func getStats(c *gin.Context) {
	c.JSON(200, gin.H{
		"ok":              true,
		"enabledServices": app.enabledServices,
		"commitVersion":   version,
		"serverMode":      gin.Mode(),
	})
}