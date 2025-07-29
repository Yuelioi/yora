package condition

import (
	"context"
	"yora/pkg/event"
)

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
