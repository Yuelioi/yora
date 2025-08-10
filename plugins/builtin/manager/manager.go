package manager

import (
	"yora/adapters/onebot/events"
	"yora/pkg/bot"
	"yora/pkg/plugin"
)

var _ plugin.Plugin = (*pluginDemo)(nil)

var pluginMeta = plugin.PluginInfo{
	ID:          "",
	Name:        "",
	Description: "",
	Version:     "",
	Author:      "",
	Usage:       "",
	Examples:    []string{},
	Group:       "",
	Extra:       map[string]any{},
}

func New() plugin.Plugin {
	return &pluginDemo{}
}

type pluginDemo struct {
}

// Matchers implements plugin.Plugin.
func (e *pluginDemo) Matchers() []*plugin.Matcher {
	panic("unimplemented")
}

// PluginInfo implements plugin.Plugin.
func (e *pluginDemo) PluginInfo() *plugin.PluginInfo {
	panic("unimplemented")
}

func (e *pluginDemo) Load() error {
	// cmdMatcher := on.OnCommand([]string{"echo"}, true, handler.NewHandler(e.echo))

	return nil
}

func (e *pluginDemo) echo(evt *events.MessageEvent, bot bot.Bot) error {
	return nil

}
