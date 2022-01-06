package plugin

import "github.com/kamikazechaser/hedwig/internal/message"

// NewPlugin descibes plugin signature
type NewPlugin func([]byte) (Plugin, error)

// Plugin descibes generic plugin
type Plugin interface {
	PluginName() string
	HealthCheck() bool
	Push(message message.Message) error
}
