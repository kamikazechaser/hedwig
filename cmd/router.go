package main

import (
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func startServer(host string, port string) {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.Default()

	if conf.Bool("api.auth") {
		router.Use(checkKey())
	}
	router.POST("/push", pushMessage)

	log.Debug().Msg("starting api server")
	go endless.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router)
}
