package matcher

import (
	"context"
	"yora/internal/event"
)

// Rule 规则接口
type Rule interface {
	Check(ctx context.Context, event event.Event) bool
}

// RuleFunc 规则函数类型
type RuleFunc func(ctx context.Context, event event.Event) bool

func (f RuleFunc) Check(ctx context.Context, event event.Event) bool {
	return f(ctx, event)
}

// 全部规则
func AllRules(rules ...Rule) Rule {
	return RuleFunc(func(ctx context.Context, event event.Event) bool {
		for _, rule := range rules {
			if !rule.Check(ctx, event) {
				return false
			}
		}
		return true
	})
}

// 任意规则
func AnyRule(rules ...Rule) Rule {
	return RuleFunc(func(ctx context.Context, event event.Event) bool {
		for _, rule := range rules {
			if rule.Check(ctx, event) {
				return true
			}
		}
		return false
	})
}

// 非该规则
func NotRule(rule Rule) Rule {
	return RuleFunc(func(ctx context.Context, event event.Event) bool {
		return !rule.Check(ctx, event)
	})
}
