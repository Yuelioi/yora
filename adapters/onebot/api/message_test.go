package api

import (
	"strconv"
	"testing"
	"time"

	"yora/adapters/onebot/message"
	"yora/adapters/onebot/models"
)

// 消息API测试辅助类
type MessageTestHelper struct {
	*TestHelper
}

func NewMessageTestHelper(t *testing.T) *MessageTestHelper {
	h := NewTestHelper(t)
	return &MessageTestHelper{
		TestHelper: h,
	}
}

// 发送群消息并获取消息ID
func (h *MessageTestHelper) sendGroupMessageAndGetID() (int, func()) {
	resp, err := h.api.SendMessage(0, GID, message.New("测试消息"))
	h.StatusOk(resp, err, "发送群消息")

	// 等待消息发送成功
	time.Sleep(time.Second * 2)

	callback := func() {
		resp2, err := h.api.DeleteMessage(resp.Data.MessageID)
		h.StatusOk(resp2, err, "撤回消息")
	}

	return resp.Data.MessageID, callback

}

// 发送私聊消息并获取消息ID
func (h *MessageTestHelper) sendPrivateMessageAndGetID() (int, func()) {
	resp, err := h.api.SendMessage(UID, 0, message.New("测试消息"))
	h.StatusOk(resp, err, "发送私聊消息")

	time.Sleep(time.Second * 3)

	callback := func() {
		resp2, err := h.api.DeleteMessage(resp.Data.MessageID)
		h.StatusOk(resp2, err, "撤回消息")
	}

	return resp.Data.MessageID, callback

}

// 测试删除精华消息
func TestDeleteEssenceMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	messageID := 12345

	// 执行测试
	resp, err := h.api.DeleteEssenceMessage(messageID)
	h.StatusOk(resp, err, "删除精华消息")

}

// 测试撤回消息
func TestDeleteMessage(t *testing.T) {

	// 测试撤回群消息
	h := NewMessageTestHelper(t)
	// mid, _ := h.sendGroupMessageAndGetID()
	// resp, err := h.api.DeleteMessage(mid)
	// h.StatusOk(resp, err, "撤回群聊消息")

	// 测试撤回私聊消息
	pid, _ := h.sendPrivateMessageAndGetID()
	resp2, err2 := h.api.DeleteMessage(pid)
	h.StatusOk(resp2, err2, "撤回私聊消息")

}

// !测试私聊戳一戳
func TestPrivatePoke(t *testing.T) {
	h := NewMessageTestHelper(t)

	h.t.Skip("私聊戳一戳接口未实现")

	// 执行测试
	resp, err := h.api.PrivatePoke(UID)
	h.StatusOk(resp, err, "私聊戳一戳")

}

// !测试获取精华消息列表
func TestGetEssenceMessageList(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 执行测试
	resp, err := h.api.GetEssenceMessageList(GID)

	// 验证结果
	h.StatusOk(resp, err, "获取精华消息列表应该成功")

}

// !测试获取合并转发消息
func TestGetForwardMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	messageID := "test_forward_id"

	// 执行测试
	resp, err := h.api.GetForwardMessage(messageID)

	// 验证结果
	h.StatusOk(resp, err, "获取合并转发消息应该成功")

}

// ! 测试获取好友历史聊天记录
func TestGetFriendChatHistory(t *testing.T) {
	h := NewMessageTestHelper(t)

	mid, call := h.sendPrivateMessageAndGetID()
	defer call()

	// 测试数据
	userID := UID

	count := 20

	// 执行测试
	resp, err := h.api.GetFriendChatHistory(userID, mid, count)

	// 验证结果
	h.StatusOk(resp, err, "获取好友历史聊天记录应该成功")
}

// * 测试获取群历史聊天记录
func TestGetGroupChatHistory(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	count := 20
	messageID, cleanup := h.sendGroupMessageAndGetID()
	defer cleanup()

	// 执行测试
	resp, err := h.api.GetGroupChatHistory(GID, strconv.Itoa(messageID), count)

	// 验证结果
	h.StatusOk(resp, err, "获取群历史聊天记录应该成功")
	h.t.Logf("群聊历史消息: %v", len(resp.Data.Messages))

}

