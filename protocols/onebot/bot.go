package onebot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"yora/internal/adapter"
	"yora/internal/event"
	"yora/internal/log"
	"yora/internal/plugin"
	"yora/protocols/onebot/message"

	"github.com/rs/zerolog"
)

// 确保 Bot 实现了 adapter.Bot 接口
var _ adapter.Bot = (*Bot)(nil)

var (
	BOT *Bot // 用于注入依赖
)

// Bot OneBot 协议机器人实现
type Bot struct {
	logger   zerolog.Logger   // 日志记录器
	pending  sync.Map         // 待处理请求映射表
	registry *AdapterRegistry // 适配器注册表
	manager  *PluginManager   // 插件管理器
	server   *http.Server     // HTTP 服务器
	mu       sync.RWMutex     // 读写锁
	config   map[string]any   // 配置信息
	running  bool             // 运行状态
}

// NewBot 创建新的机器人实例
func NewBot() *Bot {
	logger := log.NewBot("月灵Bot")

	bot := &Bot{
		logger:   logger,
		pending:  sync.Map{},
		manager:  NewPluginManager(),
		registry: NewAdapterRegistry(),
		mu:       sync.RWMutex{},
		config:   make(map[string]any),
		running:  false,
	}

	bot.logger.Info().Msg("机器人实例创建成功")

	BOT = bot
	return bot
}

// Platform 返回机器人平台标识
func (b *Bot) Platform() string {
	return "onebot"
}

// SelfID 返回机器人自身ID
func (b *Bot) SelfID() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if selfID, exists := b.config["self_id"]; exists {
		if id, ok := selfID.(string); ok {
			return id
		}
		b.logger.Warn().Interface("self_id", selfID).Msg("self_id 类型错误，期望 string")
	}

	b.logger.Warn().Msg("未找到 self_id 配置")
	return ""
}

// SetConfig 设置机器人配置
func (b *Bot) SetConfig(key string, value any) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.config[key] = value
	b.logger.Debug().Str("配置键", key).Interface("配置值", value).Msg("配置已更新")
}

// GetConfig 获取机器人配置
func (b *Bot) GetConfig(key string) (any, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	value, exists := b.config[key]
	return value, exists
}

// Run 启动机器人服务
func (b *Bot) Run() error {
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

	// 启动插件
	if err := b.manager.LoadPlugins(); err != nil {
		b.logger.Error().Err(err).Msg("插件加载失败")
		return fmt.Errorf("插件加载失败: %w", err)
	}

	// 创建HTTP服务器
	b.server = &http.Server{
		Addr:         ":12001",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	b.logger.Info().Str("地址", b.server.Addr).Msg("启动 HTTP 服务器")

	// 启动服务器（非阻塞）
	go func() {
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			b.logger.Error().Err(err).Msg("HTTP 服务器启动失败")
		}
	}()

	b.logger.Info().Msg("机器人服务启动完成")
	return nil
}

// setupRoutes 设置HTTP路由
func (b *Bot) setupRoutes() {
	// 健康检查端点
	http.HandleFunc("/", b.handleHealthCheck)
	// WebSocket 端点
	http.HandleFunc("/onebot/v11/ws", b.handleWebSocket)

	b.logger.Debug().Msg("HTTP 路由设置完成")
}

// handleHealthCheck 健康检查处理器
func (b *Bot) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	b.logger.Debug().
		Str("方法", r.Method).
		Str("路径", r.URL.Path).
		Str("客户端IP", r.RemoteAddr).
		Msg("收到健康检查请求")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := `{"status":"ok","message":"月灵Bot 运行正常","platform":"onebot"}`
	w.Write([]byte(response))
}

