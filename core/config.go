package main

import (
	"encoding/json"
	"fmt"
	"plugin"

	"github.com/kamikazechaser/hedwig/internal/svcplugin"
	kson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
)

func initConfig() (map[string]svcplugin.Service, error) {
	if err := conf.Load(file.Provider("config.json"), kson.Parser()); err != nil {
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
		f, ok := newFunc.(func([]byte) (interface{}, error))
		if !ok {
			return nil, fmt.Errorf("plugin:New() function is of invalid type (%T) in %s", newFunc, service)
		}

		var svcConf svcplugin.ServiceConf

		// Unmarshal koanf-loaded json service configuration into map struct
		// Ref: https://github.com/knadh/koanf/issues/76#issuecomment-853754910
		if err := conf.Unmarshal(fmt.Sprintf("services.%s", service), &svcConf); err != nil {
			return nil, err
		}

		if len(svcConf.Config) == 0 {
			return nil, fmt.Errorf("no config found for %s", service)
		}

		marshalledConf, err := json.Marshal(svcConf.Config)

		if err != nil {
			return nil, fmt.Errorf("unable to marshal into byte string for service %s", service)
		}

		// Init plugin
		loadSvc, err := f(marshalledConf)
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