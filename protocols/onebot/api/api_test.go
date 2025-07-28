package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"yora/protocols/onebot/client"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	LocalFile  string
	LocalImage string
	ImageURL   string
	GID        int
	UID        int
	TID        int
)

var (
	onceTest     sync.Once
	instanceTest *TestHelper
)

type Config struct {
	Tests struct {
		LocalFile  string `mapstructure:"local_file"`
		LocalImage string `mapstructure:"local_image"`
		ImageUrl   string `mapstructure:"image_url"`
		GID        int    `mapstructure:"gid"`
		UID        int    `mapstructure:"uid"`
		TID        int    `mapstructure:"tid"` // 被玩弄的对象
	} `mapstructure:"tests"`
}

type TestHelper struct {
	api *API
	t   *testing.T
	ctx context.Context
}

func NewTestHelperInstance(t *testing.T) *TestHelper {
	api := GetAPI()
	ctx := context.Background()
	initAPITestServer(ctx)
	return &TestHelper{
		api: api,
		t:   t,
		ctx: ctx,
	}
}

func NewTestHelper(t *testing.T) *TestHelper {
	onceTest.Do(func() {
		instanceTest = NewTestHelperInstance(t)
	})
	return instanceTest
}
func (h *TestHelper) StatusOk(resp any, err error, messages ...any) {
	require.NoError(h.t, err, messages...)
	require.NotNil(h.t, resp, messages...)

	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	statusField := v.FieldByName("Status")
	retcodeField := v.FieldByName("Retcode")

	require.True(h.t, statusField.IsValid(), "Status field not found")
	require.True(h.t, retcodeField.IsValid(), "Retcode field not found")

	assert.Equal(h.t, "ok", statusField.String(), messages...)
	assert.Equal(h.t, 0, int(retcodeField.Int()), messages...)
	messages = append([]any{fmt.Sprintf("测试 %s", h.t.Name())}, messages...)
	h.t.Log(messages...)
}

// 获取本地测试绝对路径
func (h *TestHelper) getTestFilePath() (string, error) {
	f := filepath.Join("..", "..", "..", LocalImage)
	return filepath.Abs(f)
}

// 获取本地测试图片的base64编码
func (h *TestHelper) getTestImageBase64() (string, error) {
	imagePath := filepath.Join("..", "..", "..", LocalImage)
	absPath, err := filepath.Abs(imagePath)
	if err != nil {
		return "", fmt.Errorf("获取图片绝对路径失败: %v", err)
	}

	imageData, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("读取图片文件失败: %v", err)
	}

	return "base64://" + base64.StdEncoding.EncodeToString(imageData), nil
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

func init() {
	c, err := loadConfig()
	if err != nil {
		panic(err)
	}
	ImageURL = c.Tests.ImageUrl
	GID = c.Tests.GID
	UID = c.Tests.UID
	LocalFile = c.Tests.LocalFile
	LocalImage = c.Tests.LocalImage
	TID = c.Tests.TID
}

// 封装初始化函数
func initAPITestServer(ctx context.Context) {
	http.HandleFunc("/onebot/v11/ws", func(w http.ResponseWriter, r *http.Request) {
		client := client.GetClient(ctx)
		client.HandleWebSocket(w, r, func(msg []byte) {
			fmt.Printf("收到消息: %s\n", string(msg))
		})
	})

	go http.ListenAndServe(":12001", nil)

	fmt.Println("WebSocket服务器启动在: ws://localhost:12001/onebot/v11/ws")
	fmt.Println("等待连接...")

}
