package models

type ImageData struct {
	File     string `json:"file"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	SubType  string `json:"subType"`
	Summary  string `json:"summary"`
}

type ReplyData struct {
	MessageID string `json:"message_id"`
}

type TextData struct {
	Text string `json:"text"`
}

type AtData struct {
	UserID string `json:"user_id"`
}
