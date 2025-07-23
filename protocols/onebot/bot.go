package onebot

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"yora/internal/adapter"
	"yora/internal/event"
	"yora/internal/log"
	"yora/internal/matcher"
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/message"

	"github.com/rs/zerolog"
)

type Bot struct {
	logger     zerolog.Logger
	ctx        context.Context
	pending    sync.Map
	client     *client.Client
	registry   *AdapterRegistry
	dispatcher *Dispatcher
	mu         sync.RWMutex
}

func NewBot(ctx context.Context) *Bot {

	logger := log.NewDatabaseLogger("全局")

	bot := &Bot{
		logger:     logger,
		ctx:        ctx,
		pending:    sync.Map{},
		client:     client.NewClient(ctx),
		registry:   NewAdapterRegistry(),
		dispatcher: NewDispatcher(),
		mu:         sync.RWMutex{},
	}

	http.HandleFunc("/onebot/v11/ws", bot.HandleWebSocket)

	// 健康检查端点
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Bot!"))
	})

	return bot
}

func (b *Bot) CallAPI(params ...any) (any, error) {
	for i := 0; i < len(params); i++ {
		if params[i] == nil {
			return nil, fmt.Errorf("params[%d] is nil", i)
		}
	}
	if len(params) != 2 {
		return nil, fmt.Errorf("params length is not 2")
	}

	action, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("action is not string")
	}

	ps, ok := params[1].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("params[1] is not map[string]any")
	}

	return b.client.CallAPI(action, ps)
}

func (b *Bot) Send(messageType string, userId string, groupId string, message message.Message) (any, error) {
	if message == nil {
		return nil, fmt.Errorf("message is nil")
	}

	gid, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, fmt.Errorf("groupId or userId is not int")
	}
	uid, err := strconv.Atoi(userId)

	if err != nil {
		return nil, fmt.Errorf("groupId or userId is not int")
	}

	return b.client.SendMessage(messageType, uid, gid, message)
}

func (b *Bot) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	b.client.HandleWebSocket(w, r, func(message []byte) { b.Publish(message, adapter.ProtocolOneBot) })
}

func (b *Bot) Client() *client.Client {
	return b.client
}

// Register 注册协议适配器
func (b *Bot) Register(adapter adapter.Adapter) error {
	return b.registry.Register(adapter)
}

// Unregister 注销协议适配器
func (b *Bot) Unregister(protocol adapter.Protocol) error {
	return b.registry.Unregister(protocol)
}

// Publish 发布事件
func (b *Bot) Publish(raw any, protocol adapter.Protocol) error {
	adapter, err := b.registry.Get(protocol)
	if err != nil {
		return err
	}

	evt, err := adapter.ParseEvent(raw)
	if err != nil {
		return fmt.Errorf("parse event failed: %w", err)
	}

	if err := adapter.ValidateEvent(evt); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}

	// 构建中间件链
	handler := func(ctx context.Context, e event.Event) error {
		return b.dispatcher.Dispatch(ctx, e)
	}

	for i := len(b.registry.middlewares) - 1; i >= 0; i-- {
		m := b.registry.middlewares[i]
		next := handler
		handler = func(ctx context.Context, e event.Event) error {
			return m.Process(ctx, e, next)
		}
	}

	return handler(b.ctx, evt)

}

// Subscribe 订阅事件
func (b *Bot) AddMatcher(matcher matcher.Matcher) error {
	return b.dispatcher.Subscribe(matcher)
}

// Unsubscribe 取消订阅
func (b *Bot) RemoveMatcher(matcher matcher.Matcher) error {
	return b.dispatcher.Unsubscribe(matcher)
}

// GetAdapter 获取协议适配器
func (b *Bot) GetAdapter(protocol adapter.Protocol) (adapter.Adapter, error) {
	return b.registry.Get(protocol)
}

func (b *Bot) AddMiddleware(m adapter.Middleware) {
	b.registry.AddMiddleware(m)
}
