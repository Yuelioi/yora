package plugin

import (
	"context"
	"yora/internal/event"
)

type PluginManager interface {
	LoadPlugins() error
	GetPlugin(name string) (Plugin, error)
	Plugins() []Plugin
	RegisterPlugin(plugin Plugin) error

	Dispatch(ctx context.Context, event event.Event) error
}
