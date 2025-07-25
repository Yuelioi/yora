package condition

import (
	"context"
	"yora/internal/event"
)

// 通用事件规则
type Condition interface {
	Match(ctx context.Context, e event.Event) bool
}

type multiCondition struct {
	conditions []Condition
	combiner   func(results []bool) bool
}

func (mc *multiCondition) Match(ctx context.Context, e event.Event) bool {
	results := make([]bool, len(mc.conditions))
	for i, cond := range mc.conditions {
		results[i] = cond.Match(ctx, e)
	}
	return mc.combiner(results)
}

type singleCondition struct {
	condition Condition
	modifier  func(result bool) bool
}

func (sc *singleCondition) Match(ctx context.Context, e event.Event) bool {
	originalResult := sc.condition.Match(ctx, e)
	return sc.modifier(originalResult)
}

// Not 返回一个 Condition，当传入的 condition Match 为 true 时，返回 false；反之，返回 true。
func Not(condition Condition) Condition {
	return &singleCondition{
		condition: condition,
		modifier: func(result bool) bool {
			return !result
		},
	}
}

// Any 返回一个 Condition，当且仅当任意一个传入的 conditions Match 为 true 时，才返回 true。
func Any(conditions ...Condition) Condition {
	return &multiCondition{
		conditions: conditions,
		combiner: func(results []bool) bool {
			for _, r := range results {
				if r {
					return true // 只要有一个符合，就返回 true
				}
			}
			return false // 所有都不符合
		},
	}
}

// All 返回一个 Condition，当且仅当所有传入的 conditions 都 Match 为 true 时，才返回 true。
func All(conditions ...Condition) Condition {
	return &multiCondition{
		conditions: conditions,
		combiner: func(results []bool) bool {
			for _, r := range results {
				if !r {
					return false // 只要有一个不符合，就返回 false
				}
			}
			return true // 所有都符合
		},
	}
}
