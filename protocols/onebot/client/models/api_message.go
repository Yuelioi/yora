package models

import (
	"yora/internal/message"
)

type MessageRequest struct {
	Message message.Message `json:"message"`
	// 消息类型
	MessageType string `json:"message_type"`
	// 用户UIN
	UserID *int `json:"user_id,omitempty"`
	// 是否不解析CQ码
	AutoEscape *bool `json:"auto_escape,omitempty"`
	// 群UIN
	GroupID *int `json:"group_id,omitempty"`
}

// MessageSegment 富文本消息段列表
// 支持多种消息类型：At消息段、表情消息段、图片消息段、回复消息段、文本消息段、
// 单个消息段列表、骰子消息段、转发消息段、Json消息段、定位消息段、长消息消息段、
// 商城表情消息段、音乐消息段、戳一戳消息段、语音消息段、猜拳消息段、视频消息段、
// Markdown消息段列表、Markdown消息段、按钮消息段（只能跟在Markdown消息段下）

// MessageSegmentData 消息段数据
type MessageSegmentData struct {
	// 显示的文本
	Name *string `json:"name,omitempty"`
	// 用户UIN
	Qq *string `json:"qq,omitempty"`
	// 表情ID / 消息ID / 转发ID / 长消息ID / 戳一戳ID
	ID *string `json:"id,omitempty"`
	// 是否大表情
	Large *bool `json:"large,omitempty"`
	// 图片链接，支持http/https/file/base64 / file链接，支持http/https/file/base64 / 视频链接，支持http/https/file/base64
	File *string `json:"file,omitempty"`
	// 图片名称
	Filename *string `json:"filename,omitempty"`
	// 图片子类型
	SubType *string `json:"subType,omitempty"`
	// 图片说明 / 表情说明
	Summary *string `json:"summary,omitempty"`
	// 图片链接 / 表情URL / 跳转URL / 语音链接 / 视频链接
	URL *string `json:"url,omitempty"`
	// 文本内容
	Text *string `json:"text,omitempty"`
	// Json数据
	Data *string `json:"data,omitempty"`
	// 纬度
	Lat *string `json:"lat,omitempty"`
	// 经度
	Lon *string `json:"lon,omitempty"`
	// 标题
	Title *string `json:"title,omitempty"`
	// 表情ID
	EmojiID *string `json:"emoji_id,omitempty"`
	// 表情包ID
	EmojiPackageID *int64 `json:"emoji_package_id,omitempty"`
	// 表情Key
	Key *string `json:"key,omitempty"`
	// 音乐URL
	Audio *string `json:"audio,omitempty"`
	// 图片
	Image *string `json:"image,omitempty"`
	// 音乐类型 / 戳一戳类型
	Type *string `json:"type,omitempty"`
	// 应用ID
	Appid *string `json:"appid,omitempty"`
	// 应用包名
	PackageName *string `json:"package_name,omitempty"`
	// 应用签名
	Sign *string `json:"sign,omitempty"`
	// 戳一戳强度
	Strength *string `json:"strength,omitempty"`
	// 文件Hash
	FileHash *string `json:"file_hash,omitempty"`
	// 文件ID
	FileID *string `json:"file_id,omitempty"`
	// 文件名
	FileName *string `json:"file_name,omitempty"`
}

// SingleMessageSegment 单个消息段
// 支持多种消息类型：骰子消息段、转发消息段、Json消息段、定位消息段、长消息消息段、
// 商城表情消息段、音乐消息段、戳一戳消息段、语音消息段、猜拳消息段、视频消息段、
// At消息段、表情消息段、图片消息段、回复消息段（只能书写一个）、文本消息段、
// Markdown消息段、按钮消息段（只能跟在Markdown消息段下）、文件消息段

// Type 消息段类型枚举
type Type string

const (
	At       Type = "at"       // At消息
	Dict     Type = "dict"     // 字典
	Face     Type = "face"     // 表情
	Forward  Type = "forward"  // 转发
	Image    Type = "image"    // 图片
	JSON     Type = "json"     // JSON数据
	Keyboard Type = "keyboard" // 键盘
	Location Type = "location" // 位置
	Longmsg  Type = "longmsg"  // 长消息
	Mface    Type = "mface"    // 商城表情
	Music    Type = "music"    // 音乐
	Poke     Type = "poke"     // 戳一戳
	Record   Type = "record"   // 语音
	Reply    Type = "reply"    // 回复
	Rps      Type = "rps"      // 猜拳
	Text     Type = "text"     // 文本
	Video    Type = "video"    // 视频
	// Markdown Type = "markdown" // Markdown

)

// EssenceMessage 精华消息相关结构体

// DeleteEssenceMessageRequest 删除精华消息请求
type DeleteEssenceMessageRequest struct {
	MessageID int `json:"message_id"` // 消息ID
}

// SetEssenceMessageRequest 设置精华消息请求
type SetEssenceMessageRequest struct {
	MessageID int `json:"message_id"` // 消息ID
}

