package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name        string `mapstructure:"name"`
		Version     string `mapstructure:"version"`
		Description string `mapstructure:"description"`
	} `mapstructure:"app"`

	Plugins struct {
		Chat struct {
			APIKey string `mapstructure:"api_key"`
		} `mapstructure:"chat"`
	} `mapstructure:"plugins"`
}

func findProjectRoot(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		cfgPath := filepath.Join(dir, filename)
		if _, err := os.Stat(cfgPath); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // 已经到根目录
		}
		dir = parent
	}
	return "", fmt.Errorf("找不到配置文件 %s", filename)
}

func loadConfig() (*Config, error) {
	root, err := findProjectRoot("yora.yaml")
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("yora")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(root) // 指定根目录路径

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func TestChatAI(t *testing.T) {
	c, err := loadConfig()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
		return
	}
	APIKey = c.Plugins.Chat.APIKey
	fmt.Printf("APIKey: %s\n", APIKey)

	// 示例用法
	userInfo := GroupMemberInfo{
		UserID:   123456789,
		Nickname: "测试用户",
		Card:     "群名片",
	}

	rawMessages := []RawMessage{
		{
			UserID: 123456789,
			Time:   time.Now().Unix(),
			Message: []MessageItem{
				{Type: "text", Data: map[string]interface{}{"text": "你好"}},
			},
			Sender: map[string]interface{}{"nickname": "测试用户"},
		},
	}

	response, err := ChatAI("你好呀", userInfo, rawMessages)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("回复: %s\n", response)
}
