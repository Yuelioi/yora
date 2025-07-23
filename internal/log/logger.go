package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Component string

const (
	BotComponent Component = "bot"
	APIComponent Component = "api"

	MiddlewareComponent Component = "middleware"
	PluginComponent     Component = "plugin"
	DependComponent     Component = "depend"
	HandlerComponent    Component = "handler"
	EventComponent      Component = "event"

	WebsocketComponent Component = "websocket"

	DatabaseComponent Component = "database"

	DefaultComponent Component = "default"
)

// NewLogger 创建带主题的日志记录器
func NewLogger(component Component, name string) zerolog.Logger {
	// 根据组件类型选择主题
	theme := getComponentTheme(component)

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",

		FormatLevel: func(i interface{}) string {
			level := strings.ToUpper(fmt.Sprintf("%s", i))
			color := getLevelColor(level)
			return fmt.Sprintf("%s[%-5s]%s", color, level, "\x1b[0m")
		},

		FormatMessage: func(i interface{}) string {
			if i == nil {
				return ""
			}
			return fmt.Sprintf("[%s %s] %s", theme.Icon, name, i)
		},

		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("\x1b[90m%s\x1b[0m:", i) // 灰色字段名
		},

		FormatFieldValue: func(i interface{}) string {
			fieldStr := fmt.Sprintf("%s", i)
			if strings.Contains(fieldStr, string(component)) {
				return fmt.Sprintf("%s%s%s", theme.Color, fieldStr, "\x1b[0m")
			}
			return fmt.Sprintf("\x1b[97m%s\x1b[0m", fieldStr) // 亮白色
		},

		FormatTimestamp: func(i interface{}) string {
			t, ok := i.(string)
			if !ok {
				return ""
			}
			parsed, err := time.Parse(time.RFC3339, t)
			if err != nil {
				return t
			}
			return fmt.Sprintf("\x1b[90m%s\x1b[0m", parsed.Format("01-02 15:04:05"))
		},
	}

	return zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()
}

// ComponentTheme 组件主题结构
type ComponentTheme struct {
	Icon  string
	Color string
}

// getComponentTheme 根据组件获取主题
func getComponentTheme(component Component) ComponentTheme {
	themes := map[Component]ComponentTheme{
		BotComponent:        {"🤖", "\x1b[94m"},  // 蓝色
		DependComponent:     {"⚙️", "\x1b[90m"}, // 黑色
		WebsocketComponent:  {"🔌", "\x1b[95m"},  // 紫色
		APIComponent:        {"🌐", "\x1b[96m"},  // 青色
		PluginComponent:     {"🔧", "\x1b[93m"},  // 黄色
		HandlerComponent:    {"⚡", "\x1b[92m"},  // 绿色
		DatabaseComponent:   {"🗄️", "\x1b[91m"}, // 红色
		MiddlewareComponent: {"🔗", "\x1b[97m"},  // 白色
		EventComponent:      {"📨", "\x1b[35m"},  // 紫红色
		DefaultComponent:    {"📋", "\x1b[37m"},  // 白色

	}

	if theme, exists := themes[component]; exists {
		return theme
	}
	return ComponentTheme{"📋", "\x1b[37m"} // 默认主题
}

// getLevelColor 获取日志级别颜色
func getLevelColor(level string) string {
	colors := map[string]string{
		"DEBUG": "\x1b[36m", // 青色
		"INFO":  "\x1b[32m", // 绿色
		"WARN":  "\x1b[33m", // 黄色
		"ERROR": "\x1b[31m", // 红色
		"FATAL": "\x1b[35m", // 紫色
		"PANIC": "\x1b[41m", // 红色背景
	}

	if color, exists := colors[level]; exists {
		return color
	}
	return "\x1b[37m" // 默认白色
}

// NewLoggerWithLevel 创建带日志级别的记录器
func NewLoggerWithLevel(component Component, name string, level zerolog.Level) zerolog.Logger {
	logger := NewLogger(component, name)
	return logger.Level(level)
}

// NewGlobalLogger 设置全局日志记录器
func NewGlobalLogger(component Component, name string) {
	logger := NewLogger(component, name)
	zerolog.DefaultContextLogger = &logger
}

func NewBotLogger(name string) zerolog.Logger {
	return NewLogger(BotComponent, name)
}
func NewDependLogger(name string) zerolog.Logger {
	return NewLogger(DependComponent, name)
}

func NewWebsocketLogger(name string) zerolog.Logger {
	return NewLogger(WebsocketComponent, name)
}

func NewAPILogger(name string) zerolog.Logger {
	return NewLogger(APIComponent, name)
}

func NewPluginLogger(name string) zerolog.Logger {
	return NewLogger(PluginComponent, name)
}

func NewHandlerLogger(name string) zerolog.Logger {
	return NewLogger(HandlerComponent, name)
}

func NewDatabaseLogger(name string) zerolog.Logger {
	return NewLogger(DatabaseComponent, name)
}

func NewMiddlewareLogger(name string) zerolog.Logger {
	return NewLogger(MiddlewareComponent, name)
}

func NewEventLogger(name string) zerolog.Logger {
	return NewLogger(EventComponent, name)
}
func NewDefaultLogger(name string) zerolog.Logger {
	return NewLogger(DefaultComponent, name)
}