// GetEssenceMessageListRequest 获取精华消息列表请求
type GetEssenceMessageListRequest struct {
	GroupID int `json:"group_id"` // 群ID
}

// GetEssenceMessageListResponse 获取精华消息列表响应
type GetEssenceMessageListResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    []struct {
		SenderID     int    `json:"sender_id"`     // 发送者ID
		SenderNick   string `json:"sender_nick"`   // 发送者昵称
		SenderTime   int    `json:"sender_time"`   // 发送时间
		OperatorID   int    `json:"operator_id"`   // 操作者ID
		OperatorNick string `json:"operator_nick"` // 操作者昵称
		OperatorTime int    `json:"operator_time"` // 操作时间
		MessageID    int    `json:"message_id"`    // 消息ID
		Content      []struct {
			Type string `json:"type"` // 消息段类型
			Data struct {
				Name string `json:"name"` // 名称
				Qq   string `json:"qq"`   // QQ号
			} `json:"data"` // 消息段数据
		} `json:"content"` // 消息内容
	} `json:"data"` // 精华消息列表
}

// MessageOperations 消息操作相关结构体

// RecallMessageRequest 撤回消息请求
type RecallMessageRequest struct {
	MessageID int `json:"message_id"` // 消息ID
}

// MarkMessageAsReadRequest 标记消息为已读请求
type MarkMessageAsReadRequest struct {
	MessageID int `json:"message_id"` // 消息ID
}

// GetMessageRequest 获取消息请求
type GetMessageRequest struct {
	MessageID int `json:"message_id"` // 消息ID
}

// GetMessageResponse 获取消息响应
type GetMessageResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		Time        int    `json:"time"`         // 发送时间
		MessageType string `json:"message_type"` // 消息类型
		MessageID   int    `json:"message_id"`   // 消息ID
		RealID      int    `json:"real_id"`      // 真实ID
		Sender      struct {
			UserID   int    `json:"user_id"`  // 用户ID
			Nickname string `json:"nickname"` // 昵称
			Sex      string `json:"sex"`      // 性别
			Age      int    `json:"age"`      // 年龄
		} `json:"sender"` // 发送者信息
		Message map[string]any `json:"message"` // 消息内容
	} `json:"data"` // 消息数据
}

// Poke 戳一戳相关结构体

// PrivatePokeRequest 私聊戳一戳请求
type PrivatePokeRequest struct {
	UserID int `json:"user_id"` // 用户ID
}

// GroupPokeRequest 群里戳一戳请求
type GroupPokeRequest struct {
	GroupID int `json:"group_id"` // 群ID
	UserID  int `json:"user_id"`  // 用户ID
}

// ForwardMessage 合并转发消息相关结构体

// GetForwardMessageRequest 获取合并转发消息请求
type GetForwardMessageRequest struct {
	ID string `json:"id"` // 合并转发消息ID
}

// GetForwardMessageResponse 获取合并转发消息响应
type GetForwardMessageResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		Message []struct {
			Type string `json:"type"` // 消息段类型
			Data struct {
				UserID   string `json:"user_id"`  // 用户ID
				Nickname string `json:"nickname"` // 昵称
				Content  []struct {
					Type string `json:"type"` // 消息段类型
					Data struct {
						Name string `json:"name"` // 名称
						Qq   string `json:"qq"`   // QQ号
					} `json:"data"` // 消息段数据
				} `json:"content"` // 消息内容
			} `json:"data"` // 消息段数据
		} `json:"message"` // 消息列表
	} `json:"data"` // 转发消息数据
}

// ConstructForwardMessageRequest 构造合并转发消息请求
type ConstructForwardMessageRequest struct {
	Messages []struct {
		Type string `json:"type"` // 消息段类型
		Data struct {
			UserID   string `json:"user_id"`  // 用户ID
			Nickname string `json:"nickname"` // 昵称
			Content  []struct {
				Type string `json:"type"` // 消息段类型
				Data struct {
					Name string `json:"name"` // 名称
					Qq   string `json:"qq"`   // QQ号
				} `json:"data"` // 消息段数据
			} `json:"content"` // 消息内容
		} `json:"data"` // 消息段数据
	} `json:"messages"` // 消息列表
}

// SendGroupForwardMessageRequest 发送群聊合并转发消息请求
type SendGroupForwardMessageRequest struct {
	GroupID  int `json:"group_id"` // 群ID
	Messages []struct {
		Type string `json:"type"` // 消息段类型
		Data struct {
			UserID   string `json:"user_id"`  // 用户ID
			Nickname string `json:"nickname"` // 昵称
			Content  []struct {
				Type string `json:"type"` // 消息段类型
				Data struct {
					Name string `json:"name"` // 名称
					Qq   string `json:"qq"`   // QQ号
				} `json:"data"` // 消息段数据
			} `json:"content"` // 消息内容
		} `json:"data"` // 消息段数据
	} `json:"messages"` // 消息列表
}

// SendGroupForwardMessageResponse 发送群聊合并转发消息响应
type SendGroupForwardMessageResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		MessageID int    `json:"message_id"` // 消息ID
		ForwardID string `json:"forward_id"` // 转发消息ID
	} `json:"data"` // 响应数据
}

