package api

import (
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/client/models"
)

// 获取自定义表情
func (api *API) FetchCustomFace() (*models.FetchCustomFaceResponse, error) {
	req := models.FetchCustomFaceRequest{}
	return client.Call[models.FetchCustomFaceRequest, models.FetchCustomFaceResponse](api.client, "fetch_custom_face", req)

}

// 获取商城表情 key
//
// 参数：
//   - emojiIDs: 表情 ID 列表
func (api *API) FetchMfaceKey(emojiIDs []string) (*models.FetchMfaceKeyResponse, error) {
	req := models.FetchMfaceKeyRequest{
		Emoji_IDs: emojiIDs,
	}
	return client.Call[models.FetchMfaceKeyRequest, models.FetchMfaceKeyResponse](api.client, "fetch_mface_key", req)

}

// 加入好友表情接龙
//
// 参数：
//   - userID: 用户 ID
//   - messageID: 消息 ID
//   - emojiID: 表情 ID
func (api *API) JoinFriendEmojiChain(userID int, messageID int, emojiID int) (*models.JoinFriendEmojiChainResponse, error) {
	req := models.JoinFriendEmojiChainRequest{
		UserID:    userID,
		MessageID: messageID,
		EmojiID:   emojiID,
	}
	return client.Call[models.JoinFriendEmojiChainRequest, models.JoinFriendEmojiChainResponse](api.client, ".join_friend_emoji_chain", req)

}

// 获取群 Ai 声色
//
// 参数：
//   - groupID: 群 ID
//   - chatType: 聊天类型（如 1 表示群聊）
func (api *API) GetAICharacters(groupID int, chatType int) (*models.GetAICharactersResponse, error) {
	req := models.GetAICharactersRequest{
		GroupID:  groupID,
		ChatType: chatType,
	}
	return client.Call[models.GetAICharactersRequest, models.GetAICharactersResponse](api.client, "get_ai_characters", req)

}

// 获取 Cookies
//
// 参数：
//   - domain: 域名，如 ".qq.com"
func (api *API) GetCookies(domain string) (*models.GetCookiesResponse, error) {
	req := models.GetCookiesRequest{
		Domain: domain,
	}
	return client.Call[models.GetCookiesRequest, models.GetCookiesResponse](api.client, "get_cookies", req)

}

// 获取 QQ 接口凭证
//
// 参数：
//   - domain: 域名，如 ".qq.com"
func (api *API) GetCredentials(domain string) (*models.GetCredentialsResponse, error) {
	req := models.GetCredentialsRequest{
		Domain: domain,
	}
	return client.Call[models.GetCredentialsRequest, models.GetCredentialsResponse](api.client, "get_credentials", req)

}

// 获取 CSRF Token
func (api *API) GetCSRFToken() (*models.GetCSRFTokenResponse, error) {
	req := models.GetCSRFTokenRequest{}
	return client.Call[models.GetCSRFTokenRequest, models.GetCSRFTokenResponse](api.client, "get_csrf_token", req)

}

// 加入群聊表情接龙
//
// 参数：
//   - groupID: 群 ID
//   - messageID: 消息 ID
//   - emojiID: 表情 ID
func (api *API) JoinGroupEmojiChain(groupID int, messageID int, emojiID int) (*models.JoinGroupEmojiChainResponse, error) {
	req := models.JoinGroupEmojiChainRequest{
		GroupID:   groupID,
		MessageID: messageID,
		EmojiID:   emojiID,
	}
	return client.Call[models.JoinGroupEmojiChainRequest, models.JoinGroupEmojiChainResponse](api.client, ".join_group_emoji_chain", req)

}

// OCR 图像识别
//
// 参数：
//   - image: http/https/file/base64
func (api *API) OCRImage(image string) (*models.OCRImageResponse, error) {
	req := models.OCRImageRequest{
		Image: image,
	}
	return client.Call[models.OCRImageRequest, models.OCRImageResponse](api.client, "ocr_image", req)

}

// 设置 QQ 头像
//
// 参数：
//   - file: http/https/file/base64
func (api *API) SetQQAvatar(file string) (*models.SetQQAvatarResponse, error) {
	req := models.SetQQAvatarRequest{
		File: file,
	}
	return client.Call[models.SetQQAvatarRequest, models.SetQQAvatarResponse](api.client, "set_qq_avatar", req)

}

// 点赞用户资料
//
// 参数：
//   - userID: 用户 ID
//   - times: 点赞次数（通常为 1~10）
func (api *API) SendLike(userID int, times int) (*models.SendLikeResponse, error) {
	req := models.SendLikeRequest{
		UserID: userID,
		Times:  times,
	}
	return client.Call[models.SendLikeRequest, models.SendLikeResponse](api.client, "send_like", req)

}

// 删除好友
//
// 参数：
//   - userID: 用户 ID
//   - block: 是否拉黑（true 表示拉黑该好友）
func (api *API) DeleteFriend(userID string, block bool) (*models.DeleteFriendResponse, error) {
	req := models.DeleteFriendRequest{
		UserID: userID,
		Block:  block,
	}
	return client.Call[models.DeleteFriendRequest, models.DeleteFriendResponse](api.client, "delete_friend", req)

}

// 获取 rkey
func (api *API) GetRKey() (*models.GetRKeyResponse, error) {
	return client.Call[interface{}, models.GetRKeyResponse](api.client, "get_rkey", struct{}{})

}
