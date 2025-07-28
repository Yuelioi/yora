package models

// 获取自定义表情请求
type FetchCustomFaceRequest struct{}

// 获取自定义表情响应（表情 URL 列表）
type FetchCustomFaceResponse = Response[[]string]

// 获取商城表情 key 请求
type FetchMfaceKeyRequest struct {
	Emoji_IDs []string `json:"emoji_ids"` // 表情 ID 列表
}

// 获取商城表情 key 响应
type FetchMfaceKeyResponse = Response[[]string]

// 加入好友表情接龙请求
type JoinFriendEmojiChainRequest struct {
	UserID    int `json:"user_id"`    // 用户 ID
	MessageID int `json:"message_id"` // 消息 ID
	EmojiID   int `json:"emoji_id"`   // 表情 ID
}

// 加入好友表情接龙响应
type JoinFriendEmojiChainResponse = Response[any]

// 获取群 Ai 声色请求
type GetAICharactersRequest struct {
	GroupID  int `json:"group_id"`  // 群 ID
	ChatType int `json:"chat_type"` // 聊天类型（如 1 表示群聊）
}

// 获取群 Ai 声色响应
type GetAICharactersResponse = Response[[]struct {
	Type       string `json:"type"` // 类型
	Characters []struct {
		CharacterID   string `json:"character_id"`   // 声色 ID
		CharacterName string `json:"character_name"` // 声色名称
		PreviewURL    string `json:"preview_url"`    // 试听地址
	} `json:"characters"`
}]

// 获取 Cookies 请求
type GetCookiesRequest struct {
	Domain string `json:"domain"` // 域名
}

// 获取 Cookies 响应
type GetCookiesResponse = Response[struct {
	Cookies string `json:"cookies"` // Cookies 字符串
}]

// 获取 QQ 接口凭证请求
type GetCredentialsRequest struct {
	Domain string `json:"domain"` // 域名
}

// 获取 QQ 接口凭证响应
type GetCredentialsResponse = Response[struct {
	Cookies   string `json:"cookies"`    // Cookies
	CsrfToken string `json:"csrf_token"` // CSRF Token
}]

// 获取 CSRF Token 请求
type GetCSRFTokenRequest struct{}

// 获取 CSRF Token 响应
type GetCSRFTokenResponse = Response[struct {
	Token int `json:"token"` // CSRF Token
}]

// 加入群聊表情接龙请求
type JoinGroupEmojiChainRequest struct {
	GroupID   int `json:"group_id"`   // 群 ID
	MessageID int `json:"message_id"` // 消息 ID
	EmojiID   int `json:"emoji_id"`   // 表情 ID
}

// 加入群聊表情接龙响应
type JoinGroupEmojiChainResponse = Response[any]

// OCR 图像识别请求
type OCRImageRequest struct {
	Image string `json:"image"` // 图像地址（或 base64）
}

// OCR 图像识别响应
type OCRImageResponse = Response[struct {
	Texts []struct {
		Text        string `json:"text"`       // 识别文本
		Confidence  int    `json:"confidence"` // 置信度
		Coordinates []struct {
			X int `json:"x"` // X 坐标
			Y int `json:"y"` // Y 坐标
		} `json:"coordinates"` // 文本区域坐标
	} `json:"texts"`
	Language string `json:"language"` // 检测语言
}]

// 设置 QQ 头像请求
type SetQQAvatarRequest struct {
	File string `json:"file"` // 图像内容或路径
}

// 设置 QQ 头像响应
type SetQQAvatarResponse = Response[any]

// 点赞用户资料请求
type SendLikeRequest struct {
	UserID int `json:"user_id"` // 用户 ID
	Times  int `json:"times"`   // 点赞次数
}

// 点赞用户资料响应
type SendLikeResponse = Response[any]

// 删除好友请求
type DeleteFriendRequest struct {
	UserID string `json:"user_id"` // 用户 ID
	Block  bool   `json:"block"`   // 是否拉黑
}

// 删除好友响应
type DeleteFriendResponse = Response[any]

// 获取 rkey 响应
type GetRKeyResponse = Response[struct {
	RKeys []struct {
		Type      string `json:"type"`       // 类型（如 private/group）
		RKey      string `json:"rkey"`       // RKey 字符串
		CreatedAt int64  `json:"created_at"` // 创建时间戳
		TTL       int    `json:"ttl"`        // 有效时间（秒）
	} `json:"rkeys"`
}]
