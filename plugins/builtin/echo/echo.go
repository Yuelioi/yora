package echo

import (
	"regexp"
	"strings"
	"time"
	"yora/adapters/onebot/events"
	"yora/adapters/onebot/messages"
	"yora/pkg/bot"
	"yora/pkg/handler"
	"yora/pkg/log"
	"yora/pkg/on"
	"yora/pkg/plugin"

	"github.com/rs/zerolog"
)

var _ plugin.Plugin = (*echo)(nil)

func New() plugin.Plugin {
	return &echo{
		logger: log.NewPlugin("echo"),
		metadatas: &plugin.PluginInfo{
			ID:          "echo",
			Name:        "Echo",
			Description: "重复发送消息",
			Version:     "0.1.0",
			Author:      "月离",
			Usage:       "echo <message>",
			Extra:       make(map[string]any),
			Group:       "builtin",
		},
	}
}

type echo struct {
	metadatas *plugin.PluginInfo
	matchers  []*plugin.Matcher
	logger    zerolog.Logger
}

// PluginInfo implements plugin.Plugin.
func (e *echo) PluginInfo() *plugin.PluginInfo {
	return e.metadatas
}

func (e *echo) Matchers() []*plugin.Matcher {
	cmdMatcher := on.OnCommand([]string{"echo"}, true, handler.NewHandler(e.echo)).SetPlugin(e)

	return []*plugin.Matcher{
		cmdMatcher,
	}
}

func (e *echo) echo(evt *events.MessageEvent, bot bot.Bot) error {
	var msgs = messages.NewMessage()

	var echoRegex = regexp.MustCompile(`(?i)echo`)

	for _, seg := range evt.Message().Segments() {
		if seg.Type() == "text" {
			content := seg.String()
			cleaned := strings.TrimSpace(echoRegex.ReplaceAllString(content, ""))
			if cleaned != "" {
				s := messages.NewTextSegment(cleaned)
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
