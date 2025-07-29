package chat

// GroupMemberInfo QQ群成员信息类型定义
type GroupMemberInfo struct {
	GroupID         int    `json:"group_id"`
	UserID          int    `json:"user_id"`
	Nickname        string `json:"nickname"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int    `json:"age"`
	Area            string `json:"area"`
	JoinTime        int    `json:"join_time"`
	LastSentTime    int    `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	TitleExpireTime int    `json:"title_expire_time"`
	CardChangeable  bool   `json:"card_changeable"`
}

// SimplifiedMessage 极简版群消息结构
type SimplifiedMessage struct {
	UserID   int    `json:"user_id"`  // 发送者QQ号
	Nickname string `json:"nickname"` // 群昵称
	Text     string `json:"text"`     // 合并后的消息文本
	Time     string `json:"time"`     // 消息时间
}

// GroupMessageHistory 群消息历史记录
type GroupMessageHistory struct {
	Messages []SimplifiedMessage `json:"messages"`
}

// RelationshipInfo 关系信息
type RelationshipInfo struct {
	Status       string `json:"status"`
	Attitude     string `json:"attitude"`
	Mode         string `json:"mode"`
	Relationship string `json:"relationship"`
}

// OpenAI API 请求结构
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

// 原始消息结构
type RawMessage struct {
	UserID  int                    `json:"user_id"`
	Time    int64                  `json:"time"`
	Message []MessageItem          `json:"message"`
	Sender  map[string]interface{} `json:"sender"`
}

type MessageItem struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
