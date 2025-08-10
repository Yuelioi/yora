package help

import (
	"yora/adapters/onebot/events"
	"yora/adapters/onebot/messages"
	"yora/pkg/bot"
	"yora/pkg/handler"
	"yora/pkg/on"
	"yora/pkg/params"
	"yora/pkg/plugin"
	"yora/pkg/provider"
)

var _ plugin.Plugin = (*helper)(nil)

var pluginMeta = &plugin.PluginInfo{
	ID:          "help",
	Name:        "帮助插件",
	Description: "提供帮助信息",
	Version:     "0.1.0",
	Author:      "月离",
	Usage:       "help [插件ID]",
	Group:       "builtin",
	Extra:       nil,
}

func New() plugin.Plugin {
	return &helper{}
}

type helper struct {
}

func (h *helper) Matchers() []*plugin.Matcher {
	helpHandler := handler.NewHandler(h.listPlugins).RegisterProviders(provider.CommandArgs([]string{"help"}))
	helpMatcher := on.OnCommand([]string{"help"}, true, helpHandler).SetPlugin(h)

	return []*plugin.Matcher{helpMatcher}

}

func (h *helper) PluginInfo() *plugin.PluginInfo {
	return pluginMeta
}

// PluginInfo implements plugin.Plugin.

func (h *helper) listPlugins(bot bot.Bot, event *events.MessageEvent, params *params.CommandArgs) {
	plugins := bot.Plugins()
	// 基于 group 分组

	filtered := make(map[string][]plugin.Plugin)
	for _, p := range plugins {
		if p.PluginInfo().Group == "" {
			filtered[""] = append(filtered[""], p)
		} else {
			filtered[p.PluginInfo().Group] = append(filtered[p.PluginInfo().Group], p)
		}
	}

	msgs := ""

	for group, ps := range filtered {
		if group == "" {
			msgs += "Available plugins:\n"
			for _, p := range ps {
				msgs += "- " + p.PluginInfo().Name + "\n"
			}
		} else {
			msgs += "Available plugins in group " + group + ":\n"
			for _, p := range ps {
				msgs += "- " + p.PluginInfo().Name + "\n"
			}
		}
	}

	if event.IsGroup() {
		bot.Send("0", event.ChatID(), messages.New(msgs))
	} else {
		bot.Send(event.UserID(), "0", messages.New(msgs))
	}
}
