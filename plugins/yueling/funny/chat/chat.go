package chat

import (
	"yora/adapters/onebot/events"
	"yora/adapters/onebot/rules"
	"yora/pkg/bot"
	"yora/pkg/handler"
	"yora/pkg/on"
	"yora/pkg/plugin"
)

var _ plugin.Plugin = (*chat)(nil)

type chat struct {
}

func New() plugin.Plugin {
	return &chat{}
}

func (c *chat) PluginInfo() *plugin.PluginInfo {
	return &plugin.PluginInfo{
		ID:          "chat",
		Name:        "Chat",
		Description: "Chat with other users",
		Version:     "0.0.1",
		Author:      "Yoram",
		Usage:       "@bot/bot_name [message]",
		Examples:    []string{"chat hello", "chat hi"},
		Group:       "funny",
		Extra:       map[string]any{},
	}
}

func (c *chat) Matchers() []*plugin.Matcher {
	handler := handler.NewHandler(c.chat)
	m := on.OnMessage(handler).AppendRule(rules.ToMe())
	return []*plugin.Matcher{
		m,
	}
}

func (c *chat) chat(bot bot.Bot, event *events.MessageEvent) error {

	// _, err := bot..Send("", event.ChatID(), ms)
	// if err != nil {
	// 	return err
	// }
	return nil

}
