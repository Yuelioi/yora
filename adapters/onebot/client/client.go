package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"yora/adapters/onebot/models"
	"yora/pkg/log"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	clientOnce     sync.Once
	clientInstance *Client
)

func GetClient(ctx context.Context) *Client {
	clientOnce.Do(func() {
		clientInstance = newClient(ctx)
	})
	return clientInstance
}

type Client struct {
	pending sync.Map
	logger  zerolog.Logger
	conn    *websocket.Conn
	sendCh  chan any
	ctx     context.Context
	mu      sync.RWMutex

	connCtx    context.Context    // 当前连接的上下文
	connCancel context.CancelFunc // 用于取消当前连接的所有 goroutine

	// 连接状态管理
	connClosed int64 // 原子操作标记连接是否已关闭

	// 监控指标
	metrics struct {
		activeGoroutines int64 // 当前活跃的goroutine数量
		messagesSent     int64 // 已发送的消息数量
		messagesReceived int64 // 已接收的消息数量
		reconnectCount   int64 // 重连次数
	}
}

func newClient(ctx context.Context) *Client {
	log := log.NewAPI("api")
	return &Client{
		logger:     log,
		sendCh:     make(chan any, 100),
		ctx:        ctx,
		pending:    sync.Map{},
		connClosed: 1, // 初始状态为已关闭
	}
}

// HandleWebSocket 用于处理 OneBot 反向连接
func (c *Client) HandleWebSocket(w http.ResponseWriter, r *http.Request, handleReceivedMessage func(message []byte)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error().Err(err).Msg("WebSocket升级失败")
		return
	}

	// 连接替换
	c.replaceConnection(conn)

	c.logger.Info().Msg("WebSocket连接已建立")

	// 启动连接管理器
	go c.connectionManager(handleReceivedMessage)
}

// 使用 Context 的连接替换方法，避免 channel 重新创建的问题
func (c *Client) replaceConnection(newConn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 取消旧连接的所有goroutine
	if c.connCancel != nil {
		c.logger.Info().Msg("取消旧连接的所有goroutine")
		c.connCancel()

		// 等待一小段时间让goroutine优雅退出
		time.Sleep(100 * time.Millisecond)
	}

	// 关闭旧连接
	if c.conn != nil {
		c.conn.Close()
	}

	// 设置新连接和新的context
	c.conn = newConn
	c.connCtx, c.connCancel = context.WithCancel(c.ctx)
	atomic.StoreInt64(&c.connClosed, 0)

	// 记录重连次数
	atomic.AddInt64(&c.metrics.reconnectCount, 1)

	c.logger.Debug().Msg("连接替换完成")
}

// 连接管理器 - 统一管理发送和接收goroutine的生命周期
func (c *Client) connectionManager(handleReceivedMessage func(message []byte)) {
	atomic.AddInt64(&c.metrics.activeGoroutines, 1)
	defer atomic.AddInt64(&c.metrics.activeGoroutines, -1)

	c.logger.Debug().Msg("连接管理器启动")

	// 启动发送循环
	sendDone := make(chan struct{})
	go func() {
		defer close(sendDone)
		c.sendLoop()
	}()

	// 启动接收循环
	receiveDone := make(chan struct{})
	go func() {
		defer close(receiveDone)
		c.receiveLoop(handleReceivedMessage)
	}()

	// 等待任一循环结束或context取消
	select {
	case <-sendDone:
		c.logger.Info().Msg("发送循环结束，连接管理器退出")
	case <-receiveDone:
		c.logger.Info().Msg("接收循环结束，连接管理器退出")
	case <-c.connCtx.Done():
		c.logger.Info().Msg("连接context取消，连接管理器退出")
	}

	// 确保连接被关闭
	c.closeConnection()

	// 等待所有goroutine结束
	<-sendDone
	<-receiveDone

	c.logger.Debug().Msg("连接管理器完全退出")
}

