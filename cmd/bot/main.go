package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"yora/middleware"
	"yora/plugins/builtin/echo"
	"yora/plugins/builtin/help"
	"yora/plugins/community/repeater"
	"yora/protocols/onebot/adapter"
	"yora/protocols/onebot/bot"
)

func main() {

	qqAdapter := adapter.NewAdapter()
	bot := bot.NewBot()

	// 添加中间件
	bot.AddMiddleware(middleware.LoggingMiddleware()) // 日志
	bot.AddMiddleware(middleware.RecoveryMiddleware())
	bot.AddMiddleware(middleware.RateLimitMiddleware(10, 1*time.Minute)) // 频率限制

	// 注册适配器
	bot.RegisterAdapter(qqAdapter)

	// 注册插件

	bot.LoadPlugins(echo.Echo, help.Helper, repeater.Repeater)

	// 启动机器人
	if err := bot.Run(); err != nil {
		panic(err)
	}

	// 等待中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

}
