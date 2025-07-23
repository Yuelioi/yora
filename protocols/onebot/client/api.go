package client

import (
	"encoding/json"
	"fmt"
	"time"
	"yora/protocols/onebot/client/models"
	"yora/protocols/onebot/message"
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

// GetPrivateFile 获取私聊文件
func (c *Client) GetPrivateFile(userID string, fileID string, fileHash string) (*models.GetPrivateFileResponse, error) {
	req := models.GetPrivateFileRequest{
		UserID:   userID,
		FileID:   fileID,
		FileHash: fileHash,
	}
	resp, err := Call[models.GetPrivateFileRequest, models.GetPrivateFileResponse](c, "get_private_file", req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SendMessage 发送消息
func (c *Client) SendMessage(messageType string, UserId int, GroupId int, message message.Message) (*models.SendMessageResponse, error) {

	req := models.MessageRequest{
		MessageType: messageType,
		UserID:      &UserId,
		GroupID:     &GroupId,
		Message:     message,
	}

	resp, err := Call[models.MessageRequest, models.SendMessageResponse](c, "send_msg", req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SendPrivateMessage 发送私聊消息
func (c *Client) SendPrivateMessage(userID int, message message.Message) (*models.SendPrivateMessageResponse, error) {
	req := models.MessageRequest{
		UserID:      &userID,
		Message:     message,
		MessageType: "private",
	}

	resp, err := Call[models.MessageRequest, models.SendPrivateMessageResponse](c, "send_private_msg", req)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

// SendGroupMessage 发送群消息
func (c *Client) SendGroupMessage(groupId int, message message.Message) (*models.SendMessageResponse, error) {
	req := models.MessageRequest{
		GroupID:     &groupId,
		Message:     message,
		MessageType: "group",
	}

	resp, err := Call[models.MessageRequest, models.SendMessageResponse](c, "send_group_msg", req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetMsg(id int) (*models.GetMessageResponse, error) {
	req := models.GetMessageRequest{
		MessageID: id,
	}

	resp, err := Call[models.GetMessageRequest, models.GetMessageResponse](c, "get_msg", req)
	if err != nil {
		return nil, err
	}
	return resp, nil
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