func (c *Client) receiveLoop(handleReceivedMessage func(message []byte)) {
	atomic.AddInt64(&c.metrics.activeGoroutines, 1)
	defer atomic.AddInt64(&c.metrics.activeGoroutines, -1)

	c.logger.Debug().Msg("接收循环启动")

	for {
		select {
		case <-c.connCtx.Done():
			c.logger.Debug().Msg("接收循环收到退出信号")
			return
		default:
			// 继续处理消息
		}

		// 检查连接状态
		if atomic.LoadInt64(&c.connClosed) == 1 {
			c.logger.Debug().Msg("连接已关闭，接收循环退出")
			return
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			c.logger.Debug().Msg("连接为空，接收循环退出")
			return
		}

		// 设置读取超时，避免长时间阻塞
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error().Err(err).Msg("WebSocket连接意外关闭")
			} else {
				c.logger.Debug().Err(err).Msg("WebSocket读取消息失败")
			}
			c.closeConnection()
			return
		}

		// 统计接收的消息数量
		atomic.AddInt64(&c.metrics.messagesReceived, 1)

		// 解析消息
		var baseResp map[string]any
		if err := json.Unmarshal(message, &baseResp); err != nil {
			c.logger.Error().Err(err).Str("raw_message", string(message)).Msg("解析基础响应失败")
			continue
		}

		// 处理不同类型的消息
		if echo, ok := baseResp["echo"].(string); ok {
			// API响应消息
			c.handleAPIResponse(message, echo)
		} else if _, ok := baseResp["post_type"]; ok {
			// 事件消息，异步处理
			go func(msg []byte) {
				atomic.AddInt64(&c.metrics.activeGoroutines, 1)
				defer atomic.AddInt64(&c.metrics.activeGoroutines, -1)
				handleReceivedMessage(msg)
			}(message)
		} else {
			c.logger.Warn().
				Str("message", string(message)).
				Msg("未知消息格式")
		}
	}
}

type Response struct {
	Echo string `json:"echo"`
}

func (c *Client) handleAPIResponse(data []byte, echo string) {
	var resp models.Response[any]

	if err := json.Unmarshal(data, &resp); err != nil {
		c.logger.Error().
			Err(err).
			Str("echo", echo).
			Str("data", string(data)).
			Msg("解析 API 响应失败")
		return
	}

	// 查找等待中的请求
	chVal, found := c.pending.Load(echo)
	if !found {
		c.logger.Warn().
			Str("echo", echo).
			Msg("收到未知echo的API响应")
		return
	}

	// 发送响应到等待的channel
	if ch, ok := chVal.(chan *models.Response[any]); ok {
		select {
		case ch <- &resp:
			c.logger.Debug().
				Str("echo", echo).
				Msg("API响应已发送到等待通道")
		case <-time.After(5 * time.Second):
			c.logger.Warn().
				Str("echo", echo).
				Msg("发送API响应超时，通道可能已满或已关闭")
		}
	} else {
		c.logger.Error().
			Str("echo", echo).
			Msg("响应通道类型断言失败")
	}
}

func (c *Client) sendLoop() {
	atomic.AddInt64(&c.metrics.activeGoroutines, 1)
	defer atomic.AddInt64(&c.metrics.activeGoroutines, -1)

	c.logger.Debug().Msg("发送循环启动")

	// 心跳定时器
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-c.sendCh:
			// 发送消息
			if err := c.writeMessage(msg); err != nil {
				c.logger.Error().Err(err).Msg("发送消息失败")
				c.closeConnection()
				return
			}
			atomic.AddInt64(&c.metrics.messagesSent, 1)

		case <-ticker.C:
			// 发送心跳
			if err := c.sendPing(); err != nil {
				c.logger.Error().Err(err).Msg("发送心跳失败")
				c.closeConnection()
				return
			}
			c.logger.Debug().Msg("心跳发送成功")

		case <-c.connCtx.Done():
			c.logger.Debug().Msg("发送循环收到退出信号")
			return
		}
	}
}

