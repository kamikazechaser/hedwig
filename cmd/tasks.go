package main

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"strings"

	"github.com/rs/zerolog/log"

	service "github.com/kamikazechaser/hedwig/internal/plugin"
)

func loadPlugins(pluginsPath string) map[string]service.Plugin {
	services := make(map[string]service.Plugin)

	plugins, err := ioutil.ReadDir(pluginsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load plugins path")
	}

	if len(plugins) < 1 {
		log.Fatal().Err(err).Msg("atleast one plugin needs to be loaded")
	}

	log.Debug().Msg("loading plugins")
	for _, pluginFile := range plugins {
		pluginFileName := pluginFile.Name()

		isPlugin := strings.Split(pluginFileName, ".")

		if isPlugin[1] == "plugin" {
			log.Debug().Msg(fmt.Sprintf("%s/%s", pluginsPath, pluginFileName))
			// UNIX style path
			plg, err := plugin.Open(fmt.Sprintf("%s/%s", pluginsPath, pluginFileName))
			fmt.Println(err)
			if err != nil {
				log.Fatal().Err(err).Msgf("cannot load plugin file %s", pluginFileName)
				return nil
			}

			newFunc, err := plg.Lookup("New")
			if err != nil {
				log.Fatal().Err(err).Msgf("plugin:New() function not found in %s", pluginFileName)
				return nil
			}

			pluginFunc, ok := newFunc.(func() (interface{}, error))
			if !ok {
				log.Fatal().Err(err).Msgf("plugin:New() function is of invalid type (%T) in %s", newFunc, pluginFileName)
				return nil
			}

			initPlugin, err := pluginFunc()
			if err != nil {
				log.Fatal().Err(err).Msgf("error initializing plugin %s", pluginFileName)
				return nil
			}

			svcPlg, ok := initPlugin.(service.Plugin)
			if !ok {
				log.Fatal().Err(err).Msgf("loaded service plugin %s does not match plugin:Plugin interface", pluginFileName)
				return nil
			}

			services[svcPlg.PluginName()] = svcPlg
			log.Debug().Msgf("loaded %s", pluginFileName)
		}
	}

	return services
}
