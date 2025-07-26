package models

type Status = struct {
	Yes bool `json:"yes"`
}

type File struct {
	File string `json:"file"`
}

// 检查是否可以发送图片
type CanSendImageRequest struct{}

// 检查是否可以发送图片
type CanSendImageResponse = Response[Status]

// 检查是否可以发送语音
type CanSendRecordRequest struct{}

type CanSendRecordResponse = Response[Status]

// 上传图片
type UploadImageRequest = File

type UploadImageResponse = Response[string]
