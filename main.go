package main

import (
	"fmt"
	"go-aigc-agent-demo/business/engine"
	"go-aigc-agent-demo/business/workerid"
	"go-aigc-agent-demo/clients/alitts"
	qwenCli "go-aigc-agent-demo/clients/qwen"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/alibaba/speech"
	chat_gpt "go-aigc-agent-demo/pkg/azureopenai/chat-gpt"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var buildTimeStamp string

func main() {
	fmt.Println("buildTimeStamp:", buildTimeStamp)
	var err error

	// 加载配置文件
	if err = config.Init("./config/chat-robot.toml"); err != nil {
		panic(fmt.Sprintf("基于配置文件初始化配置失败:%s", err))
	}

	cfg := config.Inst()

	// 初始化日志
	logger.Init(cfg.Log.File, cfg.Log.Level, map[any]any{"uuid": workerid.UUID})
	logger.Info(fmt.Sprintf("buildTimeStamp:%s, config:%+v", buildTimeStamp, cfg))

	if err = initDependency(cfg); err != nil {
		logger.Error(err.Error(), slog.String("func", "initDependency"))
		fmt.Printf("InitDependency执行失败: %v", err)
		os.Exit(1)
	}

	// 初始化 engine
	em, err := engine.InitEngine()
	if err != nil {
		logger.Error(err.Error(), slog.String("func", "engine.InitEngine"))
		os.Exit(1)
	}
	logger.Info("EngineManager初始化成功...")

	// 启动 engine
	if err = em.Run(); err != nil {
		logger.Error("fail to ")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig

	logger.Info(fmt.Sprintf("收到退出信号%v，即将退出", s))
	return
}

// initDependency 初始化各种依赖（pkg/client包中的全局变量）
func initDependency(cfg *config.Config) error {
	if err := speech.InitToken(cfg.STT.Ali.AKID, cfg.STT.Ali.AKKey); err != nil {
		return fmt.Errorf("[speech.InitToken]%v", err)
	}

	if err := alitts.Init(cfg.TTS.Ali.URL, cfg.TTS.Ali.AppKey, speech.TOKEN); err != nil {
		return fmt.Errorf("[alitts.Init]%v", err)
	}

	// 初始化「llm」
	switch cfg.LLM.ModelSelect {
	case config.LLMQwen:
		if err := qwenCli.Init(cfg.LLM.QWen.URL, cfg.LLM.QWen.ApiKey); err != nil {
			return fmt.Errorf("[qwenCli.Init]%v", err)
		}
	case config.LLMChatGPT4o:
		if err := chat_gpt.InitChatGPT(chat_gpt.NewConfig(cfg.LLM.ChatGPT4o.Key, cfg.LLM.ChatGPT4o.EndPoint)); err != nil {
			return fmt.Errorf("[chat_gpt.InitChatGPT 4o]%v", err)
		}
	}

	return nil
}
