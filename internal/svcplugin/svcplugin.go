package svcplugin

import "github.com/kamikazechaser/hedwig/internal/message"

// ServiceConf represents the configuration of the service loaded from the config file
type ServiceConf struct {
	Config map[string]interface{} `json:"conf"`
}

// NewService represents an init function that returns a Provider
type NewService func(jsonCfg map[string]interface{}) (Service, error)

// Service represents a messaging service e.g. Telegram
type Service interface {
	// ServiceName returns the messaging service identifier
	ServiceName() string

	// Push enqueues the message for delivery
	Push(message message.Message) error
}