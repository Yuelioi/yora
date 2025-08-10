package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"yora/adapters/onebot/adapter"
	"yora/middleware"
	"yora/pkg/bot"
	"yora/pkg/conf"
	"yora/plugins/builtin/echo"
	"yora/plugins/builtin/help"

	"github.com/rs/zerolog"
)

func main() {
	// 设置日志
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	qqAdapter := adapter.NewAdapter()
	bot := bot.NewBot(conf.NewBotConfig())

	// 添加中间件
	bot.RegisterMiddlewares(
		middleware.LoggingMiddleware(),
		middleware.RateLimitMiddleware(10, 1*time.Minute), // 频率限制
		middleware.RecoveryMiddleware(),
	)

	// 注册适配器
	bot.RegisterAdapters(qqAdapter)

	// 注册插件
	bot.RegisterPlugins(echo.New(), help.New())
	// bot.RegisterPlugins(funny.Plugins...)

	// 启动机器人
	if err := bot.Run(); err != nil {
		panic(err)
	}

	// 等待中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

}
