package models

// HandleFriendRequestRequest 处理加好友请求
type HandleFriendRequestRequest struct {
	Flag    string `json:"flag"`    // 加好友请求的标识符
	Approve bool   `json:"approve"` // 是否同意请求
	Reason  string `json:"reason"`  // 拒绝理由（可选）
}

// GroupRequest 群请求相关结构体

// HandleGroupRequestRequest 处理加群请求/邀请
type HandleGroupRequestRequest struct {
	Flag    string `json:"flag"`    // 加群请求/邀请的标识符
	Approve bool   `json:"approve"` // 是否同意请求/邀请
	Reason  string `json:"reason"`  // 拒绝理由（可选）
}
