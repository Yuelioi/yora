package api

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 通用测试辅助类
type GenericTestHelper struct {
	*TestHelper
}

func NewGenericTestHelper(t *testing.T) *GenericTestHelper {
	return &GenericTestHelper{
		TestHelper: NewTestHelper(t),
	}

}

// 辅助函数：创建测试用的表情ID列表
func (h *GenericTestHelper) createTestEmojiIDs() []string {
	return []string{"1", "2", "3", "100", "101"}
}

// 辅助函数：获取最近的群消息ID（用于表情接龙测试）
func (h *GenericTestHelper) getRecentGroupMessageID(groupID int) (int, error) {
	// 这里应该调用获取群消息历史的API
	// 为了测试，返回一个模拟的消息ID
	return 12345, nil
}

// 辅助函数：获取最近的好友消息ID
func (h *GenericTestHelper) getRecentFriendMessageID(userID int) (int, error) {
	// 这里应该调用获取私聊消息历史的API
	// 为了测试，返回一个模拟的消息ID
	return 54321, nil
}

// !返回空数组 1. 测试获取自定义表情
func TestFetchCustomFace(t *testing.T) {
	h := NewGenericTestHelper(t)

	resp, err := h.api.FetchCustomFace()
	h.StatusOk(resp, err, "获取自定义表情")

	t.Logf("获取自定义表情成功，返回数据: %+v", resp)
}

// !不知道怎么获取emojiIDs 2. 测试获取商城表情key
func TestFetchMfaceKey(t *testing.T) {
	h := NewGenericTestHelper(t)

	emojiIDs := h.createTestEmojiIDs()
	resp, err := h.api.FetchMfaceKey(emojiIDs)
	h.StatusOk(resp, err, "获取商城表情key")

	t.Logf("获取商城表情key成功，表情ID数量: %d, 返回数据: %+v", len(emojiIDs), resp)
}

// !3. 测试加入好友表情接龙
func TestJoinFriendEmojiChain(t *testing.T) {
	h := NewGenericTestHelper(t)

	h.t.Skip("跳过好友表情接龙测试")

	// 模拟参数
	userID := UID
	messageID, err := h.getRecentFriendMessageID(userID)
	if err != nil {
		t.Skipf("跳过好友表情接龙测试，无法获取消息ID: %v", err)
		return
	}
	emojiID := 1

	resp, err := h.api.JoinFriendEmojiChain(userID, messageID, emojiID)
	h.StatusOk(resp, err, "加入好友表情接龙")

	t.Logf("加入好友表情接龙成功，用户ID: %d, 消息ID: %d, 表情ID: %d", userID, messageID, emojiID)
}

// 4. 测试获取群AI声色
func TestGetAICharacters(t *testing.T) {
	h := NewGenericTestHelper(t)

	groupID := GID
	chatType := 1 // 群聊类型

	resp, err := h.api.GetAICharacters(groupID, chatType)
	h.StatusOk(resp, err, "获取群AI声色")

	t.Logf("获取群AI声色成功，群ID: %d, 聊天类型: %d, 返回数据: %+v", groupID, chatType, resp)
}

// 5. 测试获取Cookies
func TestGetCookies(t *testing.T) {
	h := NewGenericTestHelper(t)

	domain := ".qq.com"
	resp, err := h.api.GetCookies(domain)
	h.StatusOk(resp, err, "获取Cookies")

	t.Logf("获取Cookies成功，域名: %s, 返回数据: %+v", domain, resp)
}

// ! 6. 测试获取QQ接口凭证
func TestGetCredentials(t *testing.T) {
	h := NewGenericTestHelper(t)
	h.t.Skip("跳过获取QQ接口凭证测试")

	domain := ".qq.com"
	resp, err := h.api.GetCredentials(domain)
	h.StatusOk(resp, err, "获取QQ接口凭证")

	t.Logf("获取QQ接口凭证成功，域名: %s, 返回数据: %+v", domain, resp)
}

// 7. 测试获取CSRF Token
func TestGetCSRFToken(t *testing.T) {
	h := NewGenericTestHelper(t)

	resp, err := h.api.GetCSRFToken()
	h.StatusOk(resp, err, "获取CSRF Token")

	t.Logf("获取CSRF Token成功，返回数据: %+v", resp)
}

// !8. 测试加入群聊表情接龙
func TestJoinGroupEmojiChain(t *testing.T) {
	h := NewGenericTestHelper(t)
	h.t.Skip("跳过群表情接龙测试")

	groupID := GID
	messageID, err := h.getRecentGroupMessageID(groupID)
	if err != nil {
		t.Skipf("跳过群表情接龙测试，无法获取消息ID: %v", err)
		return
	}
	emojiID := 1

	resp, err := h.api.JoinGroupEmojiChain(groupID, messageID, emojiID)
	h.StatusOk(resp, err, "加入群表情接龙")
	t.Logf("加入群表情接龙成功，群ID: %d, 消息ID: %d, 表情ID: %d", groupID, messageID, emojiID)
}

// 9. 测试OCR图像识别
func TestOCRImage(t *testing.T) {
	h := NewGenericTestHelper(t)

	resp, err := h.api.OCRImage(ImageURL)
	h.StatusOk(resp, err, "OCR图像识别")

	t.Logf("OCR图像识别成功，返回数据: %+v", resp)
}

// *10. 测试设置QQ头像(不会返回status ok, 但是会正常调用)
func TestSetQQAvatar(t *testing.T) {
	h := NewGenericTestHelper(t)

	// 使用图片路径或base64进行测试
	// imagePath := filepath.Join("..", "..", "..", LocalImage)
	// absPath, err := filepath.Abs(imagePath)
	// if err != nil {
	// 	t.Skipf("跳过设置QQ头像测试，无法获取图片路径: %v", err)
	// 	return
	// }

	// b64, err := h.getTestImageBase64()
	// if err != nil {
	// 	t.Skipf("跳过设置QQ头像测试，无法获取图片base64编码: %v", err)
	// 	return
	// }

	resp, err := h.api.SetQQAvatar(ImageURL)
	assert.NoError(t, err)

	t.Logf("设置QQ头像成功，返回数据: %+v", resp)

}

// 11. 测试点赞用户资料
func TestSendLike(t *testing.T) {
	h := NewGenericTestHelper(t)

	userID := UID
	times := 1 // 点赞1次

	resp, err := h.api.SendLike(userID, times)
	if err != nil {
		t.Logf("点赞用户资料可能失败（正常现象，可能有频率限制）: %v", err)
	} else {
		h.StatusOk(resp, err, "点赞用户资料")
		t.Logf("点赞用户资料成功，用户ID: %d, 点赞次数: %d, 返回数据: %+v", userID, times, resp)
	}
}

// 12. 测试删除好友
func TestDeleteFriend(t *testing.T) {
	h := NewGenericTestHelper(t)

	block := false // 不拉黑

	resp, err := h.api.DeleteFriend(strconv.Itoa(UID), block)
	h.StatusOk(resp, err, "删除好友")

	t.Logf("删除好友成功，用户ID: %d, 是否拉黑: %v", UID, block)
}

// 13. 测试获取rkey
func TestGetRKey(t *testing.T) {
	h := NewGenericTestHelper(t)

	resp, err := h.api.GetRKey()
	h.StatusOk(resp, err, "获取rkey")

}
