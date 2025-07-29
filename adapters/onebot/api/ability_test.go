package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanSendImage(t *testing.T) {
	h := NewTestHelper(t)

	resp, err := h.api.CanSendImage()

	h.StatusOk(resp, err, "可以发送图片")

	assert.Equal(t, resp.Data.Yes, true, "可以发送图片")
}

func TestCanSendRecord(t *testing.T) {
	h := NewTestHelper(t)

	resp, err := h.api.CanSendRecord()

	h.StatusOk(resp, err, "可以发送语音")
	assert.Equal(t, resp.Data.Yes, true, "可以发送语音")
}

func TestUploadImage(t *testing.T) {
	h := NewTestHelper(t)

	resp, err := h.api.UploadImage(ImageURL)
	h.StatusOk(resp, err, "上传图片")
	assert.True(t, strings.HasPrefix(resp.Data, "http"), "上传图片")
}
