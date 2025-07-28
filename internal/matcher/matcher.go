package matcher

import (
	"context"
	"yora/internal/condition"
	"yora/internal/depends"
	"yora/internal/event"
	"yora/internal/permission"
	"yora/internal/rule"
)

type Matcher struct {
	Rule       rule.Rule             // 规则(必须全部满足)
	Permission permission.Permission // 权限(任意满足即可)
	Priority   int                   // 优先级(越大越优先)
	Block      bool                  // 是否阻止事件传播
	Handlers   []*Handler            // 处理器
}

func (m *Matcher) SetPriority(priority int) *Matcher {
	m.Priority = priority
	return m
}

func (m *Matcher) SetBlock(block bool) *Matcher {
	m.Block = block
	return m
}

func (m *Matcher) AppendRule(rule rule.Rule) *Matcher {
	m.Rule = condition.All(m.Rule, rule)
	return m
}

func (m *Matcher) AppendPermission(permission permission.Permission) *Matcher {
	m.Permission = condition.Any(m.Permission, permission)
	return m
}

func (m *Matcher) AppendHandler(handler *Handler) *Matcher {
	m.Handlers = append(m.Handlers, handler)
	return m
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

func (m *Matcher) Handle(ctx context.Context, e event.Event, deps ...depends.Dependent) error {
	for _, h := range m.Handlers {
		if h != nil {
			// 注入运行时依赖
			h.RegisterDependent(deps...)
			// 构建依赖缓存
			err := h.BuildDependentType(ctx, e)
			if err != nil {
				return err
			}

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
