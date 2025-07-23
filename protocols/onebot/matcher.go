package onebot

import (
	"context"
	"yora/internal/event"
	"yora/internal/matcher"
)

type Matcher struct {
	typ        string
	priority   int
	block      bool
	rule       matcher.Rule
	permission matcher.Permission
	Handlers   []matcher.Dependent
}

var _ matcher.Matcher = (*Matcher)(nil)

// NewMatcher 创建新的匹配器
func NewMatcher(typ string, priority int, block bool, rule matcher.Rule, permission matcher.Permission, handler Handler) *Matcher {

	return &Matcher{
		typ:        typ,
		rule:       rule,
		permission: permission,
		priority:   priority,
		block:      block,
	}
}

func (m *Matcher) Name() string {
	return m.typ
}

func (m *Matcher) Priority() int {
	return m.priority
}

func (m *Matcher) Block() bool {
	return m.block
}

// Check 检查事件是否匹配
func (m *Matcher) CanHandle(ctx context.Context, e event.Event) bool {
	// 检查事件类型
	if m.typ != "" && m.typ != e.Type() {
		return false
	}

	// 检查权限
	if m.permission != nil && !m.permission.Check(ctx, e) {
		return false
	}

	if m.rule != nil && !m.rule.Check(ctx, e) {
		return false
	}

	return true
}

// 处理事件
func (m *Matcher) Handle(ctx context.Context, e event.Event) error {
	for _, handler := range m.Handlers {
		if _, err := handler.Resolve(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// Append 添加处理函数
func (m *Matcher) Append(handler *Handler) *Matcher {
	m.Handlers = append(m.Handlers, handler)
	return m
}
