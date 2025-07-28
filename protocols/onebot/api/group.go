package api

import (
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/client/models"
)

// 删除群公告
//
// 参数：
//   - groupID: 群 ID
//   - noticeID: 公告 ID
func (api *API) DeleteGroupNotice(groupID int, noticeID string) (*models.Response[any], error) {
	req := models.DeleteGroupNoticeRequest{
		GroupID:  groupID,
		NoticeID: noticeID,
	}
	return client.Call[models.DeleteGroupNoticeRequest, models.Response[any]](api.client, "_del_group_notice", req)
}

// 获取群公告
//
// 参数：
//   - groupID: 群 ID
func (api *API) GetGroupNotice(groupID int) (*models.GetGroupNoticeResponse, error) {
	req := models.GetGroupNoticeRequest{
		GroupID: groupID,
	}
	return client.Call[models.GetGroupNoticeRequest, models.GetGroupNoticeResponse](api.client, "_get_group_notice", req)

}

// 获取 AI 语音
//
// 参数：
//   - character: 声色 ID
//   - groupID: 群 ID
//   - text: 语音内容
//   - chatType: 聊天类型（如 1 表示群聊）
func (api *API) GetAIRecord(character string, groupID int, text string, chatType int) (*models.GetAIRecordResponse, error) {
	req := models.GetAIRecordRequest{
		Character: character,
		GroupID:   groupID,
		Text:      text,
		ChatType:  chatType,
	}
	return client.Call[models.GetAIRecordRequest, models.GetAIRecordResponse](api.client, "get_ai_record", req)

}

// 获取群荣耀信息
//
// 参数：
//   - groupID: 群 ID
//   - honorType: 荣耀类型（如 "all", "talkative", "active" 等）
func (api *API) GetGroupHonorInfo(groupID int, honorType string) (*models.GetGroupHonorInfoResponse, error) {
	req := models.GetGroupHonorInfoRequest{
		GroupID: groupID,
		Type:    honorType,
	}
	return client.Call[models.GetGroupHonorInfoRequest, models.GetGroupHonorInfoResponse](api.client, "get_group_honor_info", req)

}

// 设置群管理员
//
// 参数：
//   - groupID: 群 ID
//   - userID: 用户 ID
//   - enable: true 设置为管理员，false 取消
func (api *API) SetGroupAdmin(groupID int, userID int, enable bool) (*models.Response[any], error) {
	req := models.SetGroupAdminRequest{
		GroupID: groupID,
		UserID:  userID,
		Enable:  enable,
	}
	return client.Call[models.SetGroupAdminRequest, models.Response[any]](api.client, "set_group_admin", req)
}

// 设置群成员禁言
//
// 参数：
//   - userID: 用户 ID
//   - groupID: 群 ID
//   - duration: 禁言时长（单位：秒）
func (api *API) SetGroupBan(userID int, groupID int, duration int) (*models.Response[any], error) {
	req := models.SetGroupBanRequest{
		UserID:   userID,
		GroupID:  groupID,
		Duration: duration,
	}
	return client.Call[models.SetGroupBanRequest, models.Response[any]](api.client, "set_group_ban", req)
}

// 设置群 Bot 发言状态
//
// 参数：
//   - groupID: 群 ID
//   - botID: Bot ID
//   - enable: 0 禁用，1 启用
func (api *API) SetGroupBotStatus(groupID int, botID int, enable int) (*models.SetGroupBotStatusResponse, error) {
	req := models.SetGroupBotStatusRequest{
		GroupID: groupID,
		BotID:   botID,
		Enable:  enable,
	}
	return client.Call[models.SetGroupBotStatusRequest, models.SetGroupBotStatusResponse](api.client, "set_group_bot_status", req)

}

// 调用群机器人回调
//
// 参数：
//   - groupID: 群 ID
//   - botID: Bot ID
//   - data1: 回调参数 1
//   - data2: 回调参数 2
func (api *API) SendGroupBotCallback(groupID int, botID int, data1 string, data2 string) (*models.SendGroupBotCallbackResponse, error) {
	req := models.SendGroupBotCallbackRequest{
		GroupID: groupID,
		BotID:   botID,
		Data1:   data1,
		Data2:   data2,
	}
	return client.Call[models.SendGroupBotCallbackRequest, models.SendGroupBotCallbackResponse](api.client, "send_group_bot_callback", req)

}

