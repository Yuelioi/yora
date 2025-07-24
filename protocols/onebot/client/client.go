package client

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"yora/internal/log"

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

type Client struct {
	pending sync.Map
	logger  zerolog.Logger
	conn    *websocket.Conn
	sendCh  chan any
	ctx     context.Context
}

func NewClient(ctx context.Context) *Client {
	log := log.NewAPI("api")
	return &Client{
		logger:  log,
		sendCh:  make(chan any, 100),
		ctx:     ctx,
		pending: sync.Map{},
	}
}

// HandleWebSocket 用于处理 OneBot 反向连接
func (c *Client) HandleWebSocket(w http.ResponseWriter, r *http.Request, handleReceivedMessage func(message []byte)) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 在接受新连接时，关闭并清理旧连接（如果存在）
	if c.conn != nil {
		c.conn.Close()
	}

	c.conn = conn

	// 每次有新连接，就为这个连接启动新的发送和接收循环
	go c.sendLoop()
	c.receiveLoop(handleReceivedMessage)
}

func (c *Client) receiveLoop(handleReceivedMessage func(message []byte)) {
	for {
		if c.conn == nil {
			return
		}

		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.conn = nil
			return
		}

		var baseResp map[string]any
		if err := json.Unmarshal(message, &baseResp); err != nil {
			c.logger.Error().Err(err).Msg("解析基础响应失败")
			continue
		}

		if echo, ok := baseResp["echo"].(string); ok {
			c.handleAPIResponse(c.ctx, message, echo)
		} else if _, ok := baseResp["post_type"]; ok {
			go handleReceivedMessage(message)

		} else {
			c.logger.Warn().Str("message", string(message)).Msg("未知消息格式")
		}
	}
}

func (c *Client) handleAPIResponse(ctx context.Context, data []byte, echo string) {
	var resp APIResponse

	if err := json.Unmarshal(data, &resp); err != nil {
		c.logger.Error().Err(err).Str("echo", echo).Str("data", string(data)).Msg("解析 API 响应失败")
		return
	}

	chVal, found := c.pending.Load(echo)
	if found {
		if ch, ok := chVal.(chan *APIResponse); ok {
			select {
			case ch <- &resp:
			default:
				c.logger.Warn().Str("echo", echo).Msg("响应通道已满或已关闭")
			}
		} else {
			c.logger.Error().Str("echo", echo).Msg("响应通道类型断言失败")
		}
	}
}

func (c *Client) sendLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		if c.conn == nil {
			time.Sleep(time.Second)
			continue
		}

		select {
		case msg := <-c.sendCh:
			if err := c.conn.WriteJSON(msg); err != nil {
				c.logger.Error().Err(err).Msg("发送消息失败")
				c.conn = nil
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Error().Err(err).Msg("发送心跳失败")
				c.conn = nil
				return
			}
		case <-c.ctx.Done():
			c.logger.Info().Msg("sendLoop 被取消")
			return
		}
	}
}