// handleWebSocket WebSocket连接处理器
func (b *Bot) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	b.logger.Info().
		Str("客户端IP", r.RemoteAddr).
		Str("User-Agent", r.Header.Get("User-Agent")).
		Msg("收到 WebSocket 连接请求")

	// 查找 OneBot 适配器
	b.registry.mu.RLock()
	defer b.registry.mu.RUnlock()

	for protocol, adapter := range b.registry.adapters {
		if protocol == adapter.Protocol() {
			b.logger.Debug().Msg("找到 OneBot 适配器，处理 WebSocket 连接")

			if onebotAdapter, ok := adapter.(*Adapter); ok {
				onebotAdapter.Client.HandleWebSocket(w, r, func(message []byte) {
					if err := b.Publish(message, adapter.Protocol()); err != nil {
						b.logger.Error().Err(err).Msg("处理 WebSocket 消息失败")
					}
				})
				return
			}

			b.logger.Error().Msg("OneBot 适配器类型转换失败")
			http.Error(w, "适配器类型错误", http.StatusInternalServerError)
			return
		}
	}

	b.logger.Error().Msg("未找到 OneBot 适配器")
	http.Error(w, "协议不支持", http.StatusNotFound)
}

// ShutDown 关闭机器人服务
func (b *Bot) ShutDown() error {
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

// CallAPI 调用协议API
func (b *Bot) CallAPI(params ...any) (any, error) {
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

	// 查找 OneBot 适配器
	b.registry.mu.RLock()
	defer b.registry.mu.RUnlock()

	for protocol, adapter := range b.registry.adapters {
		if protocol == adapter.Protocol() {
			b.logger.Debug().Msg("使用 OneBot 适配器调用API")

			result, err := adapter.CallAPI(action, apiParams)
			if err != nil {
				b.logger.Error().
					Err(err).
					Str("动作", action).
					Msg("API调用失败")
				return nil, fmt.Errorf("API调用失败: %w", err)
			}

			b.logger.Debug().
				Str("动作", action).
				Msg("API调用成功")
			return result, nil
		}
	}

	b.logger.Error().Str("动作", action).Msg("未找到支持的协议适配器")
	return nil, fmt.Errorf("未找到支持的协议")
}

// Send 发送消息
func (b *Bot) Send(messageType string, userId string, groupId string, msg message.Message) (any, error) {
	if msg == nil {
		b.logger.Error().Msg("发送消息失败：消息内容为空")
		return nil, fmt.Errorf("消息内容不能为空")
	}

	if messageType == "" {
		b.logger.Error().Msg("发送消息失败：消息类型为空")
		return nil, fmt.Errorf("消息类型不能为空")
	}

	b.logger.Debug().
		Str("消息类型", messageType).
		Str("用户ID", userId).
		Str("群组ID", groupId).
		Msg("发送消息")

	// 查找 OneBot 适配器
	b.registry.mu.RLock()
	defer b.registry.mu.RUnlock()

	for protocol, adapter := range b.registry.adapters {
		if protocol == adapter.Protocol() {
			b.logger.Debug().Msg("使用 OneBot 适配器发送消息")

			result, err := adapter.Send(messageType, userId, groupId, msg)
			if err != nil {
				b.logger.Error().
					Err(err).
					Str("消息类型", messageType).
					Str("用户ID", userId).
					Str("群组ID", groupId).
					Msg("消息发送失败")
				return nil, fmt.Errorf("消息发送失败: %w", err)
			}

			b.logger.Info().
				Str("消息类型", messageType).
				Str("用户ID", userId).
				Str("群组ID", groupId).
				Msg("消息发送成功")
			return result, nil
		}
	}

	b.logger.Error().Msg("未找到支持的协议适配器")
	return nil, fmt.Errorf("未找到支持的协议")
}

// Register 注册协议适配器
func (b *Bot) Register(adapter adapter.Adapter) error {
	if adapter == nil {
		b.logger.Error().Msg("注册适配器失败：适配器为空")
		return fmt.Errorf("适配器不能为空")
	}

	b.logger.Info().
		Str("协议", string(adapter.Protocol())).
		Msg("注册协议适配器")

	if err := b.registry.Register(adapter); err != nil {
		b.logger.Error().
			Err(err).
			Str("协议", string(adapter.Protocol())).
			Msg("适配器注册失败")
		return fmt.Errorf("适配器注册失败: %w", err)
	}

	b.logger.Info().
		Str("协议", string(adapter.Protocol())).
		Msg("适配器注册成功")
	return nil
}

// Unregister 注销协议适配器
func (b *Bot) Unregister(protocol adapter.Protocol) error {
	b.logger.Info().Str("协议", string(protocol)).Msg("注销协议适配器...")

	if err := b.registry.Unregister(protocol); err != nil {
		b.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Msg("适配器注销失败")
		return fmt.Errorf("适配器注销失败: %w", err)
	}

	b.logger.Info().Str("协议", string(protocol)).Msg("适配器注销成功")
	return nil
}

// Publish 发布事件到插件系统
func (b *Bot) Publish(raw any, protocol adapter.Protocol) error {
	if raw == nil {
		b.logger.Error().Msg("发布事件失败：原始数据为空")
		return fmt.Errorf("原始数据不能为空")
	}

	// 获取适配器
	adapter, err := b.registry.GetAdapter(protocol)
	if err != nil {
		b.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Msg("获取适配器失败")
		return fmt.Errorf("获取适配器失败: %w", err)
	}

	// 解析事件
	evt, err := adapter.ParseEvent(raw)
	if err != nil {
		b.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Msg("事件解析失败")
		return fmt.Errorf("事件解析失败: %w", err)
	}

	// 验证事件
	if err := adapter.ValidateEvent(evt); err != nil {
		b.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件验证失败")
		return fmt.Errorf("事件验证失败: %w", err)
	}

	b.logger.Debug().
		Str("协议", string(protocol)).
		Str("事件类型", fmt.Sprintf("%T", evt)).
		Msg("事件解析和验证成功")

	// 构建中间件链
	handler := func(ctx context.Context, e event.Event) error {
		return b.manager.Dispatch(ctx, e)
	}

	middlewares := b.registry.Middlewares()
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		next := handler
		handler = func(ctx context.Context, e event.Event) error {
			return middleware.Process(ctx, e, next)
		}
	}

	b.logger.Debug().
		Int("中间件数量", len(middlewares)).
		Msg("中间件链构建完成，处理事件")

	// 执行处理链
	ctx := context.Background()
	if err := handler(ctx, evt); err != nil {
		b.logger.Error().
			Err(err).
			Str("事件类型", fmt.Sprintf("%T", evt)).
			Msg("事件处理失败")
		return fmt.Errorf("事件处理失败: %w", err)
	}

	b.logger.Debug().Str("事件类型", fmt.Sprintf("%T", evt)).Msg("事件处理完成")
	return nil
}

