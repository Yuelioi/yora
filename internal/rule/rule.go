package rule

import (
	"context"
	"yora/internal/condition"
	"yora/internal/event"
)

// 规则
type Rule condition.Condition

// RuleFunc 规则函数类型
type RuleFunc func(ctx context.Context, event event.Event) bool

func (f RuleFunc) Match(ctx context.Context, event event.Event) bool {
	return f(ctx, event)
}

// todo
func ToMe() RuleFunc {
	return func(ctx context.Context, event event.Event) bool {
		return event.SelfID() == "me"
	}
}
