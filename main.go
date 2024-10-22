package main

import (
	"fmt"
	"github.com/google/uuid"
	"go-aigc-agent-demo/business/aigcCtx/sentence"
	"go-aigc-agent-demo/business/engine"
	"go-aigc-agent-demo/clients/alitts"
	qwenCli "go-aigc-agent-demo/clients/qwen"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/alibaba/speech"
	chat_gpt "go-aigc-agent-demo/pkg/azureopenai/chat-gpt"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var buildTimeStamp string

func main() {
	fmt.Println("buildTimeStamp:", buildTimeStamp)
	var err error

	// load config
	if err = config.Init("./config/chat-robot.toml"); err != nil {
		panic(fmt.Sprintf("Failed to initialize configuration based on the configuration file:%s", err))
	}

	cfg := config.Inst()

	// init log
	logger.AddContextHook(sentence.LogHook)
	logger.Init(cfg.Log.File, cfg.Log.Level, map[any]any{"uuid": strings.Join(strings.Split(uuid.New().String(), "-"), "")})
	logger.Info(fmt.Sprintf("buildTimeStamp:%s, config:%+v", buildTimeStamp, cfg))

	if err = initDependency(cfg); err != nil {
		logger.Error(err.Error(), slog.String("func", "initDependency"))
		fmt.Printf("InitDependency execution failed: %v", err)
		os.Exit(1)
	}

	// init engine
	em, err := engine.InitEngine()
	if err != nil {
		logger.Error(err.Error(), slog.String("func", "engine.InitEngine"))
		os.Exit(1)
	}
	logger.Info("EngineManager initialized successfully...")

	// start engine
	if err = em.Run(); err != nil {
		logger.Error("[em.Run]fail", slog.String("err", err.Error()))
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig

	logger.Info(fmt.Sprintf("Received exit signal %v, exiting soon", s))
	return
}

// initDependency Initialize various dependencies(global variables in the pkg/client package)
func initDependency(cfg *config.Config) error {
	if cfg.STT.Select == config.AliSTT || cfg.TTS.Select == config.AliTTS {
		if err := speech.InitToken(cfg.STT.Ali.AKID, cfg.STT.Ali.AKKey); err != nil {
			return fmt.Errorf("[speech.InitToken]%v", err)
		}
	}

	if cfg.TTS.Select == config.AliTTS {
		if err := alitts.Init(cfg.TTS.Ali.URL, cfg.TTS.Ali.AppKey, speech.TOKEN); err != nil {
			return fmt.Errorf("[alitts.Init]%v", err)
		}
	}

	// init「llm」
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
