package bot

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"yora/pkg/adapter"
	"yora/pkg/depends"
	"yora/pkg/event"
	"yora/pkg/handler"
)

// 事件包装器
type EventWrapper struct {
	Event    event.Event
	Adapter  adapter.Adapter
	Protocol adapter.Protocol
	RawData  []byte
	RecvTime time.Time
}

// 为每个适配器处理连接
func (b *botImpl) handleAdapterConnection(w http.ResponseWriter, r *http.Request, a adapter.Adapter, p adapter.Protocol) {
	a.HandleWebSocket(w, r, func(message []byte) {
		go func() {
			if err := b.processRawMessage(message, a, p); err != nil {
				b.logger.Error().
					Err(err).
					Str("协议", string(p)).
					Msg("处理 WebSocket 消息失败")
			}
		}()
	})
}

// 处理原始消息
func (b *botImpl) processRawMessage(raw []byte, a adapter.Adapter, protocol adapter.Protocol) error {
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
		b.logger.Error().
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
		b.logger.Error().
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
	case b.eventQueue <- wrapper:
		b.logger.Debug().
			Str("协议", string(protocol)).
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件已加入处理队列")
	case <-time.After(time.Second): // 避免无限阻塞
		b.logger.Warn().
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件处理队列已满，丢弃事件")
		return fmt.Errorf("事件队列已满")
	}

	return nil
}

// 处理事件包装器
func (b *botImpl) handleEventWrapper(workerID int, wrapper EventWrapper) {
	startTime := time.Now()

	b.logger.Debug().
		Int("workerID", workerID).
		Str("协议", string(wrapper.Protocol)).
		Str("事件类型", fmt.Sprintf("%T", wrapper.Event)).
		Msg("开始处理事件")

	// 设置依赖注入
	deps := []depends.Dependent{
		depends.Ctx(),
		depends.Event(),
		depends.MessageEvent(),
		depends.MetaEvent(),
		depends.RequestEvent(),
		depends.NoticeEvent(),
		getBotDependent(),
	}

	handler.GetHandlerRegistry().ResetCache()
	handler.GetHandlerRegistry().RegisterDependent(deps...)

	// 通过分发器处理事件
	ctx := context.WithValue(context.Background(), "workerID", workerID)
	ctx = context.WithValue(ctx, "adapter", wrapper.Adapter)
	ctx = context.WithValue(ctx, "protocol", wrapper.Protocol)

	if err := b.dispatcher.DispatchEvent(ctx, wrapper.Event); err != nil {
		b.logger.Error().
			Err(err).
			Int("workerID", workerID).
			Str("事件类型", fmt.Sprintf("%T", wrapper.Event)).
			Msg("事件分发失败")
	}

	duration := time.Since(startTime)
	b.logger.Debug().
		Int("workerID", workerID).
		Str("协议", string(wrapper.Protocol)).
		Str("事件类型", fmt.Sprintf("%T", wrapper.Event)).
		Dur("处理时长", duration).
		Msg("事件处理完成")
}

// 启动事件处理循环
func (b *botImpl) startEventLoop(workerCount int) {

	b.logger.Info().Msg("启动事件处理器")
	b.logger.Debug().Int("worker数量", workerCount).Msg("启动事件处理循环")

	// 启动多个worker处理事件
	for i := 0; i < workerCount; i++ {
		go b.eventWorker(i)
	}
}

// 事件处理worker
func (b *botImpl) eventWorker(workerID int) {
	b.logger.Debug().Int("workerID", workerID).Msg("事件处理worker启动")

	for {
		select {
		case wrapper := <-b.eventQueue:
			b.handleEventWrapper(workerID, wrapper)
		case <-b.shutdownCh:
			b.logger.Debug().Int("workerID", workerID).Msg("事件处理worker关闭")
			return
		}
	}
}
