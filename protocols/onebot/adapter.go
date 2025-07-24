package onebot

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"yora/internal/adapter"
	"yora/internal/event"
	"yora/protocols/onebot/client"
	"yora/protocols/onebot/message"
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
		Client: client.NewClient(ctx),
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

	var e Event

	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func (a *Adapter) Send(messageType string, userId string, groupId string, message message.Message) (any, error) {

	gid, err := strconv.Atoi(groupId)
	if err != nil {
		return nil, fmt.Errorf("groupId or userId is not int")
	}
	uid, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("groupId or userId is not int")
	}

	return a.Client.SendMessage(messageType, uid, gid, message)
}

// ParseMessage implements adapter.Adapter.
func (a *Adapter) ParseMessage(raw string) ([]event.Segment, error) {
	var segments []event.Segment
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

type PlatformInfo struct {
	Platform        string
	PlatformVersion string
	AppVersion      string
}

// Extra implements event.PlatformInfo.
func (p *PlatformInfo) Extra() map[string]any {
	m := make(map[string]any)
	m["app_version"] = p.AppVersion
	return m
}

// Name implements event.PlatformInfo.
func (p *PlatformInfo) Name() string {
	return p.Platform
}

// Version implements event.PlatformInfo.
func (p *PlatformInfo) Version() string {
	return p.PlatformVersion
}
