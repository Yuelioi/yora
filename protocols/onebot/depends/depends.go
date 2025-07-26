package depends

import (
	"context"
	"fmt"
	"strings"
	"yora/internal/depends"
	"yora/internal/event"
	"yora/internal/params"
	onebotEvent "yora/protocols/onebot/event"
)

// 获取 onebot 的 message 事件
func OnebotEvent() depends.Dependent {
	return depends.DependentFunc(func(ctx context.Context, e event.Event) any {
		if e2, ok := e.(*onebotEvent.Event); ok {
			return e2
		}

		panic(fmt.Sprintf("event is not onebot event: %T", e))
	})
}

// 获取配置
// func Config() event.Handler {
// 	return event.DependencyFunc[*conf.BotConfig](func(ec *event.EventContext) (*conf.BotConfig, error) {
// 		return conf.GetBotConfig(), nil
// 	})
// }

func CommandArgs(cmds []string) depends.Dependent {
	return depends.DependentFunc(func(ctx context.Context, e event.Event) any {
		msgEvent, ok := e.(*onebotEvent.MessageEvent)
		if !ok {
			panic(fmt.Sprintf("message event is not onebot event: %T", e))
		}

		msg := msgEvent.Message()
		var argsText []string
		var matched bool

		for i, seg := range msg.Segments() {
			if seg.Type() != "text" {
				continue
			}

			content := seg.String()

			if !matched {
				// 第一个 text，尝试匹配命令
				for _, cmd := range cmds {
					if strings.HasPrefix(content, cmd) {
						rest := strings.TrimSpace(content[len(cmd):])
						if rest != "" {
							argsText = append(argsText, rest)
						}
						matched = true
						break
					}
				}
			} else {
				// 已匹配，其他 text 也收集
				if trimmed := strings.TrimSpace(content); trimmed != "" {
					argsText = append(argsText, trimmed)
				}
			}

			if matched && i == len(msg.Segments())-1 {
				break // 优化：末尾已处理完
			}
		}

		if !matched {
			return make([]string, 0)
		}

		// 空白切分参数
		args := strings.Fields(strings.Join(argsText, " "))
		cmdArgs := params.CommandArgs(args)
		return &cmdArgs
	})
}

// 获取用户信息
func UserInfo() depends.Dependent {
	return depends.DependentFunc(func(ctx context.Context, e event.Event) any {
		if msgEvent, ok := e.(*onebotEvent.MessageEvent); ok {
			return &UInfo{
				UID:      msgEvent.Sender().ID(),
				Nickname: msgEvent.Sender().Username(),
			}
		}

		return &UInfo{}
	})
}

type UInfo struct {
	UID      string `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}
