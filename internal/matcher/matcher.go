package matcher

import (
	"context"
	"yora/internal/condition"
	"yora/internal/event"
	"yora/internal/permission"
	"yora/internal/rule"
)

type Matcher struct {
	Rule       rule.Rule
	Permission permission.Permission
	Priority   int
	Block      bool
	Handlers   []*Handler
}

func (m *Matcher) SetPriority(priority int) {
	m.Priority = priority
}

func (m *Matcher) SetBlock(block bool) {
	m.Block = block
}

func (m *Matcher) AppendPermission(permission permission.Permission) {
	m.Permission = condition.Any(m.Permission, permission)
}

func (m *Matcher) Match(ctx context.Context, e event.Event) bool {
	if m.Rule != nil && !m.Rule.Match(ctx, e) {
		return false
	}
	if m.Permission != nil && !m.Permission.Match(ctx, e) {
		return false
	}
	return true
}

func (m *Matcher) Handle(ctx context.Context, e event.Event) error {
	for _, h := range m.Handlers {
		if h != nil {
			if err := h.Call(ctx, e); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewMatcher(rule rule.Rule, handlers ...*Handler) *Matcher {
	return &Matcher{
		Rule:       rule,
		Permission: permission.Everyone(),
		Priority:   10,
		Block:      false,
		Handlers:   handlers,
	}
}

func (m *Matcher) AppendHandler(handler *Handler) {
	m.Handlers = append(m.Handlers, handler)
}
