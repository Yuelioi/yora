package bot

import (
	"context"
	"fmt"
	"sync"
	"time"
	"yora/pkg/depends"
	"yora/pkg/event"
	"yora/pkg/log"
	"yora/pkg/middleware"
	"yora/pkg/plugin"

	"github.com/rs/zerolog"
)

type EventDispatcher struct {
	eventQueue  chan EventWrapper
	middlewares []middleware.Middleware
	logger      zerolog.Logger
	mr          *plugin.MatcherRegistry
	mu          sync.RWMutex
	stats       EventStats
}

type EventStats struct {
	TotalEvents   int64
	SuccessEvents int64
	FailedEvents  int64
	LastEventTime time.Time
	mu            sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		eventQueue: make(chan EventWrapper, 100),
		mr:         plugin.GetMatcherRegistry(),
		logger:     log.NewMatcher("event_dispatcher"),
	}
}

func (ed *EventDispatcher) RegisterMiddlewares(middlewares ...middleware.Middleware) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	for _, m := range middlewares {
		ed.logger.Info().Msgf("添加中间件：%s", m.Name())
		ed.middlewares = append(ed.middlewares, middlewares...)
	}

	ed.logger.Info().Int("中间件数量", len(ed.middlewares)).Msg("注册中间件")
	return nil
}

// 构建中间件链
func (ed *EventDispatcher) buildMiddlewareChain() func(ctx context.Context, event event.Event) error {
	ed.mu.RLock()
	middlewares := make([]middleware.Middleware, len(ed.middlewares))
	copy(middlewares, ed.middlewares)
	ed.mu.RUnlock()

	return middleware.Chain(middlewares, func(ctx context.Context, e event.Event) error {
		return ed.dispatchToMatchers(ctx, e)
	})
}

// 分发到匹配器
func (ed *EventDispatcher) dispatchToMatchers(ctx context.Context, e event.Event) error {
	matched := ed.mr.MatchedMatchers(ctx, e)

	ed.logger.Debug().
		Int("匹配数量", len(matched)).
		Str("事件类型", fmt.Sprintf("%T", e)).
		Msg("找到匹配的处理器")

	var lastErr error
	successCount := 0

	for _, matcher := range matched {
		if err := matcher.Call(ctx, e); err != nil {
			ed.logger.Error().
				Err(err).
				Str("匹配器", fmt.Sprintf("%T", matcher)).
				Msg("处理器执行失败")
			lastErr = err
		} else {
			successCount++
			ed.logger.Debug().
				Str("匹配器", fmt.Sprintf("%T", matcher)).
				Msg("处理器执行成功")
		}
	}

	ed.logger.Info().
		Int("匹配总数", len(matched)).
		Int("成功数量", successCount).
		Int("失败数量", len(matched)-successCount).
		Msg("事件分发完成")

	return lastErr
}

// 更新统计信息
func (ed *EventDispatcher) updateStats(success bool) {
	ed.stats.mu.Lock()
	defer ed.stats.mu.Unlock()

	ed.stats.TotalEvents++
	if success {
		ed.stats.SuccessEvents++
	} else {
		ed.stats.FailedEvents++
	}
	ed.stats.LastEventTime = time.Now()
}

// DispatchEvent 分发单个事件
func (ed *EventDispatcher) DispatchEvent(ctx context.Context, e event.Event) error {
	if ctx == nil {
		return fmt.Errorf("上下文不能为空")
	}
	if e == nil {
		return fmt.Errorf("事件不能为空")
	}

	ed.updateStats(true)

	ed.logger.Debug().
		Str("Bot ID", e.SelfID()).
		Str("事件类型", fmt.Sprintf("%T", e)).
		Msg("开始分发事件")

	// 构建中间件链
	handler := ed.buildMiddlewareChain()

	err := handler(ctx, e)
	if err != nil {
		ed.updateStats(false)
		ed.logger.Error().
			Err(err).
			Str("事件类型", fmt.Sprintf("%T", e)).
			Msg("事件分发失败")
	} else {
		ed.updateStats(true)
	}

	return err
}

func getBotDependent() depends.Dependent {
	return depends.DependentFunc(func(ctx context.Context, e event.Event) any {
		return GetBot()
	})
}
