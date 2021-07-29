package main

import (
	"fmt"
	"log"
	"plugin"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"

	"github.com/kamikazechaser/hoodwink/internal/svcplugin"
)

var (
	conf = koanf.New(".")
)

func initConfig() (map[string]svcplugin.Service, error) {
	if err := conf.Load(file.Provider("config.json"), json.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config")
	}	

	enabledServices := conf.Strings("enabledServices")

	if len(enabledServices) < 1 { return nil, fmt.Errorf("no services enabled") }

	services := make(map[string]svcplugin.Service)

	for _, service := range enabledServices {
		svcFile := service + ".svc"

		plg, err := plugin.Open(svcFile)

		if err != nil {
			return nil, fmt.Errorf("could not load %s: %v", service, err)
		}

		newFunc, err := plg.Lookup("New")
		if err != nil {
			return nil, fmt.Errorf("New() function not found in %s: %v", service, err)
		}

		f, ok := newFunc.(func(map[string]interface{}) (interface{}, error))
		if !ok {
			return nil, fmt.Errorf("New() function is of invalid type (%T) in %s", newFunc, service)
		}

		var svcConf svcplugin.ServiceConf

		conf.Unmarshal("services." + service, &svcConf)
		if len(svcConf.Config) == 0 {
			return nil, fmt.Errorf("no config found for %s", service)
		}

		loadSvc, err := f(map[string]interface{}(svcConf.Config))
		if err != nil {
			return nil, fmt.Errorf("error initializing service %s: %v", loadSvc, err)
		}

		svcPlg, ok := loadSvc.(svcplugin.Service)
		if !ok {
			return nil, fmt.Errorf("loaded service plugin does not match svcplugin.Service interface")
		}

		services[svcPlg.ServiceName()] = svcPlg
	}

	return services, nil
}

// func injectConfig(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), configKey, config)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func main() {
	services, err := initConfig()

	if err != nil {
		log.Fatalf("failed to load service: %v", err)
	}

	fmt.Printf("%v\n", services)

	// chiRouter := chi.NewRouter()
	// chiRouter.Route("/notification", func(router chi.Router) {
	// 	router.Use(injectConfig)
	// 	router.Get("/telegram", telegramHandler)
	// })

	// log.Info().Msg("server listening on port " + conf.String("server.port"))

	// if err := http.ListenAndServe(":" +conf.String("server.port"), chiRouter); err != nil {
	// 	log.Fatal().Err(err).Msg("server couldn't start")
	// }
}
