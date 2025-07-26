package client

import (
	"fmt"
	"time"
)

func Call[ReqType any, RespType any](c *Client, action string, req ReqType) (*RespType, error) {

	// 2. 调用 c.CallAPI 发送请求
	apiResp, err := c.CallAPI(action, req)
	if err != nil {
		return nil, fmt.Errorf("调用 API %s 失败: %w", action, err)
	}

	// 3. 将通用的 API 响应转换为特定的响应结构体
	specificResp, err := convertStruct[any, RespType](apiResp)
	if err != nil {
		return nil, fmt.Errorf("转换 API %s 响应失败: %w", action, err)
	}

	return &specificResp, nil
}

func (c *Client) CallAPI(action string, params any) (*APIResponse, error) {
	echo := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())

	request := APIRequest{
		Action: action,
		Params: params,
		Echo:   echo,
	}

	ch := make(chan *APIResponse, 1)
	c.pending.Store(echo, ch)

	// 设置清理函数
	cleanup := func() {
		if chVal, exists := c.pending.Load(echo); exists {
			c.pending.Delete(echo)
			if c, ok := chVal.(chan *APIResponse); ok {
				close(c) // 关闭通道避免 goroutine 泄露
			}
		}
	}

	// 发送请求
	select {
	case c.sendCh <- request:
		c.logger.Debug().Msgf("发送 API 请求: %s (echo: %s)", action, echo)
	default:
		cleanup()
		c.logger.Error().Msg("发送队列已满，丢弃消息")
		return nil, fmt.Errorf("发送队列已满")
	}

	// 等待响应（带超时）
	timeout := time.NewTimer(10 * time.Second)
	defer timeout.Stop()

	select {
	case resp, ok := <-ch:
		cleanup()
		if !ok {
			c.logger.Error().Msg("通道已关闭")
			return nil, fmt.Errorf("通道已关闭")
		}
		c.logger.Debug().Msgf("收到 API 响应: %s (echo: %s) %v", action, echo, resp.Data)

		return resp, nil
	case <-timeout.C:
		cleanup()
		c.logger.Error().Msg("等待响应超时")
		return nil, fmt.Errorf("等待响应超时")
	}
}
