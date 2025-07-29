package rule

import (
	"context"
	obEvent "yora/adapters/onebot/event"
	"yora/adapters/onebot/message"
	"yora/pkg/event"
	"yora/pkg/rule"
)

// todo
func ToMe() rule.RuleFunc {
	return func(ctx context.Context, e event.Event) bool {
		if e2, ok := e.(*obEvent.MessageEvent); ok {
			if msg, ok := e2.Message().(message.Message); ok {
				helper := message.NewHelper(msg)
				return helper.IsAtMe(e.SelfID())
			}

		}
		return false
	}
}
