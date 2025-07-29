package client

import (
	"encoding/json"
	"fmt"
	"time"
	"yora/adapters/onebot/message"
	"yora/adapters/onebot/models"
	"yora/pkg/log"
)

func Call[ReqType any, RespType any](c *Client, action string, req ReqType) (*RespType, error) {

	l := log.NewAPI("Call")

	// 2. 调用 c.CallAPI 发送请求
	apiResp, err := c.CallAPI(action, req)
	if err != nil {
		l.Error().Err(err).Msgf("调用 API %s 失败", action)
		return nil, fmt.Errorf("调用 API %s 失败: %w", action, err)
	}
	if apiResp == nil {
		l.Error().Msgf("API %s 调用失败: 响应为空", action)
		return nil, fmt.Errorf("API %s 调用失败: 响应为空", action)
	}
	// if apiResp.Status != "ok" {
	// 	return nil, fmt.Errorf("API %s 调用失败: %s", action, apiResp.Status)
	// }

	// 3. 将通用的 API 响应转换为特定的响应结构体
	specificResp, err := convertStruct[any, RespType](apiResp)
	if err != nil {
		l.Error().Err(err).Msgf("转换 API %s 响应失败", action)
		return nil, fmt.Errorf("转换 API %s 响应失败: %w", action, err)
	}

	return &specificResp, nil
}

// 发送消息
func (c *Client) Send(userID int, GroupId int, message message.Message) (*models.SendMessageResponse, error) {
	messageType := "private"
	if GroupId != 0 {
		messageType = "group"
	}

	req := models.MessageRequest{
		MessageType: messageType,
		UserID:      &userID,
		GroupID:     &GroupId,
		Message:     message,
	}
	return Call[models.MessageRequest, models.SendMessageResponse](c, "send_msg", req)

}

func (c *Client) CallAPI(action string, params any) (*models.Response[any], error) {
	echo := fmt.Sprintf("%s-%d", action, time.Now().UnixNano())

	request := models.APIRequest{
		Action: action,
		Params: params,
		Echo:   echo,
	}

	ch := make(chan *models.Response[any], 1)
	c.pending.Store(echo, ch)

	// 设置清理函数
	cleanup := func() {
		if chVal, exists := c.pending.Load(echo); exists {
			c.pending.Delete(echo)
			if c, ok := chVal.(chan *models.Response[any]); ok {
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

func convertStruct[A any, B any](input A) (B, error) {
	var out B
	raw, err := json.Marshal(input)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(raw, &out)
	return out, err
}
