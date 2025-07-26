package manager

import (
	"yora/internal/plugin"
)

var Manager = &manager{}

var _ plugin.Plugin = (*manager)(nil)

type manager struct {
	plugin.BasePlugin
}

// Load implements plugin.Plugin.
func (m *manager) Load() error {
	panic("unimplemented")
}
