package rule

import (
	"context"
	"yora/pkg/condition"
	"yora/pkg/event"
)

// 规则
type Rule condition.Condition

// RuleFunc 规则函数类型
type RuleFunc func(ctx context.Context, event event.Event) bool

func (f RuleFunc) Match(ctx context.Context, event event.Event) bool {
	return f(ctx, event)
}
