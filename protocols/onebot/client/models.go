package client

// 发送消息的结构体
type APIRequest struct {
	Action string `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo,omitempty"`
}

type APIResponse struct {
	Status  string         `json:"status"`
	RetCode int            `json:"retcode"`
	Data    map[string]any `json:"data"`
	Echo    string         `json:"echo,omitempty"`
}
