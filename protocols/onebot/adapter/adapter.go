package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"yora/internal/adapter"
	"yora/internal/event"
	"yora/internal/message"
	"yora/protocols/onebot/client"

	onebotEvent "yora/protocols/onebot/event"
)

var _ adapter.Adapter = (*Adapter)(nil)

type Adapter struct {
	Client *client.Client
}

// CallAPI implements adapter.Adapter.
func (a *Adapter) CallAPI(action string, params any) (any, error) {
	return a.Client.CallAPI(action, params)
}

func NewAdapter() *Adapter {

	ctx := context.Background()
	return &Adapter{
		Client: client.GetClient(ctx),
	}
}

// GetCapabilities implements adapter.Adapter.
func (a *Adapter) GetCapabilities() adapter.Capabilities {
	return adapter.Capabilities{
		SupportsGroupChat:   true,
		SupportsPrivateChat: true,
		SupportsFileUpload:  true,
		SupportsRichText:    false,
		SupportsReply:       true,
		SupportsForward:     true,
		SupportsEdit:        false,
		SupportsDelete:      true,
		SupportedSegmentTypes: []string{
			"dice",
			"forward",
			"json",
			"location",
			"longmsg",
			"mface",
			"music",
			"poke",
			"record",
			"rps",
			"video",
			"at",
			"face",
			"image",
			"reply",
			"text",
			"file",
		},
		MaxMessageLength: 0,
		MaxFileSize:      0,
		Extra:            map[string]any{},
	}
}

// ParseEvent implements adapter.Adapter.
func (a *Adapter) ParseEvent(raw any) (event.Event, error) {

	data, ok := raw.([]byte)
	if !ok {
		return nil, fmt.Errorf("ParseEvent: raw 类型应为 []byte，实际为 %T", raw)
	}

	var base struct {
		Type string `json:"post_type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("解析事件类型失败: %w", err)
	}

	switch base.Type {
	case "message":
		var e onebotEvent.MessageEvent
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("解析 MessageEvent 失败: %w", err)
		}
		return &e, nil
	case "notice":
		var e onebotEvent.NoticeEvent
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("解析 NoticeEvent 失败: %w", err)
		}
		return &e, nil
	case "meta_event":
		var e onebotEvent.MetaEvent
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("解析 MetaEvent 失败: %w", err)
		}
		return &e, nil
	case "request":
		var e onebotEvent.RequestEvent
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("解析 RequestEvent 失败: %w", err)
		}
		return &e, nil
	default:
		return nil, fmt.Errorf("未知事件类型: %s", base.Type)
	}

}

func (a *Adapter) Send(messageType string, userId string, groupId string, msg message.Message) (any, error) {

	gid, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, fmt.Errorf("groupId or userId is not int")
	}
	uid, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("groupId or userId is not int")
	}

	return a.Client.SendMessage(messageType, uid, gid, msg)
}

// ParseMessage implements adapter.Adapter.
func (a *Adapter) ParseMessage(raw string) ([]message.Segment, error) {
	var segments []message.Segment
	return segments, nil
}

// Protocol implements adapter.Adapter.
func (a *Adapter) Protocol() adapter.Protocol {
	return adapter.ProtocolOneBot
}

// ValidateEvent implements adapter.Adapter.
func (a *Adapter) ValidateEvent(event event.Event) error {

	supportedEventTypes := []string{"message", "notice", "request", "meta_event"}

	// 校验事件类型
	if !slices.Contains(supportedEventTypes, event.Type()) {
		return fmt.Errorf("unsupported event type")
	}

	return nil

}
