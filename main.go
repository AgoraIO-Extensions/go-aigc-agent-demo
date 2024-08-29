package main

import (
	"fmt"
	"go-aigc-agent-demo/business/engine"
	"go-aigc-agent-demo/business/workerid"
	"go-aigc-agent-demo/clients/alitts"
	qwenCli "go-aigc-agent-demo/clients/qwen"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/alibaba/speech"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var buildTimeStamp string

func main() {
	fmt.Println("buildTimeStamp:", buildTimeStamp)
	var err error

	// 加载配置文件
	if err = config.Init("./config/alibaba.toml"); err != nil {
		panic(fmt.Sprintf("基于配置文件初始化配置失败:%s", err))
	}

	cfg := config.Inst()

	// 初始化日志
	if err = logger.Init(cfg.Log.File, cfg.Log.Level, workerid.UUID); err != nil {
		fmt.Printf("init logger failed, err:%s\n", err)
		os.Exit(1)
	}
	logger.Inst().Info(fmt.Sprintf("buildTimeStamp:%s, config:%+v", buildTimeStamp, cfg))

	if err = initDependency(cfg); err != nil {
		logger.Inst().Error(err.Error(), zap.String("func", "initDependency"))
		fmt.Printf("InitDependency执行失败: %v", err)
		os.Exit(1)
	}

	// 初始化 engine
	em, err := engine.InitEngine()
	if err != nil {
		logger.Inst().Error(err.Error(), zap.String("func", "engine.InitEngine"))
		os.Exit(1)
	}
	logger.Inst().Info("EngineManager初始化成功...")

	// 启动 engine
	if err = em.Run(); err != nil {
		logger.Inst().Fatal(err.Error(), zap.String("func", "em.Run"))
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig

	logger.Inst().Info(fmt.Sprintf("收到退出信号%v，即将退出", s))
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
	}

	return nil
}