// 测试获取消息
func TestGetMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	mid, callback := h.sendGroupMessageAndGetID()
	defer callback()

	// 执行测试
	resp, err := h.api.GetMessage(mid)

	// 验证结果
	h.StatusOk(resp, err, "获取消息应该成功")

}

// !测试群里戳一戳
func TestGroupPoke(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	groupID := GID
	userID := UID

	// 执行测试
	resp, err := h.api.GroupPoke(groupID, userID)

	// 验证结果
	h.StatusOk(resp, err, "群里戳一戳应该成功")
}

// 测试标记消息为已读
func TestMarkMessageAsRead(t *testing.T) {
	h := NewMessageTestHelper(t)

	mid, callback := h.sendGroupMessageAndGetID()
	defer callback()

	// 执行测试
	resp, err := h.api.MarkMessageAsRead(mid)

	h.StatusOk(resp, err, "标记消息为已读")

}

// !测试构造合并转发消息
func TestConstructForwardMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	h.t.Skip("构造合并转发消息接口有异常")

	messages := models.NewMessages()
	messages.AddNode(strconv.Itoa(UID), "张三").AddContentToLast(message.NewAtSegment(strconv.Itoa(TID)))

	// 执行测试
	resp, err := h.api.ConstructForwardMessage(messages.Messages)

	// 验证结果
	h.StatusOk(resp, err, "构造合并转发消息应该成功")

}

// !测试发送群AI语音
func TestSendGroupAIVoice(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	character := "test_character"
	groupID := GID
	text := "测试语音文本"
	chatType := 1

	// 执行测试
	resp, err := h.api.SendGroupAIVoice(character, groupID, text, chatType)

	// 验证结果
	h.StatusOk(resp, err, "发送群AI语音应该成功")

}

// 测试发送群聊合并转发消息
func TestSendGroupForwardMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	groupID := GID

	messages := models.NewMessages()
	messages.AddNode(strconv.Itoa(UID), "张三").AddContentToLast(message.NewAtSegment(strconv.Itoa(TID)))

	// 执行测试
	resp, err := h.api.SendGroupForwardMessage(groupID, messages.Messages)

	// 验证结果
	h.StatusOk(resp, err, "发送群聊合并转发消息应该成功")

}

// 测试发送消息
func TestSendMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	userID := UID
	groupID := GID
	msg := message.NewMessageBuilder().Append(message.NewAtSegment(strconv.Itoa(TID))).Append(message.NewAtSegment(strconv.Itoa(UID)))

	// 执行测试
	resp, err := h.api.SendMessage(userID, groupID, msg)

	// 验证结果
	h.StatusOk(resp, err, "发送消息应该成功")

}

// 测试发送私聊合并转发消息
func TestSendPrivateForwardMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	userID := UID
	messages := models.NewMessages()
	messages.AddNode(strconv.Itoa(UID), "张三").AddContentToLast(message.NewAtSegment(strconv.Itoa(TID)))

	// 执行测试
	resp, err := h.api.SendPrivateForwardMessage(userID, messages.Messages)

	// 验证结果
	h.StatusOk(resp, err, "发送私聊合并转发消息应该成功")

}

// 测试设置精华消息
func TestSetEssenceMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	mid, call := h.sendGroupMessageAndGetID()
	defer call()

	// 执行测试
	resp, err := h.api.SetEssenceMessage(mid)

	// 验证结果
	h.StatusOk(resp, err, "设置精华消息应该成功")
}

// 测试发送私聊消息
func TestSendPrivateMessage(t *testing.T) {
	h := NewMessageTestHelper(t)

	// 测试数据
	userID := UID

	// 执行测试
	resp, err := h.api.SendPrivateMessage(userID, message.New("测试"))

	// 验证结果
	h.StatusOk(resp, err, "发送私聊消息应该成功")

}
