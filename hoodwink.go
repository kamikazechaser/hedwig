package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Config struct {
	TelegramToken string
	RedditCredentials *reddit.Credentials
}

var (
	conf = koanf.New(".")
	config = &Config{}
)

func initConfig() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := conf.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	config.TelegramToken = conf.String("telegram.token")
	config.RedditCredentials = &reddit.Credentials{
		ID: conf.String("reddit.app_id"),
		Secret: conf.String("reddit.app_secret"),
		Username: conf.String("reddit.username"),
		Password: conf.String("reddit.password"),
	}
}

func injectConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "config", config)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	initConfig()

	chiRouter := chi.NewRouter()
	chiRouter.Route("/notification", func(router chi.Router) {
		router.Use(injectConfig)
		router.Get("/telegram", telegramHandler)   
	})

	log.Info().Msg("server listening on port " + conf.String("server.port"))

	if err := http.ListenAndServe(":" + conf.String("server.port"), chiRouter); err != nil {
		log.Fatal().Err(err).Msg("server couldn't start")
	}
}