package echo

import (
	"regexp"
	"strings"
	"yora/internal/bot"
	"yora/internal/matcher"
	"yora/internal/params"
	"yora/internal/plugin"
	"yora/protocols/onebot/event"
	"yora/protocols/onebot/message"
)

var _ plugin.Plugin = (*echo)(nil)

type echo struct {
	plugin.BasePlugin
}

var Echo = &echo{}

func (e *echo) Load() error {
	cmdMatcher := matcher.OnCommand([]string{"echo"}, true, matcher.NewHandler(e.echo).RegisterDependent(matcher.Event()))
	e.RegisterMatcher(cmdMatcher)

	e.SetMetadata(&plugin.Metadata{
		ID:          "echo",
		Name:        "Echo",
		Description: "Echo back the message",
		Version:     "0.1.0",
		Author:      "YueLi",
		Usage:       "echo <message>",
		Extra:       nil,
		Group:       "builtin",
	})

	return nil
}

func (e *echo) echo(event event.MessageEvent, bot bot.Bot, params params.CommandArgs) error {
	var msgs message.Message = message.New(nil)

	var echoRegex = regexp.MustCompile(`(?i)echo`)

	for _, seg := range event.Message().Segments() {

		if seg.Type() == "text" {
			content := seg.String()
			cleaned := strings.TrimSpace(echoRegex.ReplaceAllString(content, ""))
			if cleaned != "" {
				s := message.NewTextSegment(cleaned)
				msgs.Append(s)
			}
			continue
		}
		msgs.Append(seg)
	}

	if event.IsGroup() {
		bot.Send("group", "0", event.ChatID(), message.New(msgs))
	} else {
		bot.Send("private", event.UserID(), "0", message.New(msgs))
	}
	return nil

}
