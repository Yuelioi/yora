package bot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"yora/pkg/adapter"
	"yora/pkg/event"
	"yora/pkg/handler"
	"yora/pkg/log"
	"yora/pkg/middleware"
	"yora/pkg/plugin"
	"yora/pkg/provider"

	"github.com/rs/zerolog"
)

// 事件分发器
type EventDispatcher struct {
	middlewares []middleware.Middleware
	logger      zerolog.Logger
	mr          *plugin.MatcherRegistry
	mu          sync.RWMutex
	stats       EventStats

	// 事件队列
	shutdownCh chan struct{}
	eventQueue chan EventWrapper
}

// 事件包装器
type EventWrapper struct {
	Event    event.Event
	Adapter  adapter.Adapter
	Protocol adapter.Protocol
	RawData  []byte
	RecvTime time.Time
}

// 事件统计信息
type EventStats struct {
	TotalEvents   int64
	SuccessEvents int64
	FailedEvents  int64
	LastEventTime time.Time
	mu            sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {

	ed := &EventDispatcher{
		eventQueue:  make(chan EventWrapper, 100),
		mr:          plugin.GetMatcherRegistry(),
		logger:      log.NewMatcher("event_dispatcher"),
		middlewares: make([]middleware.Middleware, 0),
		stats:       EventStats{},
		shutdownCh:  make(chan struct{}),
	}

	go ed.startEventLoop()

	return ed

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

// 为每个适配器处理连接
func (ed *EventDispatcher) HandleAdapterConnection(w http.ResponseWriter, r *http.Request, a adapter.Adapter, p adapter.Protocol) {
	a.HandleWebSocket(w, r, func(message []byte) {
		go func() {
			if err := ed.processRawMessage(message, a, p); err != nil {
				ed.logger.Error().
					Err(err).
					Str("协议", string(p)).
					Msg("处理 WebSocket 消息失败")
			}
		}()
	})
}

// 处理原始消息
func (ed *EventDispatcher) processRawMessage(raw []byte, a adapter.Adapter, protocol adapter.Protocol) error {
	if len(raw) == 0 {
		return fmt.Errorf("原始数据不能为空")
	}

	// b.logger.Debug().
	// 	Str("协议", string(protocol)).
	// 	Int("数据长度", len(raw)).
	// 	Msg("开始处理原始消息")

	// 解析事件
	evt, err := a.ParseEvent(raw)
	if err != nil {
		ed.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Msg("事件解析失败")
		return fmt.Errorf("事件解析失败: %w", err)
	}

	// 忽略元事件
	if evt.Type() == "meta_event" {
		return nil
	}

	// 验证事件
	if err := a.ValidateEvent(evt); err != nil {
		ed.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件验证失败")
		return fmt.Errorf("事件验证失败: %w", err)
	}

	// 包装事件并放入队列
	wrapper := EventWrapper{
		Event:    evt,
		Adapter:  a,
		Protocol: protocol,
		RawData:  raw,
		RecvTime: time.Now(),
	}

	select {
	case ed.eventQueue <- wrapper:
		ed.logger.Debug().
			Str("协议", string(protocol)).
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件已加入处理队列")
	case <-time.After(time.Second): // 避免无限阻塞
		ed.logger.Warn().
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件处理队列已满，丢弃事件")
		return fmt.Errorf("事件队列已满")
	}

	return nil
}

// 处理事件包装器
func (ed *EventDispatcher) handleEventWrapper(wrapper EventWrapper) {
	startTime := time.Now()

	ed.logger.Debug().
		Str("协议", string(wrapper.Protocol)).
		Str("事件类型", fmt.Sprintf("%T", wrapper.Event)).
		Msg("开始处理事件")

	// 设置动态依赖注入
	provs := []provider.Provider{
		provider.Ctx(),
		provider.Event(),
		provider.MessageEvent(),
		provider.MetaEvent(),
		provider.RequestEvent(),
		provider.NoticeEvent(),
		BotProvider(),
	}

	handler.GetHandlerRegistry().ResetCache()
	handler.GetHandlerRegistry().RegisterProviders(provs...)

	// 通过分发器处理事件
	ctx := context.WithValue(context.Background(), "adapter", wrapper.Adapter)
	ctx = context.WithValue(ctx, "protocol", wrapper.Protocol)

	if err := ed.DispatchEvent(ctx, wrapper.Event); err != nil {
		ed.logger.Error().
			Err(err).
			Str("事件类型", fmt.Sprintf("%T", wrapper.Event)).
			Msg("事件分发失败")
	}

	duration := time.Since(startTime)
	ed.logger.Debug().
		Str("协议", string(wrapper.Protocol)).
		Str("事件类型", fmt.Sprintf("%T", wrapper.Event)).
		Dur("处理时长", duration).
		Msg("事件处理完成")
}

// 启动事件处理循环
func (ed *EventDispatcher) startEventLoop() {
	ed.logger.Info().Msg("启动事件处理器")

	for {
		select {
		case wrapper := <-ed.eventQueue:
			ed.handleEventWrapper(wrapper)
		case <-ed.shutdownCh:
			ed.logger.Debug().Msg("事件处理worker关闭")
			return
		}
	}
}

// 构建中间件链
func (ed *EventDispatcher) buildMiddlewareChain() func(ctx context.Context, event event.Event) error {
	ed.mu.RLock()
	middlewares := make([]middleware.Middleware, len(ed.middlewares))
	copy(middlewares, ed.middlewares)
	ed.mu.RUnlock()

	return middleware.Chain(middlewares, func(ctx context.Context, e event.Event) error {
		go ed.dispatchToMatchers(ctx, e)
		return nil
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
