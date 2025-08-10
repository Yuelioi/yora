package messages

import (
	"strings"
)

// 快速创建 Segment 对象

// 快速创建 Message 对象

// todo 快速构造消息

type MessageHelper struct {
	Message
}

func NewHelper(msg Message) *MessageHelper {
	return &MessageHelper{Message: msg}
}

// Text 添加文本段
func (mh *MessageHelper) Text(text string) *MessageHelper {
	mh.Message = append(mh.Message, NewTextSegment(text))
	return mh
}

// At 添加@某人段
func (mh *MessageHelper) At(qq string) *MessageHelper {
	mh.Message = append(mh.Message, NewAtSegment(qq))
	return mh
}

// Image 添加图片段
func (mh *MessageHelper) Image(file string) *MessageHelper {
	mh.Message = append(mh.Message, NewImageSegment(file, "image.png", "", "[图片]", "0"))
	return mh
}

// Face 添加表情段
func (mh *MessageHelper) Face(id string) *MessageHelper {
	mh.Message = append(mh.Message, NewFaceSegment(id, false))
	return mh
}

// Reply 添加回复段
func (mh *MessageHelper) Reply(id string) *MessageHelper {
	mh.Message = append(mh.Message, NewReplySegment(id))
	return mh
}

// Record 添加语音段
func (mh *MessageHelper) Record(url, file string) *MessageHelper {
	mh.Message = append(mh.Message, NewRecordSegment(url, file))
	return mh
}

// Video 添加视频段
func (mh *MessageHelper) Video(url, file string) *MessageHelper {
	mh.Message = append(mh.Message, NewVideoSegment(url, file))
	return mh
}

// GetText 获取所有文本内容
func (mh *MessageHelper) GetText() string {
	var texts []string
	for _, seg := range mh.Message {
		if seg.IsType("text") {
			if text, ok := seg.GetData("text"); ok {
				if str, ok := text.(string); ok {
					texts = append(texts, str)
				}
			}
		}
	}
	return strings.Join(texts, "")
}

// GetAtList 获取所有@的QQ号
func (mh *MessageHelper) GetAtList() []string {
	var atList []string
	for _, seg := range mh.Message {
		if seg.IsType("at") {
			if qq, ok := seg.GetData("qq"); ok {
				if str, ok := qq.(string); ok {
					atList = append(atList, str)
				}
			}
		}
	}
	return atList
}

// GetImages 获取所有图片URL
func (mh *MessageHelper) GetImages() []string {
	var images []string
	for _, seg := range mh.Message {
		if seg.IsType("image") {
			if url, ok := seg.GetData("url"); ok {
				if str, ok := url.(string); ok {
					images = append(images, str)
				}
			}
		}
	}
	return images
}

// GetReplyID 获取回复的消息ID
func (mh *MessageHelper) GetReplyID() string {
	for _, seg := range mh.Message {
		if seg.IsType("reply") {
			if id, ok := seg.GetData("id"); ok {
				if str, ok := id.(string); ok {
					return str
				}
			}
		}
	}
	return ""
}

// HasText 是否包含文本
func (mh *MessageHelper) HasText() bool {
	return mh.Message.HasType("text")
}

// HasAt 是否包含@某人
func (mh *MessageHelper) HasAt() bool {
	return mh.Message.HasType("at")
}

// HasImage 是否包含图片
func (mh *MessageHelper) HasImage() bool {
	return mh.Message.HasType("image")
}

// HasReply 是否包含回复
func (mh *MessageHelper) HasReply() bool {
	return mh.Message.HasType("reply")
}

// IsAtMe 是否@了指定的QQ号
func (mh *MessageHelper) IsAtMe(myQQ string) bool {
	atList := mh.GetAtList()
	for _, qq := range atList {
		if qq == myQQ {
			return true
		}
	}
	return false
}

// TextSegment 文本段包装器
type TextSegment struct {
	*Segment
}

func (ts *TextSegment) GetText() string {
	if text, ok := ts.GetData("text"); ok {
		if str, ok := text.(string); ok {
			return str
		}
	}
	return ""
}

// AtSegment @段包装器
type AtSegment struct {
	*Segment
}

func (as *AtSegment) GetQQ() string {
	if qq, ok := as.GetData("qq"); ok {
		if str, ok := qq.(string); ok {
			return str
		}
	}
	return ""
}

// ImageSegment 图片段包装器
type ImageSegment struct {
	*Segment
}

func (is *ImageSegment) GetURL() string {
	if url, ok := is.GetData("url"); ok {
		if str, ok := url.(string); ok {
			return str
		}
	}
	return ""
}

func (is *ImageSegment) GetFile() string {
	if file, ok := is.GetData("file"); ok {
		if str, ok := file.(string); ok {
			return str
		}
	}
	return ""
}

// GetTextSegments 获取所有文本段
func (mh *MessageHelper) GetTextSegments() []*TextSegment {
	var result []*TextSegment
	for _, seg := range mh.Message {
		if seg.IsType("text") {
			if s, ok := seg.(*Segment); ok {
				result = append(result, &TextSegment{s})
			}
		}
	}
	return result
}

// GetAtSegments 获取所有@段
func (mh *MessageHelper) GetAtSegments() []*AtSegment {
	var result []*AtSegment
	for _, seg := range mh.Message {
		if seg.IsType("at") {
			if s, ok := seg.(*Segment); ok {
				result = append(result, &AtSegment{s})
			}
		}
	}
	return result
}

// GetImageSegments 获取所有图片段
func (mh *MessageHelper) GetImageSegments() []*ImageSegment {
	var result []*ImageSegment
	for _, seg := range mh.Message {
		if seg.IsType("image") {
			if s, ok := seg.(*Segment); ok {
				result = append(result, &ImageSegment{s})
			}
		}
	}
	return result
}

// NewMessageBuilder 创建消息构建器
func NewMessageBuilder() *MessageHelper {
	return &MessageHelper{Message: Message{}}
}
