package poke

import (
	"yora/adapters/onebot/event"
	"yora/pkg/bot"
	"yora/pkg/handler"
	"yora/pkg/on"
	"yora/pkg/plugin"
)

var _ plugin.Plugin = (*pluginDemo)(nil)

var pluginMeta = plugin.Metadata{
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
	plugin.BasePlugin
}

func (e *pluginDemo) Load() error {
	cmdMatcher := on.OnNotice(handler.NewHandler(e.echo))

	e.RegisterMatcher(cmdMatcher)

	e.SetMetadata(&pluginMeta)

	return nil
}

func (e *pluginDemo) echo(evt *event.MessageEvent, bot bot.Bot) error {
	return nil

}