// 【新增】安全的消息发送方法
func (c *Client) writeMessage(msg any) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil || atomic.LoadInt64(&c.connClosed) == 1 {
		return nil // 连接已关闭，静默忽略
	}

	// 设置写入超时
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteJSON(msg)
}

// 【新增】安全的心跳发送方法
func (c *Client) sendPing() error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil || atomic.LoadInt64(&c.connClosed) == 1 {
		return nil // 连接已关闭，静默忽略
	}

	// 设置写入超时
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteMessage(websocket.PingMessage, nil)
}

// 【新增】安全的连接关闭方法
func (c *Client) closeConnection() {
	// 使用原子操作避免重复关闭
	if !atomic.CompareAndSwapInt64(&c.connClosed, 0, 1) {
		return // 已经关闭了
	}

	c.logger.Debug().Msg("开始关闭WebSocket连接")

	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()

	c.logger.Debug().Msg("WebSocket连接已关闭")
}

// 【新增】发送消息到队列的方法
func (c *Client) SendMessage(msg any) error {
	if atomic.LoadInt64(&c.connClosed) == 1 {
		return fmt.Errorf("连接已关闭")
	}

	select {
	case c.sendCh <- msg:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("发送队列已满，消息发送超时")
	}
}

// 【新增】发送API请求并等待响应的方法
func (c *Client) SendAPIRequest(request any, echo string, timeout time.Duration) (*models.Response[any], error) {
	if atomic.LoadInt64(&c.connClosed) == 1 {
		return nil, fmt.Errorf("连接已关闭")
	}

	// 创建响应通道
	respCh := make(chan *models.Response[any], 1)
	c.pending.Store(echo, respCh)

	// 确保清理pending
	defer func() {
		c.pending.Delete(echo)
		close(respCh)
	}()

	// 发送请求
	if err := c.SendMessage(request); err != nil {
		return nil, fmt.Errorf("发送API请求失败: %w", err)
	}

	// 等待响应
	select {
	case resp := <-respCh:
		return resp, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("API请求超时")
	case <-c.connCtx.Done():
		return nil, fmt.Errorf("连接已断开")
	}
}

// 【新增】获取监控指标的方法
func (c *Client) GetMetrics() map[string]int64 {
	return map[string]int64{
		"active_goroutines": atomic.LoadInt64(&c.metrics.activeGoroutines),
		"messages_sent":     atomic.LoadInt64(&c.metrics.messagesSent),
		"messages_received": atomic.LoadInt64(&c.metrics.messagesReceived),
		"reconnect_count":   atomic.LoadInt64(&c.metrics.reconnectCount),
		"connection_closed": atomic.LoadInt64(&c.connClosed),
	}
}

// 【新增】检查连接状态的方法
func (c *Client) IsConnected() bool {
	return atomic.LoadInt64(&c.connClosed) == 0
}

// 【新增】优雅关闭方法
func (c *Client) Close() error {
	c.logger.Info().Msg("开始关闭WebSocket客户端")

	// 取消所有goroutine
	if c.connCancel != nil {
		c.connCancel()
	}

	// 关闭连接
	c.closeConnection()

	// 等待所有goroutine退出（最多等待5秒）
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			activeCount := atomic.LoadInt64(&c.metrics.activeGoroutines)
			if activeCount > 0 {
				c.logger.Warn().
					Int64("active_goroutines", activeCount).
					Msg("等待goroutine退出超时，强制结束")
			}
			return nil
		case <-ticker.C:
			if atomic.LoadInt64(&c.metrics.activeGoroutines) == 0 {
				c.logger.Info().Msg("所有goroutine已优雅退出")
				return nil
			}
		}
	}
}

// 当前状态
func (c *Client) PrintStatus() {
	metrics := c.GetMetrics()
	c.logger.Info().
		Int64("active_goroutines", metrics["active_goroutines"]).
		Int64("messages_sent", metrics["messages_sent"]).
		Int64("messages_received", metrics["messages_received"]).
		Int64("reconnect_count", metrics["reconnect_count"]).
		Bool("is_connected", c.IsConnected()).
		Msg("WebSocket客户端状态")
}
