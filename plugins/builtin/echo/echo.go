package echo

import (
	"fmt"
	"regexp"
	"strings"
	"yora/internal/bot"
	"yora/internal/hook"
	"yora/internal/matcher"
	basemsg "yora/internal/message"
	"yora/internal/plugin"
	"yora/protocols/onebot/event"
	"yora/protocols/onebot/message"
)

var _ plugin.Plugin = (*echo)(nil)

type echo struct {
	plugin.BasePlugin
}

var Echo *echo

func init() {
	Echo = plugin.Register(&echo{BasePlugin: plugin.NewBasePlugin()}).
		WithHook(hook.PluginAfterLoad, func(ctx *hook.HookContext) error {
			fmt.Println("  -> Echo AfterLoad Hook Triggered")
			return nil
		}).Plugin
}

func (e *echo) Load() error {
	e.BasePlugin.Load()

	cmdMatcher := matcher.OnCommand([]string{"echo"}, true, matcher.NewHandler(e.echo))
	e.RegisterMatcher(cmdMatcher)

	e.SetMetadata(&plugin.Metadata{
		ID:          "echo",
		Name:        "Echo",
		Description: "重复发送消息",
		Version:     "0.1.0",
		Author:      "月离",
		Usage:       "echo <message>",
		Extra:       make(map[string]any),
		Group:       "builtin",
	})

	return nil
}

func (e *echo) echo(evt *event.MessageEvent, bot bot.Bot) error {
	var msgs basemsg.Message = message.NewMessage()

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

	if evt.IsGroup() {
		bot.Send("group", "0", evt.ChatID(), msgs)
	} else {
		bot.Send("private", evt.UserID(), "0", msgs)
	}
	return nil

}
