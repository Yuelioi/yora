package models

// 获取好友列表请求
type GetFriendListRequest struct{}

// 获取好友列表响应
type GetFriendListResponse = Response[[]Friend]

// 好友结构
type Friend struct {
	UserID   int    `json:"user_id"`
	QID      string `json:"q_id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
	Group    struct {
		GroupID   int    `json:"group_id"`
		GroupName string `json:"group_name"`
	} `json:"group"`
}

// 获取群信息请求
type GetGroupInfoRequest struct {
	GroupID int  `json:"group_id"`
	NoCache bool `json:"no_cache"`
}

// 获取群信息响应
type GetGroupInfoResponse = Response[GroupInfo]

// 群信息结构
type GroupInfo struct {
	GroupID        int    `json:"group_id"`         // 群号
	GroupName      string `json:"group_name"`       // 群名称
	MaxMemberCount int    `json:"max_member_count"` // 最大成员数
	MemberCount    int    `json:"member_count"`     // 当前成员数
}

// 获取群成员列表请求
type GetGroupMemberListRequest struct {
	GroupID int `json:"group_id"`
}

// 获取群成员列表响应
type GetGroupMemberListResponse = Response[[]GroupMember]

// 群成员结构
type GroupMember struct {
	UserID          int64  `json:"user_id"`           // 用户 ID
	GroupID         int64  `json:"group_id"`          // 群 ID
	Nickname        string `json:"nickname"`          // 昵称
	Card            string `json:"card"`              // 群名片
	CardChangeable  bool   `json:"card_changeable"`   // 是否允许修改群名片
	Sex             string `json:"sex"`               // 性别
	Age             int    `json:"age"`               // 年龄
	Area            string `json:"area"`              // 地区
	JoinTime        int64  `json:"join_time"`         // 加群时间（Unix 时间戳）
	LastSentTime    int64  `json:"last_sent_time"`    // 上次发言时间（Unix 时间戳）
	Level           string `json:"level"`             // 等级
	Role            string `json:"role"`              // 角色（owner/admin/member）
	Title           string `json:"title"`             // 头衔
	TitleExpireTime int64  `json:"title_expire_time"` // 头衔过期时间（Unix 时间戳）
	Unfriendly      bool   `json:"unfriendly"`        // 是否为非好友
}

// 获取群成员信息请求
type GetGroupMemberInfoRequest struct {
	GroupID int  `json:"group_id"`
	UserID  int  `json:"user_id"`
	NoCache bool `json:"no_cache"`
}

// 获取群成员信息响应
type GetGroupMemberInfoResponse = Response[GroupMember]

// 获取群列表请求
type GetGroupListRequest struct {
	NoCache bool `json:"no_cache"`
}

// 获取群列表响应
type GetGroupListResponse = Response[[]GroupInfo]

// 获取登录信息请求
type GetLoginInfoRequest struct{}

// 获取登录信息响应
type GetLoginInfoResponse = Response[LoginInfo]

// 登录信息结构
type LoginInfo struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
}

// 获取状态请求
type GetStatusRequest struct{}

// 获取状态响应
type GetStatusResponse = Response[StatusInfo]

// 状态信息结构
type StatusInfo struct {
	AppInitialized bool `json:"app_initialized"`
	AppEnabled     bool `json:"app_enabled"`
	PluginsGood    bool `json:"plugins_good"`
	AppGood        bool `json:"app_good"`
	Online         bool `json:"online"` // 是否在线
	Good           bool `json:"good"`   // 是否正常
	Memory         int  `json:"memory"`
}

// 获取陌生人信息请求
type GetStrangerInfoRequest struct {
	UserID  int  `json:"user_id"`
	NoCache bool `json:"no_cache"`
}

// 获取陌生人信息响应
type GetStrangerInfoResponse = Response[StrangerInfo]

// 陌生人信息结构
type StrangerInfo struct {
	UserID   int    `json:"user_id"`  // 用户 ID
	Avatar   string `json:"avatar"`   // 头像
	QID      string `json:"q_id"`     // QID
	Nickname string `json:"nickname"` // 昵称
	Sign     string `json:"sign"`     // 个性签名
	Sex      string `json:"sex"`      // 性别
	Age      int    `json:"age"`      // 年龄
	Level    int    `json:"level"`    // 等级
	Status   struct {
		StatusID int    `json:"status_id"`
		FaceID   int    `json:"face_id"`
		Message  string `json:"message"`
	} `json:"status"`
	RegisterTime string `json:"RegisterTime"` // 注册时间
	Business     []struct {
		Type   int    `json:"type"`
		Name   string `json:"name"`
		Level  int    `json:"level"`
		Icon   string `json:"icon"`
		Ispro  int    `json:"ispro"`
		Isyear int    `json:"isyear"`
	} `json:"Business"`
}

// 获取版本信息请求
type GetVersionInfoRequest struct{}

// 获取版本信息响应
type GetVersionInfoResponse = Response[VersionInfo]

// 版本信息结构
type VersionInfo struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	ProtocolVersion string `json:"protocol_version"`
	NtProtocol      string `json:"nt_protocol"`
}
