package models

// Response 通用响应结构体
type Response[T any] struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    T      `json:"data,omitempty"`
	Echo    string `json:"echo,omitempty"`
}

// 发送消息的结构体
type APIRequest struct {
	Action string `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo,omitempty"`
}
