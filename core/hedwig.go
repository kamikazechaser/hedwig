package main

import (
	"fmt"
	"log"
	"plugin"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"

	"github.com/kamikazechaser/hedwig/internal/svcplugin"
)

// App represents th golbal application configuration
type App struct {
	enabledServices []string
	services        map[string]svcplugin.Service
}

var (
	version = "dev"
	conf    = koanf.New(".")
)

func initConfig() (map[string]svcplugin.Service, error) {
	if err := conf.Load(file.Provider("config.json"), json.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config")
	}

	enabledServices := conf.Strings("enabledServices")
	if len(enabledServices) < 1 {
		return nil, fmt.Errorf("no services enabled")
	}

	services := make(map[string]svcplugin.Service)

	// Load all enabled services
	for _, service := range enabledServices {
		svcFile := fmt.Sprintf("%s.svc", service)

		plg, err := plugin.Open(svcFile)
		if err != nil {
			return nil, fmt.Errorf("could not open service file %s: %v", service, err)
		}

		// Ensure plugin implements New()
		newFunc, err := plg.Lookup("New")
		if err != nil {
			return nil, fmt.Errorf("plugin:New() function not found in %s: %v", service, err)
		}

		// Ensure the "New()" function signature is valid
		f, ok := newFunc.(func(map[string]interface{}) (interface{}, error))
		if !ok {
			return nil, fmt.Errorf("plugin:New() function is of invalid type (%T) in %s", newFunc, service)
		}

		var svcConf svcplugin.ServiceConf

		// Unmarshal koanf-loaded json service configuration into map struct
		// TODO: Skip unnecessary re-marshalling step inside plugins
		// Perhaps we can jsonStringify then koanf-unmarhsal it directly into a better defined config struct
		if err := conf.Unmarshal(fmt.Sprintf("services.%s", service), &svcConf); err != nil {
			return nil, err
		}

		if len(svcConf.Config) == 0 {
			return nil, fmt.Errorf("no config found for %s", service)
		}

		// Init plugin
		loadSvc, err := f(svcConf.Config)
		if err != nil {
			return nil, fmt.Errorf("error initializing service %s: %v", service, err)
		}

		// Ensure plugin matches our common service plugin interface
		svcPlg, ok := loadSvc.(svcplugin.Service)
		if !ok {
			return nil, fmt.Errorf("loaded service plugin does not match svcplugin:Service interface")
		}

		// Store loaded service plugins
		services[svcPlg.ServiceName()] = svcPlg
	}

	return services, nil
}

func main() {
	services, err := initConfig()
	if err != nil {
		log.Fatalf("failed to load service: %v", err)
	}

	app := &App{
		services:        services,
		enabledServices: conf.Strings("enabledServices"),
	}

	router := gin.Default()

	router.Use(
		func(c *gin.Context) {
			c.Set("app", app)
			c.Next()
		},
	)

	router.GET("/stats", getStats)

	endless.ListenAndServe(fmt.Sprintf(":%d", conf.Int("server.port")), router)
}
