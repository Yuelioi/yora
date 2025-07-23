package onebot

import (
	"encoding/json"
	"fmt"
	"slices"
	"yora/internal/adapter"
	"yora/internal/event"
)

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
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
	supportedSubTypes := []string{"normal"}

	// 校验事件类型
	if !slices.Contains(supportedEventTypes, event.Type()) {
		return fmt.Errorf("unsupported event type")
	}
	if !slices.Contains(supportedSubTypes, event.SubType()) {
		return fmt.Errorf("unsupported sub type")
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
