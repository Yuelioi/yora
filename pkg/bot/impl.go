package bot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"yora/pkg/adapter"
	"yora/pkg/conf"
	"yora/pkg/log"
	"yora/pkg/message"
	"yora/pkg/middleware"
	"yora/pkg/plugin"

	"github.com/rs/zerolog"
)

var (
	b    Bot
	once sync.Once
)

// GetBot 获取机器人实例
func NewBot(conf *conf.BotConfig) Bot {
	once.Do(func() {
		b = newBot(conf)
	})
	return b
}

func GetBot() Bot {
	if b == nil {
		panic("bot not initialized")
	}
	return b
}

type botImpl struct {
	config          *conf.BotConfig          // 机器人配置
	logger          zerolog.Logger           // 日志记录器
	pending         sync.Map                 // 待处理请求映射表
	adapterRegistry *adapter.AdapterRegistry // 适配器注册表
	dispatcher      *EventDispatcher         // 事件分派器
	manager         *plugin.PluginRegistry   // 插件管理器
	server          *http.Server             // HTTP 服务器
	mu              sync.RWMutex             // 读写锁
	running         bool                     // 运行状态

	// 事件队列
	shutdownCh chan struct{}
	eventQueue chan EventWrapper
}

func newBot(conf *conf.BotConfig) *botImpl {
	b := &botImpl{
		config:          conf,
		logger:          log.NewBot("yora"),
		adapterRegistry: adapter.NewAdapterRegistry(),
		manager:         plugin.GetPluginRegistry(),
		pending:         sync.Map{},
		server:          &http.Server{},
		dispatcher:      NewEventDispatcher(),
		running:         false,
		eventQueue:      make(chan EventWrapper, 100),
	}

	b.logger.Info().Msg("创建机器人")
	return b
}

func (b *botImpl) Config() *conf.BotConfig {
	return b.config
}

// 注册中间件
func (b *botImpl) RegisterMiddlewares(middlewares ...middleware.Middleware) error {
	return b.dispatcher.RegisterMiddlewares(middlewares...)
}

// 调用协议API
func (b *botImpl) CallAPI(params ...any) (any, error) {
	var (
		lastErr error
		result  any
	)
	// 参数验证
	if len(params) != 2 {
		b.logger.Error().Int("参数数量", len(params)).Msg("API调用参数数量错误")
		return nil, fmt.Errorf("参数数量错误，期望2个，实际%d个", len(params))
	}

	for i, param := range params {
		if param == nil {
			b.logger.Error().Int("参数索引", i).Msg("API调用参数为空")
			return nil, fmt.Errorf("参数[%d]不能为空", i)
		}
	}

	action, ok := params[0].(string)
	if !ok {
		b.logger.Error().Interface("参数类型", fmt.Sprintf("%T", params[0])).Msg("API动作参数类型错误")
		return nil, fmt.Errorf("动作参数必须是字符串类型，实际类型: %T", params[0])
	}

	apiParams, ok := params[1].(map[string]any)
	if !ok {
		b.logger.Error().Interface("参数类型", fmt.Sprintf("%T", params[1])).Msg("API参数类型错误")
		return nil, fmt.Errorf("参数必须是 map[string]any 类型，实际类型: %T", params[1])
	}

	b.logger.Debug().
		Str("动作", action).
		Interface("参数", apiParams).
		Msg("调用API")

	for p, a := range b.adapterRegistry.Adapters() {
		b.logger.Debug().Str("协议", string(p)).Msg("使用 Bot 适配器调用API")
		r, err := a.CallAPI(action, apiParams)
		if err != nil {
			b.logger.Error().
				Err(err).
				Str("动作", action).
				Msg("API调用失败")
			lastErr = err

		}
		result = r

		b.logger.Debug().
			Str("动作", action).
			Msg("API调用成功")
	}

	b.logger.Debug().Msg("使用 Bot 适配器调用API")

	return result, lastErr

}

func (b *botImpl) Platform() string {
	panic("unimplemented")
}

func (b *botImpl) Plugins() []plugin.Plugin {
	return b.manager.Plugins()
}

func (b *botImpl) RegisterAdapters(adapters ...adapter.Adapter) error {
	for _, a := range adapters {

		if a == nil {
			b.logger.Error().Msg("注册适配器失败：适配器为空")
			return fmt.Errorf("适配器不能为空")
		}

		b.logger.Info().
			Str("协议", string(a.Protocol())).
			Msg("注册协议适配器")

		if err := b.adapterRegistry.Register(a); err != nil {
			b.logger.Error().
				Err(err).
				Str("协议", string(a.Protocol())).
				Msg("适配器注册失败")
			return fmt.Errorf("适配器注册失败: %w", err)
		}

		b.logger.Info().
			Str("协议", string(a.Protocol())).
			Msg("适配器注册成功")
	}
	return nil
}

