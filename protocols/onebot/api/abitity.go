package api

import (
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/client/models"
)

// GetPrivateFile 获取私聊文件
func CanSendImage(userID string, fileID string, fileHash string) (*models.GetPrivateFileResponse, error) {
	req := models.GetPrivateFileRequest{
		UserID:   userID,
		FileID:   fileID,
		FileHash: fileHash,
	}
	resp, err := client.Call[models.GetPrivateFileRequest, models.GetPrivateFileResponse](c, "get_private_file", req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
