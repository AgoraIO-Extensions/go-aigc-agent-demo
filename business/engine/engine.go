package engine

import (
	"fmt"
	"go-aigc-agent-demo/business/exit"
	"go-aigc-agent-demo/business/filter"
	"go-aigc-agent-demo/business/llm"
	"go-aigc-agent-demo/business/rtc"
	"go-aigc-agent-demo/business/sentence"
	"go-aigc-agent-demo/business/stt"
	"go-aigc-agent-demo/business/tts"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
)

type Engine struct {
	exitWrapper *exit.ExitManager
	filter      *filter.Filter
	rtc         *rtc.RTC
	sttFactory  *stt.Factory
	ttsFactory  *tts.Factory
	llm         *llm.LLM
}

func InitEngine() (*Engine, error) {
	cfg := config.Inst()
	e := &Engine{}

	var err error

	// init「vad」
	e.filter = filter.NewFilter(sentence.FirstSid, cfg.Filter.Vad.StartWin, cfg.Filter.Vad.StopWin)

	// init「rtc」
	e.rtc = rtc.NewRTC(cfg.RTC.AppID, "", cfg.RTC.ChannelName, cfg.RTC.UserID, cfg.RTC.Region)
	logger.Info("RTC initialization succeeded")

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

	// init [exit]
	e.exitWrapper = exit.NewExitManager(cfg.StartTime, cfg.MaxLifeTime)

	// init「llm」
	e.llm, err = llm.NewLLM(cfg.LLM.ModelSelect, cfg.LLM.Prompt.Generate(), &cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("[llm.NewLLM]%v", err)
	}

	return e, nil
}

func (e *Engine) Run() error {
	// asynchronously: Exit automatically after reaching the maximum lifetime
	e.exitWrapper.HandlerMaxLifeTime()

	// Register user leave event handler
	e.rtc.SetOnUserLeft(e.exitWrapper.OnUserLeft)

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
