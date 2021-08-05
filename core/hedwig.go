package main

import (
	"fmt"
	"log"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/knadh/koanf"

	"github.com/kamikazechaser/hedwig/internal/svcplugin"
)

// App represents the global application configuration
type App struct {
	enabledServices []string
	services        map[string]svcplugin.Service
	key             string
}

var (
	version = "dev"
	conf    = koanf.New(".")
	app     *App
	q       *asynq.Client
)

func main() {
	// load config and plugins
	services, err := initConfig()
	if err != nil {
		log.Fatalf("failed to load service: %v", err)
	}

	app = &App{
		services:        services,
		enabledServices: conf.Strings("enabledServices"),
		key:             conf.String("secretKey"),
	}

	router := gin.Default()

	router.Use(checkKey())
	router.GET("/stats", getStats)
	router.GET("/push/:service", pushMessage)
	router.POST("/push/all", pushAll)

	go endless.ListenAndServe(fmt.Sprintf(":%d", conf.Int("server.port")), router)

	qProcessor, qMux := initQueue()
	if err := qProcessor.Run(qMux); err != nil {
		log.Fatalf("failed to start queue processor: %v", err)
	}
}
