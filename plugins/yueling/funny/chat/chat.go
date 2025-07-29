package chat

import (
	"yora/adapters/onebot/event"
	"yora/adapters/onebot/rule"
	"yora/pkg/bot"
	"yora/pkg/handler"
	"yora/pkg/on"
	"yora/pkg/plugin"
)

var _ plugin.Plugin = (*chat)(nil)

type chat struct {
	plugin.BasePlugin
}

func New() plugin.Plugin {
	return &chat{}
}

func (c *chat) Load() error {
	handler := handler.NewHandler(c.chat)
	m := on.OnMessage(handler).AppendRule(rule.ToMe())

	c.SetMetadata(&plugin.Metadata{
		ID:          "chat",
		Name:        "Chat",
		Description: "Chat with other users",
		Version:     "0.0.1",
		Author:      "Yoram",
		Usage:       "@bot/bot_name [message]",
		Examples:    []string{"chat hello", "chat hi"},
		Group:       "funny",
		Extra:       map[string]any{},
	})

	c.RegisterMatcher(m)

	return nil
}

func (c *chat) chat(bot bot.Bot, event *event.MessageEvent) error {

	// _, err := bot..Send("", event.ChatID(), ms)
	// if err != nil {
	// 	return err
	// }
	return nil

}
