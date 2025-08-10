package rules

import (
	"context"
	"yora/adapters/onebot/events"
	"yora/adapters/onebot/messages"
	"yora/pkg/event"
	"yora/pkg/rule"
)

// todo
func ToMe() rule.RuleFunc {
	return func(ctx context.Context, e event.Event) bool {
		if e2, ok := e.(*events.MessageEvent); ok {
			if msg, ok := e2.Message().(messages.Message); ok {
				helper := messages.NewHelper(msg)
				return helper.IsAtMe(e.SelfID())
			}

		}
		return false
	}
}