// 设置群名片
//
// 参数：
//   - userID: 用户 ID
//   - groupID: 群 ID
//   - card: 名片内容
func (api *API) SetGroupCard(userID int, groupID int, card string) (*models.Response[any], error) {
	req := models.SetGroupCardRequest{
		UserID:  userID,
		GroupID: groupID,
		Card:    card,
	}
	return client.Call[models.SetGroupCardRequest, models.Response[any]](api.client, "set_group_card", req)
}

// 踢出群成员
//
// 参数：
//   - userID: 用户 ID
//   - groupID: 群 ID
//   - rejectAddRequest: 是否拒绝再次加群
func (api *API) KickGroupMember(userID int, groupID int, rejectAddRequest bool) (*models.Response[any], error) {
	req := models.KickGroupMemberRequest{
		UserID:           userID,
		GroupID:          groupID,
		RejectAddRequest: rejectAddRequest,
	}
	return client.Call[models.KickGroupMemberRequest, models.Response[any]](api.client, "set_group_kick", req)
}

// 退出群（可解散）
//
// 参数：
//   - groupID: 群 ID
//   - isDismiss: 是否解散（仅群主有效）
func (api *API) LeaveGroup(groupID int, isDismiss bool) (*models.Response[any], error) {
	req := models.LeaveGroupRequest{
		GroupID:   groupID,
		IsDismiss: isDismiss,
	}
	return client.Call[models.LeaveGroupRequest, models.Response[any]](api.client, "set_group_leave", req)
}

// 发送群公告
//
// 参数：
//   - groupID: 群 ID
//   - content: 公告文本内容
//   - image: 公告图片链接（可选）
func (api *API) SendGroupNotice(groupID int, content string, image string) (*models.Response[any], error) {
	req := models.SendGroupNoticeRequest{
		GroupID: groupID,
		Content: content,
		Image:   image,
	}
	return client.Call[models.SendGroupNoticeRequest, models.Response[any]](api.client, "_send_group_notice", req)
}

// 设置群名称
//
// 参数：
//   - groupID: 群 ID
//   - groupName: 新的群名称
func (api *API) SetGroupName(groupID int, groupName string) (*models.Response[any], error) {
	req := models.SetGroupNameRequest{
		GroupID:   groupID,
		GroupName: groupName,
	}
	return client.Call[models.SetGroupNameRequest, models.Response[any]](api.client, "set_group_name", req)
}

// 设置全体禁言
//
// 参数：
//   - groupID: 群 ID
//   - enable: 是否开启禁言（true 表示开启）
func (api *API) SetGroupWholeBan(groupID int, enable bool) (*models.Response[any], error) {
	req := models.SetGroupWholeBanRequest{
		GroupID: groupID,
		Enable:  enable,
	}
	return client.Call[models.SetGroupWholeBanRequest, models.Response[any]](api.client, "set_group_whole_ban", req)
}

// 设置群头像
//
// 参数：
//   - groupID: 群 ID
//   - file: 头像链接或 Base64 图片
func (api *API) SetGroupAvatar(groupID int, file string) (*models.Response[any], error) {
	req := models.SetGroupAvatarRequest{
		GroupID: groupID,
		File:    file,
	}
	return client.Call[models.SetGroupAvatarRequest, models.Response[any]](api.client, "set_group_avatar", req)
}

// 设置群表情回复（消息表情）
//
// 参数：
//   - groupID: 群 ID
//   - messageID: 消息 ID
//   - code: 表情代码
//   - isAdd: 是否添加（true 表示添加，false 表示移除）
func (api *API) SetEmojiReaction(groupID int, messageID int, code string, isAdd bool) (*models.Response[any], error) {
	req := models.SetEmojiReactionRequest{
		GroupID:   groupID,
		MessageID: messageID,
		Code:      code,
		IsAdd:     isAdd,
	}
	return client.Call[models.SetEmojiReactionRequest, models.Response[any]](api.client, "set_emoji_reaction", req)
}

// 设置群专属头衔
//
// 参数：
//   - groupID: 群 ID
//   - userID: 用户 ID
//   - specialTitle: 头衔名称
//   - duration: 头衔有效期（单位：秒，0 为永久）
func (api *API) SetGroupSpecialTitle(groupID int, userID int, specialTitle string, duration int) (*models.Response[any], error) {
	req := models.SetGroupSpecialTitleRequest{
		GroupID:      groupID,
		UserID:       userID,
		SpecialTitle: specialTitle,
		Duration:     duration,
	}
	return client.Call[models.SetGroupSpecialTitleRequest, models.Response[any]](api.client, "set_group_special_title", req)
}
