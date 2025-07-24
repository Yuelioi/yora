package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"yora/middleware"
	"yora/plugins/builtin/echo"
	"yora/protocols/onebot"
)

func main() {

	qqAdapter := onebot.NewAdapter()
	bot := onebot.NewBot()

	// 添加中间件
	bot.AddMiddleware(middleware.LoggingMiddleware()) // 日志
	bot.AddMiddleware(middleware.RecoveryMiddleware())
	bot.AddMiddleware(middleware.RateLimitMiddleware(10, 1*time.Minute)) // 频率限制
	bot.Register(qqAdapter)

	// 注册插件
	e := echo.Echo

	bot.AddPlugin(e)
	if err := bot.Run(); err != nil {
		panic(err)
	}

	// 等待中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

}
