package matcher

import (
	"context"
	"yora/internal/event"
)

// EventHandler 事件处理器接口
type Matcher interface {

	// Name 处理器名称
	Name() string

	// Priority 优先级（数字越小优先级越高）
	Priority() int

	// Block 是否阻止后续处理器继续处理
	Block() bool

	// Handle 处理事件
	Handle(ctx context.Context, event event.Event) error

	// CanHandle 判断是否能处理该事件
	CanHandle(ctx context.Context, event event.Event) bool
}
