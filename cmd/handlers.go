package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/kamikazechaser/hedwig/internal/message"
)

type pushBody struct {
	To      string   `json:"to" binding:"required"`
	Service string   `json:"service" binding:"required"`
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Delay   int      `json:"delay"`
	Params  []string `json:"params"`
}

func pushMessage(c *gin.Context) {
	var jsonBody pushBody

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ok":      false,
			"message": "json body validation failed",
		})
		return
	}

	if service, ok := hedwig[jsonBody.Service]; ok {
		serviceDestination := service.PluginName()
		servicePayload := message.Message{
			Service: serviceDestination,
			To:      jsonBody.To,
			Title:   jsonBody.Title,
			Content: jsonBody.Content,
			Params:  nil,
		}

		if len(jsonBody.Params) > 1 {
			servicePayload.Params = jsonBody.Params
		}

		payload, _ := json.Marshal(servicePayload)

		task := asynq.NewTask(serviceDestination, payload)
		var taskArgs = []asynq.Option{}

		if jsonBody.Delay > 0 {
			taskArgs = []asynq.Option{
				asynq.ProcessIn(time.Second * time.Duration(jsonBody.Delay)),
			}
		}

		info, err := qClient.Enqueue(task, taskArgs...)
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
