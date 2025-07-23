package onebot

import (
	"context"
	"fmt"
	"strconv"
	"yora/internal/event"
	"yora/internal/matcher"
)

// 获取 onebot 的 message 事件
func OnebotEvent() matcher.Dependent {
	return matcher.NewDependency(func(ctx context.Context, e event.Event) (any, error) {
		if e2, ok := e.(*Event); ok {
			return e2, nil
		}
		return nil, fmt.Errorf("not a message event")
	})
}

// 获取配置
// func Config() event.Dependent {
// 	return event.DependencyFunc[*conf.BotConfig](func(ec *event.EventContext) (*conf.BotConfig, error) {
// 		return conf.GetBotConfig(), nil
// 	})
// }

// 获取用户信息
func UserInfo() matcher.Dependent {
	return matcher.NewDependency(func(ctx context.Context, e event.Event) (any, error) {
		if msgEvent, ok := e.(*Event); ok {
			return &UInfo{
				UID:      strconv.Itoa(msgEvent.Sender().UserID),
				Nickname: msgEvent.Sender().Nickname,
				Role:     msgEvent.Sender().Role(),
			}, nil
		}

		return nil, fmt.Errorf("not a message event")
	})
}

type UInfo struct {
	UID      string `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}