// SendPrivateForwardMessageRequest 发送私聊合并转发消息请求
type SendPrivateForwardMessageRequest struct {
	UserID   int `json:"user_id"` // 用户ID
	Messages []struct {
		Type string `json:"type"` // 消息段类型
		Data struct {
			UserID   string `json:"user_id"`  // 用户ID
			Nickname string `json:"nickname"` // 昵称
			Content  []struct {
				Type string `json:"type"` // 消息段类型
				Data struct {
					Name string `json:"name"` // 名称
					Qq   string `json:"qq"`   // QQ号
				} `json:"data"` // 消息段数据
			} `json:"content"` // 消息内容
		} `json:"data"` // 消息段数据
	} `json:"messages"` // 消息列表
}

// ChatHistory 聊天记录相关结构体

// GetFriendChatHistoryRequest 获取好友历史聊天记录请求
type GetFriendChatHistoryRequest struct {
	UserID    int `json:"user_id"`    // 用户ID
	MessageID int `json:"message_id"` // 消息ID（起始位置）
	Count     int `json:"count"`      // 获取数量
}

// GetFriendChatHistoryResponse 获取好友历史聊天记录响应
type GetFriendChatHistoryResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		Messages []struct {
			MessageType string         `json:"message_type"` // 消息类型
			SubType     string         `json:"sub_type"`     // 子类型
			MessageID   int            `json:"message_id"`   // 消息ID
			UserID      int            `json:"user_id"`      // 用户ID
			Message     map[string]any `json:"message"`      // 消息内容
			RawMessage  string         `json:"raw_message"`  // 原始消息
			Font        int            `json:"font"`         // 字体
			Sender      struct {
				UserID   int    `json:"user_id"`  // 用户ID
				Nickname string `json:"nickname"` // 昵称
				Sex      string `json:"sex"`      // 性别
				Age      int    `json:"age"`      // 年龄
			} `json:"sender"` // 发送者信息
			TargetID int `json:"target_id"` // 目标ID
		} `json:"messages"` // 消息列表
	} `json:"data"` // 聊天记录数据
}

// GetGroupChatHistoryRequest 获取群历史聊天记录请求
type GetGroupChatHistoryRequest struct {
	GroupID   int    `json:"group_id"`   // 群ID
	MessageID string `json:"message_id"` // 消息ID（起始位置）
	Count     int    `json:"count"`      // 获取数量
}

// GetGroupChatHistoryResponse 获取群历史聊天记录响应
type GetGroupChatHistoryResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		Messages []struct {
			MessageTyp string         `json:"message_typ"` // 消息类型
			SubType    string         `json:"sub_type"`    // 子类型
			MessageID  int            `json:"message_id"`  // 消息ID
			GroupID    int            `json:"group_id"`    // 群ID
			UserID     int            `json:"user_id"`     // 用户ID
			Anonymous  any            `json:"anonymous"`   // 匿名信息
			Message    map[string]any `json:"message"`     // 消息内容
			RawMessage string         `json:"raw_message"` // 原始消息
			Font       string         `json:"font"`        // 字体
			Sender     struct {
				UserID   int    `json:"user_id"`  // 用户ID
				Nickname string `json:"nickname"` // 昵称
				Card     string `json:"card"`     // 群名片
				Sex      string `json:"sex"`      // 性别
				Age      int    `json:"age"`      // 年龄
				Area     string `json:"area"`     // 地区
				Level    string `json:"level"`    // 等级
				Role     string `json:"role"`     // 角色
				Title    string `json:"title"`    // 头衔
			} `json:"sender"` // 发送者信息
		} `json:"messages"` // 消息列表
	} `json:"data"` // 聊天记录数据
}

// SendMessage 发送消息相关结构体

// SendMessageRequest 发送消息请求

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		MessageID int `json:"message_id"` // 消息ID
	} `json:"data"` // 响应数据
}

// SendPrivateMessageResponse 发送私聊消息响应
type SendPrivateMessageResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		MessageID int `json:"message_id"` // 消息ID
	} `json:"data"` // 响应数据
}

// SendGroupMessageResponse 发送群消息响应
type SendGroupMessageResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		MessageID int `json:"message_id"` // 消息ID
	} `json:"data"` // 响应数据
}

// AIVoice AI语音相关结构体

// SendGroupAIVoiceRequest 发送群AI语音请求
type SendGroupAIVoiceRequest struct {
	Character string `json:"character"` // 角色
	GroupID   int    `json:"group_id"`  // 群ID
	Text      string `json:"text"`      // 文本内容
	ChatType  int    `json:"chat_type"` // 聊天类型
}

// SendGroupAIVoiceResponse 发送群AI语音响应
type SendGroupAIVoiceResponse struct {
	Status  string `json:"status"`  // 状态
	Retcode int    `json:"retcode"` // 返回码
	Data    struct {
		MessageID int `json:"message_id"` // 消息ID
	} `json:"data"` // 响应数据
}
