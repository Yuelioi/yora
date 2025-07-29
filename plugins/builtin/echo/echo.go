package echo

import (
	"regexp"
	"strings"
	"time"
	"yora/adapters/onebot/event"
	"yora/adapters/onebot/message"
	"yora/pkg/bot"
	"yora/pkg/handler"
	"yora/pkg/on"
	"yora/pkg/plugin"
)

var _ plugin.Plugin = (*echo)(nil)

var pluginMeta = plugin.Metadata{
	ID:          "echo",
	Name:        "Echo",
	Description: "重复发送消息",
	Version:     "0.1.0",
	Author:      "月离",
	Usage:       "echo <message>",
	Extra:       make(map[string]any),
	Group:       "builtin",
}

func New() plugin.Plugin {
	return &echo{}
}

type echo struct {
	plugin.BasePlugin
}

func (e *echo) Load() error {

	cmdMatcher := on.OnCommand([]string{"echo"}, true, handler.NewHandler(e.echo)).SetPlugin(e)
	e.RegisterMatcher(cmdMatcher)

	e.SetMetadata(&pluginMeta)

	return nil
}

func (e *echo) echo(evt *event.MessageEvent, bot bot.Bot) error {
	var msgs = message.NewMessage()

	var echoRegex = regexp.MustCompile(`(?i)echo`)

	for _, seg := range evt.Message().Segments() {
		if seg.Type() == "text" {
			content := seg.String()
			cleaned := strings.TrimSpace(echoRegex.ReplaceAllString(content, ""))
			if cleaned != "" {
				s := message.NewTextSegment(cleaned)
				msgs = msgs.Append(s)
			}
			continue
		}
		msgs = msgs.Append(seg)
	}

	if msgs.IsEmpty() {
		return nil
	}
	time.Sleep(time.Second * 5)
	if evt.IsGroup() {
		bot.Send("0", evt.ChatID(), msgs)
	} else {
		bot.Send(evt.UserID(), "0", msgs)
	}
	return nil

}
