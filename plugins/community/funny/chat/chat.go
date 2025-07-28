package chat

import (
	"yora/internal/matcher"
	"yora/internal/plugin"
	"yora/protocols/onebot/bot"
	"yora/protocols/onebot/event"
	"yora/protocols/onebot/rule"
)

var _ plugin.Plugin = (*chat)(nil)

type chat struct {
	plugin.BasePlugin
}

func New() plugin.Plugin {
	return &chat{}
}

func (c *chat) Load() error {
	handler := matcher.NewHandler(c.chat)
	m := matcher.OnMessage(handler).AppendRule(rule.ToMe())

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

func (c *chat) chat(bot *bot.Bot, event *event.MessageEvent) error {

	// _, err := bot..Send("", event.ChatID(), ms)
	// if err != nil {
	// 	return err
	// }
	return nil

}
