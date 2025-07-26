package client

import (
	"encoding/json"
	ms "yora/internal/message"
	"yora/protocols/onebot/client/models"
	"yora/protocols/onebot/message"
)

// SendMessage 发送消息
func (c *Client) SendMessage(messageType string, UserId int, GroupId int, message ms.Message) (*models.SendMessageResponse, error) {

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