func (b *botImpl) RegisterPlugins(plugins ...plugin.Plugin) error {
	b.manager.RegisterPlugins(plugins...)
	return nil

}

// 检查机器人是否正在运行
func (b *botImpl) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// 获取SelfID
func (b *botImpl) SelfID() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.config.SelfID
}

// 发送消息
func (b *botImpl) Send(userId string, groupId string, msg message.Message) (any, error) {
	if msg == nil {
		b.logger.Error().Msg("发送消息失败：消息内容为空")
		return nil, fmt.Errorf("消息内容不能为空")
	}

	var lastErr error

	b.logger.Debug().
		Str("用户ID", userId).
		Str("群组ID", groupId).
		Msg("发送消息")

	for p, a := range b.adapterRegistry.Adapters() {
		b.logger.Debug().Msg("使用 Bot 适配器发送消息")

		_, err := a.Send(userId, groupId, msg)
		if err != nil {
			b.logger.Error().
				Err(err).
				Str("用户ID", userId).
				Str("群组ID", groupId).
				Str("协议", string(p)).
				Msg("消息发送失败")
			lastErr = fmt.Errorf("消息发送失败: %w", err)
		}
	}

	b.logger.Info().
		Str("用户ID", userId).
		Str("群组ID", groupId).
		Msg("消息发送成功")

	return nil, lastErr

}

// 启动机器人服务
func (b *botImpl) Run() error {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return fmt.Errorf("机器人已在运行中")
	}
	b.running = true
	b.mu.Unlock()

	b.logger.Info().Msg("启动机器人服务")

	// 设置路由
	b.setupRoutes()

	// 创建HTTP服务器
	b.server = &http.Server{
		Addr:         ":12001",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	b.logger.Info().Str("地址", b.server.Addr).Msg("启动 HTTP 服务器")

	b.startEventLoop(10)

	// 启动服务器
	go func() {
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			b.logger.Error().Err(err).Msg("HTTP 服务器启动失败")
		}
	}()

	b.logger.Info().Msg("机器人服务启动完成")
	return nil
}

// setupRoutes 设置HTTP路由
func (b *botImpl) setupRoutes() {
	// 健康检查端点
	http.HandleFunc("/", b.handleHealthCheck)
	// WebSocket 端点
	http.HandleFunc("/onebot/v11/ws", b.handleWebSocket)

	b.logger.Debug().Msg("HTTP 路由设置完成")
}

// handleHealthCheck 健康检查处理器
func (b *botImpl) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	b.logger.Debug().
		Str("方法", r.Method).Str("路径", r.URL.Path).Str("客户端IP", r.RemoteAddr).Msg("收到健康检查请求")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := `{"status":"ok","message":"月灵Bot 运行正常","platform":"onebot"}`
	w.Write([]byte(response))
}

// handleWebSocket WebSocket连接处理器
func (b *botImpl) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	b.logger.Info().
		Str("客户端IP", r.RemoteAddr).
		Str("User-Agent", r.Header.Get("User-Agent")).
		Msg("收到 WebSocket 连接请求")

	// 为每个适配器创建单独的处理逻辑
	for protocol, adapter := range b.adapterRegistry.Adapters() {
		b.handleAdapterConnection(w, r, adapter, protocol)
	}
}

// ShutDown 关闭机器人服务
func (b *botImpl) ShutDown() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		b.logger.Warn().Msg("机器人未在运行")
		return nil
	}

	b.logger.Info().Msg("关闭机器人服务...")

	// 关闭插件
	b.logger.Info().Msg("卸载插件...")
	plugins := b.manager.Plugins()
	unloadErrors := make([]error, 0)

	for _, p := range plugins {
		metadata := p.Metadata()
		b.logger.Debug().Str("插件名", metadata.Name).Msg("正在卸载插件")

		if err := p.Unload(); err != nil {
			b.logger.Error().
				Err(err).
				Str("插件名", metadata.Name).
				Msg("插件卸载失败")
			unloadErrors = append(unloadErrors, fmt.Errorf("插件 %s 卸载失败: %w", metadata.Name, err))
		} else {
			b.logger.Info().Str("插件名", metadata.Name).Msg("插件卸载成功")
		}
	}

	// 关闭 HTTP 服务器
	if b.server != nil {
		b.logger.Info().Msg("正在关闭 HTTP 服务器")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := b.server.Shutdown(ctx); err != nil {
			b.logger.Error().Err(err).Msg("HTTP 服务器关闭失败")
			return fmt.Errorf("HTTP 服务器关闭失败: %w", err)
		}

		b.logger.Info().Msg("HTTP 服务器关闭成功")
	}

	b.running = false
	b.logger.Info().Msg("机器人服务关闭完成")

	// 如果有插件卸载错误，返回第一个错误
	if len(unloadErrors) > 0 {
		return unloadErrors[0]
	}

	return nil
}
