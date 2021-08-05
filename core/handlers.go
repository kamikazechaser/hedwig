package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/kamikazechaser/hedwig/internal/message"
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
		serviceDestination := service.ServiceName()

		// TODO: Bind JSON directly to asynq payload
		payload, err := json.Marshal(message.Message{
			Service: serviceDestination,
			To:      "to",
			Title:   "title",
			Content: "body",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":      false,
				"message": "payload does not match message signature",
			})

			return
		}

		task := asynq.NewTask(serviceDestination, payload)
		info, err := q.Enqueue(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok":      false,
				"message": "sometng went wrong",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"message": info.ID,
		})

		return
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
