package models

// GroupNotice 群公告相关结构体

// DeleteGroupNoticeRequest 删除群公告请求
type DeleteGroupNoticeRequest struct {
	GroupID  int    `json:"group_id"`  // 群ID
	NoticeID string `json:"notice_id"` // 公告ID
}

// GetGroupNoticeRequest 获取群公告请求
type GetGroupNoticeRequest struct {
	GroupID int `json:"group_id"` // 群ID
}

// GetGroupNoticeResponse 获取群公告响应
type GetGroupNoticeResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    []struct {
		NoticeID    string `json:"notice_id"`    // 公告ID
		SenderID    int    `json:"sender_id"`    // 发送者ID
		PublishTime int    `json:"publish_time"` // 发布时间
		Message     struct {
			Text   string `json:"text"` // 公告文本内容
			Images []struct {
				ID     string `json:"id"`     // 图片ID
				Height string `json:"height"` // 图片高度
				Width  string `json:"width"`  // 图片宽度
			} `json:"images"` // 公告图片列表
		} `json:"message"` // 公告消息内容
	} `json:"data"` // 公告列表
}

// SendGroupNoticeRequest 发送群公告请求
type SendGroupNoticeRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	Content string `json:"content"`  // 公告内容
	Image   string `json:"image"`    // 公告图片链接，支持 http/https/file/base64
}

// GroupAdmin 群管理员相关结构体

// SetGroupAdminRequest 设置群管理员请求
type SetGroupAdminRequest struct {
	GroupID int  `json:"group_id"` // 群ID
	UserID  int  `json:"user_id"`  // 用户ID
	Enable  bool `json:"enable"`   // 是否设置为管理员
}

// GroupMute 群禁言相关结构体

// SetGroupBanRequest 设置群成员禁言请求
type SetGroupBanRequest struct {
	UserID   int `json:"user_id"`  // 用户ID
	GroupID  int `json:"group_id"` // 群ID
	Duration int `json:"duration"` // 禁言时长（秒）
}

// SetGroupWholeBanRequest 设置群全体禁言请求
type SetGroupWholeBanRequest struct {
	GroupID int  `json:"group_id"` // 群ID
	Enable  bool `json:"enable"`   // 是否开启全体禁言
}

// GroupMember 群成员相关结构体

// SetGroupCardRequest 设置群名片请求
type SetGroupCardRequest struct {
	UserID  int    `json:"user_id"`  // 用户ID
	GroupID int    `json:"group_id"` // 群ID
	Card    string `json:"card"`     // 群名片内容
}

// SetGroupSpecialTitleRequest 设置群组专属头衔请求
type SetGroupSpecialTitleRequest struct {
	GroupID      int    `json:"group_id"`      // 群ID
	UserID       int    `json:"user_id"`       // 用户ID
	SpecialTitle string `json:"special_title"` // 专属头衔
	Duration     int    `json:"duration"`      // 头衔有效期（秒）
}

// KickGroupMemberRequest 踢出群成员请求
type KickGroupMemberRequest struct {
	UserID           int  `json:"user_id"`            // 用户ID
	GroupID          int  `json:"group_id"`           // 群ID
	RejectAddRequest bool `json:"reject_add_request"` // 是否拒绝此人的加群请求
}

// LeaveGroupRequest 退出群组请求
type LeaveGroupRequest struct {
	GroupID   int  `json:"group_id"`   // 群ID
	IsDismiss bool `json:"is_dismiss"` // 是否解散群组（群主有效）
}

// GroupSettings 群设置相关结构体

// SetGroupNameRequest 设置群名请求
type SetGroupNameRequest struct {
	GroupID   int    `json:"group_id"`   // 群ID
	GroupName string `json:"group_name"` // 群名称
}

// SetGroupAvatarRequest 设置群头像请求
type SetGroupAvatarRequest struct {
	GroupID int    `json:"group_id"` // 群ID
	File    string `json:"file"`     // 头像文件链接，支持 http/https/file/base64
}

// GroupMessage 群消息相关结构体

// SetEmojiReactionRequest 表情回复操作请求
type SetEmojiReactionRequest struct {
	GroupID   int    `json:"group_id"`   // 群ID
	MessageID int    `json:"message_id"` // 消息ID
	Code      string `json:"code"`       // 表情代码
	IsAdd     bool   `json:"is_add"`     // 是否添加表情回复（true: 添加, false: 移除）
}
