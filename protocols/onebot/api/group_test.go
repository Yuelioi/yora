package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 群组测试辅助类
type GroupTestHelper struct {
	*TestHelper
}

func NewGroupTestHelper(t *testing.T) *GroupTestHelper {
	h := NewTestHelper(t)
	return &GroupTestHelper{
		TestHelper: h,
	}
}

// 辅助函数：创建测试公告并返回公告ID和清理函数
func (h *GroupTestHelper) createTestNotice(groupID int) (string, func()) {
	content := "测试公告内容 - " + time.Now().Format("2006-01-02 15:04:05")

	resp, err := h.api.SendGroupNotice(groupID, content, "")
	h.StatusOk(resp, err, "创建测试公告")

	// 获取刚创建的公告ID
	resp2, err := h.api.GetGroupNotice(groupID)
	h.StatusOk(resp2, err, "获取测试公告ID")

	require.NotEmpty(h.t, resp.Data, "公告列表不能为空")

	noticeID := resp2.Data[0].NoticeID
	require.NotEmpty(h.t, noticeID, "公告ID不能为空")

	cleanup := func() {
		resp, err := h.api.DeleteGroupNotice(groupID, noticeID)
		h.StatusOk(resp, err, "清理测试公告")
	}

	return noticeID, cleanup
}

// 辅助函数：获取最近的群消息ID
func (h *GroupTestHelper) getRecentGroupMessageID(groupID int) int {
	// 这里应该调用获取群消息历史的API
	// 为了测试，返回一个模拟的消息ID
	return 123456
}

// 1. 删除群公告测试
func TestDeleteGroupNotice(t *testing.T) {
	h := NewGroupTestHelper(t)
	groupID := GID

	// 创建一个测试公告
	noticeID, cleanup := h.createTestNotice(groupID)
	defer cleanup()

	// 删除公告
	resp, err := h.api.DeleteGroupNotice(groupID, noticeID)
	h.StatusOk(resp, err, "删除群公告")

	t.Logf("删除群公告成功，群ID: %d, 公告ID: %s", groupID, noticeID)
}

// 2. 获取群公告测试
func TestGetGroupNotice(t *testing.T) {
	h := NewGroupTestHelper(t)
	groupID := GID

	resp, err := h.api.GetGroupNotice(groupID)
	h.StatusOk(resp, err, "获取群公告")
	assert.IsType(t, []interface{}{}, resp.Data, "Data字段应为slice类型")

	t.Logf("获取群公告成功，群ID: %d, 公告数量: %d", groupID, len(resp.Data))
}

// !3. 获取AI语音测试
func TestGetAIRecord(t *testing.T) {
	h := NewGroupTestHelper(t)
	h.t.Skip("跳过获取AI语音测试")

	character := "1" // 声色ID
	groupID := GID
	text := "你好，这是一个测试语音"
	chatType := 1

	resp, err := h.api.GetAIRecord(character, groupID, text, chatType)
	h.StatusOk(resp, err, "获取AI语音")

	t.Logf("获取AI语音成功，声色ID: %s, 群ID: %d, 文本: %s", character, groupID, text)
}

// 4. 获取群荣耀信息测试
func TestGetGroupHonorInfo(t *testing.T) {
	h := NewGroupTestHelper(t)

	groupID := GID
	honorType := "all"

	resp, err := h.api.GetGroupHonorInfo(groupID, honorType)
	h.StatusOk(resp, err, "获取群荣耀信息")

	t.Logf("获取群荣耀信息成功，群ID: %d, 荣耀类型: %s", groupID, honorType)
}

// 5. 设置群管理员测试
func TestSetGroupAdmin(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过设置群管理员测试 需要群主权限")
	groupID := GID
	userID := TID

	resp, err := h.api.SetGroupAdmin(groupID, userID, true)
	h.StatusOk(resp, err, "设置群管理员")
}

// 6. 设置群成员禁言测试
func TestSetGroupBan(t *testing.T) {

	h := NewGroupTestHelper(t)

	resp, err := h.api.SetGroupBan(TID, GID, 60)
	h.StatusOk(resp, err, "设置群成员禁言")

	t.Logf("设置群成员禁言成功，用户ID: %d, 群ID: %d, 禁言时长: 60秒", TID, GID)

	// 解除禁言
	time.Sleep(1 * time.Second)
	resp2, err2 := h.api.SetGroupBan(TID, GID, 0)
	h.StatusOk(resp2, err2, "解除群成员禁言")
}

// !7. 设置群Bot发言状态测试
func TestSetGroupBotStatus(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过设置Bot发言状态测试")

	groupID := GID
	botID := 123456 // 模拟Bot ID
	enable := 1     // 启用

	resp, err := h.api.SetGroupBotStatus(groupID, botID, enable)
	h.StatusOk(resp, err, "设置群Bot发言状态")

	t.Logf("设置群Bot发言状态成功，群ID: %d, BotID: %d, 状态: %d", groupID, botID, enable)
}

// !8. 调用群机器人回调测试
func TestSendGroupBotCallback(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过调用群机器人回调测试")

	groupID := GID
	botID := 123456
	data1 := "test_data_1"
	data2 := "test_data_2"

	resp, err := h.api.SendGroupBotCallback(groupID, botID, data1, data2)
	h.StatusOk(resp, err, "调用群机器人回调")

	t.Logf("调用群机器人回调成功，群ID: %d, BotID: %d", groupID, botID)
}

