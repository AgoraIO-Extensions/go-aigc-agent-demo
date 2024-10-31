package engine

import (
	"fmt"
	"go-aigc-agent-demo/business/aigcCtx/sentence"
	"go-aigc-agent-demo/business/filter"
	"go-aigc-agent-demo/business/llm"
	"go-aigc-agent-demo/business/rtc"
	"go-aigc-agent-demo/business/rtm"
	"go-aigc-agent-demo/business/stt"
	"go-aigc-agent-demo/business/tts"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/agora-go-sdk/go_wrapper/agoraservice"
	"go-aigc-agent-demo/pkg/logger"
	"go-aigc-agent-demo/pkg/monitor"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Engine struct {
	StartTime   int64
	MaxLifeTime int64 // Maximum process uptime.
	filter      *filter.Filter
	rtc         *rtc.RTC
	sttFactory  *stt.Factory
	ttsFactory  *tts.Factory
	llm         *llm.LLM
}

func InitEngine() (*Engine, error) {
	cfg := config.Inst()
	e := &Engine{
		StartTime:   cfg.StartTime,
		MaxLifeTime: cfg.MaxLifeTime,
	}

	var err error

	// init「vad」
	e.filter = filter.NewFilter(sentence.FirstSid, cfg.Filter.Vad.StartWin, cfg.Filter.Vad.StopWin)

	// init「rtc」
	e.rtc = rtc.NewRTC(cfg.RTC.AppID, "", cfg.RTC.ChannelName, cfg.RTC.UserID, cfg.RTC.Region)
	logger.Info("RTC initialization succeeded")

	userID, err := strconv.Atoi(cfg.RTC.UserID)
	if err != nil {
		return nil, fmt.Errorf("[strconv.Atoi]%v", err)
	}

	rtm.Init(1, int32(userID), sentence.FirstSid, e.rtc)

	// init「stt」
	if e.sttFactory, err = stt.NewFactory(cfg.STT.Select, cfg.STT); err != nil {
		return nil, fmt.Errorf("[stt.NewFactory]%v", err)
	}
	logger.Info("STT initialization succeeded")

	// init tts
	if e.ttsFactory, err = tts.NewFactory(cfg.TTS.Select, 2); err != nil {
		return nil, fmt.Errorf("初始化tts失败.%v", err)
	}
	logger.Info("TTS initialization succeeded")

	// init「llm」
	e.llm, err = llm.NewLLM(cfg.LLM.ModelSelect, cfg.LLM.Prompt.Generate(), &cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("[llm.NewLLM]%v", err)
	}

	return e, nil
}

func (e *Engine) Run() error {
	// asynchronously: Exit automatically after reaching the maximum lifetime
	e.HandlerMaxLifeTime()

	// Register user leave event handler
	e.rtc.SetOnUserLeft(e.OnUserLeft)

	// Register handler for audio received from RTC
	e.rtc.SetOnReceiveAudio(e.filter.OnRcvRTCAudio)

	// stt
	sttInput := e.filter.OutputAudio()
	sttOutput := make(chan *sentenceGroupText, 20)
	go e.ProcessSTT(sttInput, sttOutput)

	// llm
	llmOutput := make(chan *llmResult, 20)
	go e.ProcessLLM(sttOutput, llmOutput)

	// tts
	ttsOutput := make(chan *ttsResult, 20)
	go e.ProcessTTS(llmOutput, ttsOutput)

	// send to rtc
	go e.ProcessSendRTC(ttsOutput)

	// connect to rtc
	if err := e.rtc.Connect(); err != nil {
		return fmt.Errorf("[rtc.Connect]%v", err)
	}
	return nil
}

// OnUserLeft Handle user departure events (currently supports only single user scenarios)
func (e *Engine) OnUserLeft(conn *agoraservice.RtcConnection, uid string, reason int) {
	logger.Info("[exit] User has left; the process is about to exit.", slog.String("uid", uid))
	monitor.LogGoroutines()
	os.Exit(0)
}

func (e *Engine) HandlerMaxLifeTime() {
	go func() {
		leftLifeTime := e.MaxLifeTime - (time.Now().Unix() - e.StartTime)
		if leftLifeTime <= 0 {
			logger.Info("Reached maximum uptime; exiting soon...")
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("Remaining uptime: %d", leftLifeTime))
		<-time.After(time.Second * time.Duration(leftLifeTime))
		logger.Info("Reached maximum uptime; exiting soon...")
		monitor.LogGoroutines()
		os.Exit(0)
	}()
}
