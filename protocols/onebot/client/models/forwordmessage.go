package models

import (
	"encoding/json"
	"yora/protocols/onebot/message"
)

// 节点数据结构
type NodeData struct {
	UserID   string            `json:"user_id"`
	Nickname string            `json:"nickname"`
	Content  []message.Segment `json:"content"`
}

// 消息节点结构
type MessageNode struct {
	Type string   `json:"type"`
	Data NodeData `json:"data"`
}

// 根消息结构
type ForwardMessages struct {
	Messages []MessageNode `json:"messages"`
}

// ============= 构造器方法 =============

// 创建新的 Messages 实例
func NewMessages() *ForwardMessages {
	return &ForwardMessages{
		Messages: make([]MessageNode, 0),
	}
}

// 创建节点数据
func NewNodeData(userID, nickname string) NodeData {
	return NodeData{
		UserID:   userID,
		Nickname: nickname,
		Content:  make([]message.Segment, 0),
	}
}

// ============= 链式构造器 =============

// MessageBuilder 用于链式构造
type ForwardMessageBuilder struct {
	messages *ForwardMessages
}

// 创建消息构造器
func NewForwardMessageBuilder() *ForwardMessageBuilder {
	return &ForwardMessageBuilder{
		messages: NewMessages(),
	}
}

// 添加消息节点
func (mb *ForwardMessageBuilder) AddNode(userID, nickname string) *NodeBuilder {
	nodeData := NewNodeData(userID, nickname)
	node := MessageNode{
		Type: "node",
		Data: nodeData,
	}

	mb.messages.Messages = append(mb.messages.Messages, node)

	// 返回节点构造器，指向刚添加的节点
	return &NodeBuilder{
		node:    &mb.messages.Messages[len(mb.messages.Messages)-1],
		builder: mb,
	}
}

// 构建最终消息
func (mb *ForwardMessageBuilder) Build() *ForwardMessages {
	return mb.messages
}

// NodeBuilder 用于构造节点内容
type NodeBuilder struct {
	node    *MessageNode
	builder *ForwardMessageBuilder
}

// 结束当前节点构造，返回消息构造器
func (nb *NodeBuilder) Done() *ForwardMessageBuilder {
	return nb.builder
}

// ============= 便捷方法 =============

// Messages 的便捷方法
func (m *ForwardMessages) AddNode(userID, nickname string) *ForwardMessages {
	nodeData := NewNodeData(userID, nickname)
	node := MessageNode{
		Type: "node",
		Data: nodeData,
	}
	m.Messages = append(m.Messages, node)
	return m
}

// 为最后一个节点添加内容
func (m *ForwardMessages) AddContentToLast(content message.Segment) *ForwardMessages {
	if len(m.Messages) > 0 {
		lastIdx := len(m.Messages) - 1
		m.Messages[lastIdx].Data.Content = append(m.Messages[lastIdx].Data.Content, content)
	}
	return m
}

// 转换为 JSON
func (m *ForwardMessages) ToJSON() (string, error) {
	data, err := json.MarshalIndent(m, "", "  ")
	return string(data), err
}