// 9. 设置群名片测试
func TestSetGroupCard(t *testing.T) {
	h := NewGroupTestHelper(t)

	userID := TID
	groupID := GID
	originalCard := "原始名片"
	testCard := "测试名片-" + time.Now().Format("15:04:05")

	// 设置测试名片
	resp, err := h.api.SetGroupCard(userID, groupID, testCard)
	h.StatusOk(resp, err, "设置群名片")

	t.Logf("设置群名片成功，用户ID: %d, 群ID: %d, 名片: %s", userID, groupID, testCard)

	// 恢复原始名片
	time.Sleep(1 * time.Second)
	resp2, err2 := h.api.SetGroupCard(userID, groupID, originalCard)
	h.StatusOk(resp2, err2, "恢复原始名片")
}

// 10. 踢出群成员测试
func TestKickGroupMember(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过踢出群成员测试")

	resp, err := h.api.KickGroupMember(TID, GID, false)
	h.StatusOk(resp, err, "踢出群成员")

	t.Logf("踢出群成员成功，用户ID: %d, 群ID: %d", TID, GID)
}

// 11. 退出群测试
func TestLeaveGroup(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过退出群测试")

	resp, err := h.api.LeaveGroup(GID, false)
	h.StatusOk(resp, err, "退出群")

	t.Logf("退出群成功，群ID: %d", GID)
}

// 12. 发送群公告测试
func TestSendGroupNotice(t *testing.T) {
	h := NewGroupTestHelper(t)

	groupID := GID
	content := "测试公告内容 - " + time.Now().Format("2006-01-02 15:04:05")
	image := ""

	resp, err := h.api.SendGroupNotice(groupID, content, image)
	h.StatusOk(resp, err, "发送群公告")

	t.Logf("发送群公告成功，群ID: %d, 内容: %s", groupID, content)

	// 清理：删除刚发送的公告
	time.Sleep(1 * time.Second)
	resp2, err := h.api.GetGroupNotice(groupID)
	h.StatusOk(resp2, err, "获取公告用于清理")

	assert.NotEmpty(t, resp.Data, "公告列表不为空")

	resp3, cleanupErr := h.api.DeleteGroupNotice(groupID, resp2.Data[0].NoticeID)
	h.StatusOk(resp3, cleanupErr, "清理测试公告")

}

// 13. 设置群名称测试
func TestSetGroupName(t *testing.T) {
	h := NewGroupTestHelper(t)

	resp, err := h.api.SetGroupName(GID, "qq测试群")
	h.StatusOk(resp, err, "设置群名称")

	t.Logf("设置群名称成功，群ID: %d, 名称: qq测试群", GID)
}

// 14. 设置全体禁言测试
func TestSetGroupWholeBan(t *testing.T) {
	h := NewGroupTestHelper(t)
	groupID := GID
	enable := true

	resp, err := h.api.SetGroupWholeBan(groupID, enable)
	h.StatusOk(resp, err, "设置全体禁言")

}

// 15. 设置群头像测试
func TestSetGroupPortrait(t *testing.T) {
	h := NewGroupTestHelper(t)
	groupID := GID

	// 使用base64编码的图片
	imageBase64, err := h.getTestImageBase64()
	require.NoError(t, err, "获取测试图片")
	require.NotEmpty(t, imageBase64, "测试图片不能为空")

	file := "data:image/jpeg;base64," + imageBase64

	resp, err := h.api.SetGroupPortrait(groupID, file)
	h.StatusOk(resp, err, "设置群头像")

	t.Logf("设置群头像成功，群ID: %d", groupID)
}

// !16. 设置群表情回复测试
func TestSetEmojiReaction(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过设置群表情回复测试")

	groupID := GID
	messageID := h.getRecentGroupMessageID(groupID)
	code := "128077" // 👍表情代码
	isAdd := true

	require.Greater(t, messageID, 0, "消息ID必须大于0")

	resp, err := h.api.SetEmojiReaction(groupID, messageID, code, isAdd)
	h.StatusOk(resp, err, "设置群表情回复")

	t.Logf("设置群表情回复成功，群ID: %d, 消息ID: %d, 表情代码: %s", groupID, messageID, code)

	// 移除表情回复
	time.Sleep(1 * time.Second)
	resp2, err2 := h.api.SetEmojiReaction(groupID, messageID, code, false)
	h.StatusOk(resp2, err2, "移除表情回复")
}

// 17. 设置群专属头衔测试
func TestSetGroupSpecialTitle(t *testing.T) {
	h := NewGroupTestHelper(t)

	h.t.Skip("跳过设置群专属头衔测试")

	groupID := GID
	userID := UID // 给自己设置头衔
	specialTitle := "测试头衔"
	duration := 3600 // 1小时

	resp, err := h.api.SetGroupSpecialTitle(groupID, userID, specialTitle, duration)
	h.StatusOk(resp, err, "设置群专属头衔")

	t.Logf("设置群专属头衔成功，群ID: %d, 用户ID: %d, 头衔: %s, 有效期: %d秒",
		groupID, userID, specialTitle, duration)

	// 清理：移除头衔
	time.Sleep(2 * time.Second)
	resp2, err2 := h.api.SetGroupSpecialTitle(groupID, userID, "", 0)
	h.StatusOk(resp2, err2, "移除专属头衔")
}
