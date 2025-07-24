package echo

import (
	"yora/internal/matcher"
	"yora/internal/plugin"
	"yora/protocols/onebot"
	"yora/protocols/onebot/message"
)

var _ plugin.Plugin = (*echo)(nil)

type echo struct {
	matchers []matcher.Matcher
}

var Echo = &echo{}

// Load implements plugin.Plugin.
func (e *echo) Load() error {
	banHandler := onebot.NewHandler(
		e.echo,
	)

	cmdMatcher := onebot.OnCommand("echo", nil, 10, false, banHandler)
	e.matchers = append(e.matchers, cmdMatcher)
	return nil
}

func (e *echo) Matchers() []matcher.Matcher {
	return e.matchers
}

// Metadata implements plugin.Plugin.
func (e *echo) Metadata() *plugin.Metadata {
	return &plugin.Metadata{
		ID:          "echo",
		Name:        "Echo",
		Description: "Echo back the message",
		Version:     "0.1.0",
		Author:      "Yoram",
		Usage:       "echo <message>",
		Extra:       nil,
	}
}

// Unload implements plugin.Plugin.
func (e *echo) Unload() error {
	return nil
}

func (e *echo) echo(event *onebot.Event, bot *onebot.Bot) error {
	bot.Send("group", "0", event.GroupID(), message.New(message.NewTextSegment("ECHO")).Append(message.NewAtSegment("435826135")))
	return nil

}
