package api

import (
	"yora/internal/message"
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/client/models"
)

// 删除精华消息
func (api *API) DeleteEssenceMessage(messageID int) (*models.Response[any], error) {
	req := models.DeleteEssenceMessageRequest{
		MessageID: messageID,
	}
	return client.Call[models.DeleteEssenceMessageRequest, models.Response[any]](api.client, "delete_essence_message", req)
}

// 撤回消息
func (api *API) DeleteMessage(messageID int) (*models.Response[any], error) {
	req := models.RecallMessageRequest{
		MessageID: messageID,
	}
	return client.Call[models.RecallMessageRequest, models.Response[any]](api.client, "delete_msg", req)
}

// 私聊戳一戳
func (api *API) PrivatePoke(userID int) (*models.Response[any], error) {
	req := models.PrivatePokeRequest{
		UserID: userID,
	}
	return client.Call[models.PrivatePokeRequest, models.Response[any]](api.client, "friend_poke", req)
}

// 获取精华消息列表
func (api *API) GetEssenceMessageList(groupID int) (*models.GetEssenceMessageListResponse, error) {
	req := models.GetEssenceMessageListRequest{
		GroupID: groupID,
	}
	return client.Call[models.GetEssenceMessageListRequest, models.GetEssenceMessageListResponse](api.client, "get_essence_msg_list", req)

}

// 获取合并转发消息
func (api *API) GetForwardMessage(id string) (*models.GetForwardMessageResponse, error) {
	req := models.GetForwardMessageRequest{
		ID: id,
	}
	return client.Call[models.GetForwardMessageRequest, models.GetForwardMessageResponse](api.client, "get_forward_msg", req)

}

// 获取好友历史聊天记录
func (api *API) GetFriendChatHistory(userID int, messageID int, count int) (*models.GetFriendChatHistoryResponse, error) {
	req := models.GetFriendChatHistoryRequest{
		UserID:    userID,
		MessageID: messageID,
		Count:     count,
	}
	return client.Call[models.GetFriendChatHistoryRequest, models.GetFriendChatHistoryResponse](api.client, "get_friend_msg_history", req)

}

// 获取群历史聊天记录
func (api *API) GetGroupChatHistory(groupID int, messageID string, count int) (*models.GetGroupChatHistoryResponse, error) {
	req := models.GetGroupChatHistoryRequest{
		GroupID:   groupID,
		MessageID: messageID,
		Count:     count,
	}
	return client.Call[models.GetGroupChatHistoryRequest, models.GetGroupChatHistoryResponse](api.client, "get_group_msg_history", req)

}

// 获取消息
func (api *API) GetMessage(messageID int) (*models.GetMessageResponse, error) {
	req := models.GetMessageRequest{
		MessageID: messageID,
	}
	return client.Call[models.GetMessageRequest, models.GetMessageResponse](api.client, "get_msg", req)

}

// 群里戳一戳
func (api *API) GroupPoke(groupID int, userID int) (*models.Response[any], error) {
	req := models.GroupPokeRequest{
		GroupID: groupID,
		UserID:  userID,
	}
	return client.Call[models.GroupPokeRequest, models.Response[any]](api.client, "group_poke", req)
}

// 标记消息为已读
func (api *API) MarkMessageAsRead(messageID int) (*models.Response[any], error) {
	req := models.MarkMessageAsReadRequest{
		MessageID: messageID,
	}
	return client.Call[models.MarkMessageAsReadRequest, models.Response[any]](api.client, "mark_msg_as_read", req)
}

// 构造合并转发消息
func (api *API) ConstructForwardMessage(messages []models.MessageNode) (*models.ConstructForwardMessageResponse, error) {
	req := models.ConstructForwardMessageRequest{
		Messages: messages,
	}
	return client.Call[models.ConstructForwardMessageRequest, models.ConstructForwardMessageResponse](api.client, "send_forward_msg", req)

}

// 发送群AI语音
func (api *API) SendGroupAIVoice(character string, groupID int, text string, chatType int) (*models.SendGroupAIVoiceResponse, error) {
	req := models.SendGroupAIVoiceRequest{
		Character: character,
		GroupID:   groupID,
		Text:      text,
		ChatType:  chatType,
	}
	return client.Call[models.SendGroupAIVoiceRequest, models.SendGroupAIVoiceResponse](api.client, "send_group_ai_voice", req)

}

// 发送群聊合并转发消息
func (api *API) SendGroupForwardMessage(groupID int, messages []models.MessageNode) (*models.SendGroupForwardMessageResponse, error) {
	req := models.SendGroupForwardMessageRequest{
		GroupID:  groupID,
		Messages: messages,
	}
	return client.Call[models.SendGroupForwardMessageRequest, models.SendGroupForwardMessageResponse](api.client, "send_group_forward_msg", req)

}

// 发送消息
func (api *API) SendMessage(userID int, GroupId int, message message.Message) (*models.SendMessageResponse, error) {
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
	return client.Call[models.MessageRequest, models.SendMessageResponse](api.client, "send_msg", req)

}

// 发送私聊合并转发消息
func (api *API) SendPrivateForwardMessage(userID int, messages []models.MessageNode) (*models.SendPrivateForwardMessageResponse, error) {
	req := models.SendPrivateForwardMessageRequest{
		UserID:   userID,
		Messages: messages,
	}
	return client.Call[models.SendPrivateForwardMessageRequest, models.SendPrivateForwardMessageResponse](api.client, "send_private_forward_msg", req)

}

// 设置精华消息
func (api *API) SetEssenceMessage(messageID int) (*models.Response[any], error) {
	req := models.SetEssenceMessageRequest{
		MessageID: messageID,
	}
	return client.Call[models.SetEssenceMessageRequest, models.Response[any]](api.client, "set_essence_msg", req)
}

// 发送私聊消息
func (api *API) SendPrivateMessage(userID int, message message.Message) (*models.SendPrivateMessageResponse, error) {
	req := models.SendPrivateMessageRequest{
		UserID:  userID,
		Message: message,
	}
	return client.Call[models.SendPrivateMessageRequest, models.SendPrivateMessageResponse](api.client, "send_private_msg", req)

}
