package models

// 删除群公告请求
type DeleteGroupNoticeRequest struct {
	GroupID  int    `json:"group_id"`  // 群 ID
	NoticeID string `json:"notice_id"` // 公告 ID
}

// 获取群公告请求
type GetGroupNoticeRequest struct {
	GroupID int `json:"group_id"` // 群 ID
}

// 获取群公告响应
type GetGroupNoticeResponse = Response[[]struct {
	NoticeID    string `json:"notice_id"`    // 公告 ID
	SenderID    int    `json:"sender_id"`    // 发布者 ID
	PublishTime int    `json:"publish_time"` // 发布时间戳
	Message     struct {
		Text   string `json:"text"` // 公告文本内容
		Images []struct {
			ID     string `json:"id"`     // 图片 ID
			Height string `json:"height"` // 图片高度
			Width  string `json:"width"`  // 图片宽度
		} `json:"images"`
	} `json:"message"`
}]

// 获取 AI 语音请求
type GetAIRecordRequest struct {
	Character string `json:"character"` // 声色 ID
	GroupID   int    `json:"group_id"`  // 群 ID
	Text      string `json:"text"`      // 语音内容
	ChatType  int    `json:"chat_type"` // 聊天类型
}

// 获取 AI 语音响应（为 null）
type GetAIRecordResponse = Response[any]

// 获取群荣耀信息请求
type GetGroupHonorInfoRequest struct {
	GroupID int    `json:"group_id"` // 群 ID
	Type    string `json:"type"`     // 荣耀类型（如 "all"）
}

// 获取群荣耀信息响应
type GetGroupHonorInfoResponse = Response[struct {
	GroupID          int `json:"group_id"` // 群 ID
	CurrentTalkative struct {
		DayCount   int    `json:"day_count"`   // 天数
		Avatar     string `json:"avatar"`      // 头像链接
		AvatarSize int    `json:"avatar_size"` // 头像尺寸
		UserID     int    `json:"user_id"`     // 用户 ID
		Nickname   string `json:"nickname"`    // 昵称
	} `json:"current_talkative"`
	TalkativeList []struct {
		Avatar      string `json:"avatar"`
		BtnText     string `json:"btnText"`
		Text        string `json:"text"`
		UserID      int    `json:"user_id"`
		Nickname    string `json:"nickname"`
		Description string `json:"description"`
	} `json:"talkative_list"`
	LegendList []struct {
		Avatar      string `json:"avatar"`
		BtnText     string `json:"btnText"`
		Icon        int    `json:"icon"`
		UserID      int    `json:"user_id"`
		Nickname    string `json:"nickname"`
		Description string `json:"description"`
	} `json:"legend_list"`
	StrongNewbieList []string `json:"strong_newbie_list"`
	EmotionList      []struct {
		Avatar      string `json:"avatar"`
		BtnText     string `json:"btnText"`
		Icon        int    `json:"icon"`
		UserID      int    `json:"user_id"`
		Nickname    string `json:"nickname"`
		Description string `json:"description"`
	} `json:"emotion_list"`
}]

// 设置群管理员请求
type SetGroupAdminRequest struct {
	GroupID int  `json:"group_id"` // 群 ID
	UserID  int  `json:"user_id"`  // 用户 ID
	Enable  bool `json:"enable"`   // 是否设置为管理员
}

// 设置群成员禁言请求
type SetGroupBanRequest struct {
	UserID   int `json:"user_id"`  // 用户 ID
	GroupID  int `json:"group_id"` // 群 ID
	Duration int `json:"duration"` // 禁言时长（秒）
}

// 设置群 Bot 发言状态请求
type SetGroupBotStatusRequest struct {
	GroupID int `json:"group_id"` // 群 ID
	BotID   int `json:"bot_id"`   // Bot ID
	Enable  int `json:"enable"`   // 启用状态（0 或 1）
}

// 群 Bot 发言状态响应
type SetGroupBotStatusResponse = Response[any]

// 调用群机器人回调请求
type SendGroupBotCallbackRequest struct {
	GroupID int    `json:"group_id"` // 群 ID
	BotID   int    `json:"bot_id"`   // Bot ID
	Data1   string `json:"data_1"`   // 回调数据 1
	Data2   string `json:"data_2"`   // 回调数据 2
}

// 群机器人回调响应
type SendGroupBotCallbackResponse = Response[any]

// 设置群名片请求
type SetGroupCardRequest struct {
	UserID  int    `json:"user_id"`  // 用户 ID
	GroupID int    `json:"group_id"` // 群 ID
	Card    string `json:"card"`     // 群名片内容
}

// 踢出群成员请求
type KickGroupMemberRequest struct {
	UserID           int  `json:"user_id"`            // 用户 ID
	GroupID          int  `json:"group_id"`           // 群 ID
	RejectAddRequest bool `json:"reject_add_request"` // 是否拒绝再次加群
}

// 退出群请求
type LeaveGroupRequest struct {
	GroupID   int  `json:"group_id"`   // 群 ID
	IsDismiss bool `json:"is_dismiss"` // 是否解散群（仅群主有效）
}

// 发送群公告请求
type SendGroupNoticeRequest struct {
	GroupID int    `json:"group_id"` // 群 ID
	Content string `json:"content"`  // 公告内容
	Image   string `json:"image"`    // 公告图片链接
}

// 设置群名称请求
type SetGroupNameRequest struct {
	GroupID   int    `json:"group_id"`   // 群 ID
	GroupName string `json:"group_name"` // 群名称
}

// 设置全体禁言请求
type SetGroupWholeBanRequest struct {
	GroupID int  `json:"group_id"` // 群 ID
	Enable  bool `json:"enable"`   // 是否启用禁言
}

// 设置群头像请求
type SetGroupPortraitRequest struct {
	GroupID int    `json:"group_id"` // 群 ID
	File    string `json:"file"`     // 头像文件链接
}

// 设置群表情回复请求
type SetEmojiReactionRequest struct {
	GroupID   int    `json:"group_id"`   // 群 ID
	MessageID int    `json:"message_id"` // 消息 ID
	Code      string `json:"code"`       // 表情代码
	IsAdd     bool   `json:"is_add"`     // 是否添加表情
}

// 设置群专属头衔请求
type SetGroupSpecialTitleRequest struct {
	GroupID      int    `json:"group_id"`      // 群 ID
	UserID       int    `json:"user_id"`       // 用户 ID
	SpecialTitle string `json:"special_title"` // 专属头衔
	Duration     int    `json:"duration"`      // 有效期（秒）
}
