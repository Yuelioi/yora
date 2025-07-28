package message

import (
	"encoding/json"
	message "yora/internal/message"
)

var _ message.Segment = (*Segment)(nil)

type Segment struct {
	TypeStr string         `json:"type"`
	DataMap map[string]any `json:"data"`
}

// NewSegment 创建新的消息段
func NewSegment(segType string, data map[string]any) *Segment {
	if data == nil {
		data = make(map[string]any)
	}
	return &Segment{
		TypeStr: segType,
		DataMap: data,
	}
}

func (s Segment) Json() map[string]any {
	return map[string]any{
		"type": s.TypeStr,
		"data": s.DataMap,
	}
}

func (s Segment) Data() map[string]any {
	return s.DataMap
}

func (s Segment) GetData(key string) (any, bool) {
	val, ok := s.DataMap[key]
	return val, ok
}

func (s Segment) IsType(segmentType string) bool {
	return s.TypeStr == segmentType
}

func (s Segment) String() string {
	switch s.TypeStr {
	case "text":
		if text, ok := s.GetData("text"); ok {
			if str, ok := text.(string); ok {
				return str
			}
		}
	case "at":
		if qq, ok := s.GetData("qq"); ok {
			return "@" + qq.(string)
		}
	case "image":
		if url, ok := s.GetData("url"); ok {
			return url.(string)
		}
	}
	return ""
}

func (s Segment) Type() string {
	return s.TypeStr
}

// -----------------------

// TODO NewDiceSegment 创建骰子片段
func NewDiceSegment(dataMap map[string]any) Segment {
	return Segment{
		TypeStr: "dice",
		DataMap: dataMap,
	}
}

// 创建转发消息段
func NewForwardSegment(id string) Segment {
	return Segment{
		TypeStr: "forward",
		DataMap: map[string]any{
			"id": id,
		},
	}
}

// 创建JSON消息段
func NewJsonSegment(data any) Segment {
	dataMap, err := json.Marshal(data)
	if err != nil {
		return Segment{}
	}
	return Segment{
		TypeStr: "json",
		DataMap: map[string]any{
			"data": string(dataMap),
		},
	}
}

// 创建定位消息段
func NewLocationSegment(lat, lon, title, content string) Segment {
	return Segment{
		TypeStr: "location",
		DataMap: map[string]any{
			"lat":     lat,
			"lon":     lon,
			"title":   title,
			"content": content,
		},
	}
}

// 创建长消息段
func NewLongMessageSegment(id string) Segment {
	return Segment{
		TypeStr: "longmsg",
		DataMap: map[string]any{
			"id": id,
		},
	}
}

// 创建商城表情消息段
func NewShopEmojisSegment(url string, emoji_package_id int, emoji_id string, key string, summary string) Segment {
	return Segment{
		TypeStr: "mface",
		DataMap: map[string]any{
			"url":              url,
			"emoji_package_id": emoji_package_id,
			"emoji_id":         emoji_id,
			"key":              key,
			"summary":          summary,
		},
	}
}

// 创建音乐消息段
func NewMusicSegment(typ string, url string, audio string, title string, content string, image string) Segment {
	return Segment{
		TypeStr: "music",
		DataMap: map[string]any{
			"type":    typ,
			"url":     url,
			"audio":   audio,
			"title":   title,
			"content": content,
			"image":   image,
		},
	}
}

// 创建戳一戳消息段
func NewPokeSegment(typ string, strength string, id string) Segment {
	return Segment{
		TypeStr: "poke",
		DataMap: map[string]any{
			"type":     typ,
			"strength": strength,
			"id":       id,
		},
	}
}

// 创建语音消息段 file 链接, 支持 http/https/file/base64
func NewRecordSegment(url string, file string) Segment {
	return Segment{
		TypeStr: "record",
		DataMap: map[string]any{
			"url":  url,
			"file": file,
		},
	}
}

// TODO 创建猜拳消息段
func NewRpsSegment(data map[string]any) Segment {
	return Segment{
		TypeStr: "rps",
		DataMap: data,
	}
}

// 创建视频消息段
func NewVideoSegment(url string, file string) Segment {
	return Segment{
		TypeStr: "video",
		DataMap: map[string]any{
			"url":  url,
			"file": file,
		},
	}
}

// NewAtSegment 创建 @ 提及片段
func NewAtSegment(qq string) Segment {
	return Segment{
		TypeStr: "at",
		DataMap: map[string]any{
			"qq": qq,
		},
	}
}

// 创建表情消息段
func NewFaceSegment(id string, large bool) Segment {
	return Segment{
		TypeStr: "face",
		DataMap: map[string]any{
			"id":    id,
			"large": large,
		},
	}
}

// NewImageSegment 创建图片片段
func NewImageSegment(file, filename, url, summary string, subType string) Segment {
	return Segment{
		TypeStr: "image",
		DataMap: map[string]any{
			"file":     file,
			"filename": filename,
			"url":      url,
			"summary":  summary,
			"subType":  subType,
		},
	}
}

// NewReplySegment 创建回复消息段
func NewReplySegment(id string) Segment {
	return Segment{
		TypeStr: "reply",
		DataMap: map[string]any{
			"id": id,
		},
	}
}

// NewTextSegment 创建文本片段
func NewTextSegment(text string) Segment {
	return Segment{
		TypeStr: "text",
		DataMap: map[string]any{
			"text": text,
		},
	}
}

func NewFileSegment(filename, file_id, file_hash, url string) Segment {
	return Segment{
		TypeStr: "file",
		DataMap: map[string]any{
			"file_name": filename,
			"file_id":   file_id,
			"file_hash": file_hash,
			"url":       url,
		},
	}
}
