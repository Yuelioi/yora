package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yora/internal/log"
	"yora/middleware"
	"yora/protocols/onebot"
	"yora/protocols/onebot/message"
)

func main() {

	globalLogger := log.NewDefaultLogger("全局")
	qqAdapter := onebot.NewAdapter()
	bot := onebot.NewBot(context.Background())

	// 添加中间件
	bot.AddMiddleware(middleware.LoggingMiddleware()) // 日志
	bot.AddMiddleware(middleware.RecoveryMiddleware())
	bot.AddMiddleware(middleware.RateLimitMiddleware(10, 1*time.Minute)) // 频率限制
	bot.Register(qqAdapter)

	banHandler := onebot.NewHandler(
		func(event *onebot.Event) error {
			bot.Send("group", "0", event.GroupID(), message.New(message.NewTextSegment("ECHO")).Append(message.NewAtSegment("435826135")))
			return nil
		},
	)

	cmdMatcher := onebot.OnCommand("echo", nil, 10, false, banHandler)

	bot.AddMatcher(cmdMatcher)

	server := &http.Server{
		Addr: ":12001",
	}

	// 启动 HTTP 服务器
	go func() {
		globalLogger.Debug().Msg("启动 HTTP 服务器")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			globalLogger.Error().Err(err).Msg("HTTP 服务器启动失败")
		}
	}()

	// 等待中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	// 关闭 HTTP 服务器
	globalLogger.Info().Msg("收到中断信号，正在关闭服务器...")

	if err := server.Close(); err != nil {
		globalLogger.Error().Err(err).Msg("关闭 HTTP 服务器失败")
	}
}
