package main

import (
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/kamikazechaser/hedwig/internal/plugin"
)

const (
	confEnvOverridePrefix = "HEDWIG_"
)

var (
	hedwig map[string]plugin.Plugin
	conf   = koanf.New(".")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	confFile := file.Provider("config.toml")

	if err := conf.Load(confFile, toml.Parser()); err != nil {
		log.Fatal().Err(err).Msg("cannot load config file")
	}

	if err := conf.Load(env.Provider(confEnvOverridePrefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, confEnvOverridePrefix)), "_", ".")
	}), nil); err != nil {
		log.Fatal().Err(err).Msg("cannot load env variables")
	}

	if conf.Bool("debug.enabled") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug mode enabled")
		conf.Print()
	}

	hedwig = loadPlugins(conf.String("plugins.path"))

	startServer(conf.String("server.host"), conf.String("server.port"))
	initQueue()
}