// AddPlugin 添加插件到管理器
func (b *Bot) AddPlugin(plugins ...plugin.Plugin) {
	if len(plugins) == 0 {
		b.logger.Warn().Msg("尝试添加空插件列表")
		return
	}

	b.logger.Info().Int("插件数量", len(plugins)).Msg("添加插件")

	for i, p := range plugins {
		if p == nil {
			b.logger.Error().Int("插件索引", i).Msg("插件为空，跳过")
			continue
		}

		metadata := p.Metadata()
		if err := b.manager.RegisterPlugin(p); err != nil {
			b.logger.Error().
				Err(err).
				Str("插件名", metadata.Name).
				Msg("插件注册失败")
		}
	}
}

// GetAdapter 获取协议适配器
func (b *Bot) GetAdapter(protocol adapter.Protocol) (adapter.Adapter, error) {
	b.logger.Debug().Str("协议", string(protocol)).Msg("获取协议适配器")

	adapter, err := b.registry.GetAdapter(protocol)
	if err != nil {
		b.logger.Error().
			Err(err).
			Str("协议", string(protocol)).
			Msg("获取适配器失败")
		return nil, fmt.Errorf("获取适配器失败: %w", err)
	}

	b.logger.Debug().Str("协议", string(protocol)).Msg("适配器获取成功")
	return adapter, nil
}

// AddMiddleware 添加中间件
func (b *Bot) AddMiddleware(m adapter.Middleware) {
	if m == nil {
		b.logger.Warn().Msg("尝试添加空中间件")
		return
	}

	b.logger.Info().Msgf("添加中间件：%s", m.Name())
	b.registry.AddMiddleware(m)
}

// IsRunning 检查机器人是否正在运行
func (b *Bot) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// GetPluginManager 获取插件管理器
func (b *Bot) GetPluginManager() *PluginManager {
	return b.manager
}

// GetAdapterRegistry 获取适配器注册表
func (b *Bot) GetAdapterRegistry() *AdapterRegistry {
	return b.registry
}
