package main

import (
	"goblin/internal/options"
	"goblin/internal/reverse"
	"os"
	"os/signal"
	"syscall"
	log "unknwon.dev/clog/v2"
)

const (
	LogConsoleNum = 100
)

func main() {
	_ = log.NewConsole(LogConsoleNum,
		log.ConsoleConfig{
			Level: log.LevelWarn,
		})
	// 命令行和全局配置初始化
	options := *options.ParseOptions()

	_ = log.NewFile(
		log.FileConfig{
			Level:    options.SetLogLevel(),
			Filename: options.LogFile,
		},
	)
	defer log.Stop()
	// 数据库配置初始化
	options.Cache.Init()
	// 检查缓存系统是否正常
	if err := options.Cache.ValidateCachePing(); err != nil {
		log.Fatal("%s", err.Error())
	}
	// IP 数据库初始化
	options.IPLocation.Init()

	// 初始化所有服务
	s := reverse.InitServerConfig(&options)
	// 开始服务
	s.Start()

	// 等待中断信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// 停止所有服务
	s.Stop()
	log.Error("All listeners shut down. Exiting.")
}
