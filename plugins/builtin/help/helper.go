package help

import (
	"yora/internal/bot"
	"yora/internal/matcher"
	"yora/internal/params"
	"yora/internal/plugin"
	"yora/protocols/onebot/depends"
	"yora/protocols/onebot/event"
	"yora/protocols/onebot/message"
)

var _ plugin.Plugin = (*helper)(nil)

var Helper = &helper{}

type helper struct {
	plugin.BasePlugin
}

// Load implements plugin.Plugin.
func (h *helper) Load() error {

	helpHandler := matcher.NewHandler(h.listplugins).RegisterDependent(depends.CommandArgs([]string{"help"}))
	helpMatcher := matcher.OnCommand([]string{"help"}, true, helpHandler)
	h.RegisterMatcher(helpMatcher)

	h.SetMetadata(&plugin.Metadata{
		ID:          "help",
		Name:        "帮助插件",
		Description: "提供帮助信息",
		Version:     "0.1.0",
		Author:      "月离",
		Usage:       "help [插件ID]",
		Group:       "builtin",
		Extra:       nil,
	})

	return nil
}

// Metadata implements plugin.Plugin.

func (h *helper) listplugins(bot bot.Bot, event *event.MessageEvent, params *params.CommandArgs) {
	plugins := bot.Plugins()
	// 基于 group 分组

	filtered := make(map[string][]plugin.Plugin)
	for _, p := range plugins {
		if p.Metadata().Group == "" {
			filtered[""] = append(filtered[""], p)
		} else {
			filtered[p.Metadata().Group] = append(filtered[p.Metadata().Group], p)
		}
	}

	msgs := ""

	for group, ps := range filtered {
		if group == "" {
			msgs += "Available plugins:\n"
			for _, p := range ps {
				msgs += "- " + p.Metadata().Name + "\n"
			}
		} else {
			msgs += "Available plugins in group " + group + ":\n"
			for _, p := range ps {
				msgs += "- " + p.Metadata().Name + "\n"
			}
		}
	}

	if event.IsGroup() {
		bot.Send("0", event.ChatID(), message.New(msgs))
	} else {
		bot.Send(event.UserID(), "0", message.New(msgs))
	}
}
