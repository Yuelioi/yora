package api

import (
	"yora/adapters/onebot/client"
	"yora/adapters/onebot/models"
)

// 检查是否可以发送图片
func (api *API) CanSendImage() (*models.CanSendImageResponse, error) {
	req := models.CanSendImageRequest{}
	return client.Call[models.CanSendImageRequest, models.CanSendImageResponse](api.client, "can_send_image", req)

}

// 检查是否可以发送语音
func (api *API) CanSendRecord() (*models.CanSendRecordResponse, error) {
	req := models.CanSendRecordRequest{}
	return client.Call[models.CanSendRecordRequest, models.CanSendRecordResponse](api.client, "can_send_record", req)

}

// 上传图片
func (api *API) UploadImage(file string) (*models.UploadImageResponse, error) {
	req := models.UploadImageRequest{
		File: file,
	}
	return client.Call[models.UploadImageRequest, models.UploadImageResponse](api.client, "upload_image", req)

}
