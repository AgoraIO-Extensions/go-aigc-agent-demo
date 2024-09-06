package ms

import (
	"fmt"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type LanguageCheckMode int

const (
	AutoCheck LanguageCheckMode = 0 // 自动检测模式
	Specify   LanguageCheckMode = 1 // 指定语言模式
)

// Config stt初始化配置
type Config struct {
	speechKey    string
	speechRegion string
	//inputFormat
	setLog                 bool // 是否启用微软sdk的日志
	languageCheckMode      LanguageCheckMode
	autoAudioCheckLanguage []string // 自动音频检测语种范围。例如：汉语、英文、日语为：{"zh-CN", "en-US", "ja-JP"}
	specifyLanguage        string   // 指定的输入音频的语言。例如：汉语为："zh-CN"
}

func NewConfig(setLog bool, languageCheckMode int, autoAudioCheckLanguage []string, specifyLanguage string, speechKey string, speechRegion string) *Config {
	return &Config{
		speechKey:              speechKey,
		speechRegion:           speechRegion,
		setLog:                 setLog,
		languageCheckMode:      LanguageCheckMode(languageCheckMode),
		autoAudioCheckLanguage: autoAudioCheckLanguage,
		specifyLanguage:        specifyLanguage,
	}
}

type STT struct {
	SID    int64
	cfg    *Config
	client *client
}

func NewSTT(sid int64, cfg *Config) (*STT, error) {
	stt := &STT{
		SID: sid,
		cfg: cfg,
	}
	c, err := newClient(sid, cfg)
	if err != nil {
		return nil, fmt.Errorf("[newClient]%v", err)
	}
	stt.client = c
	return stt, nil
}

func (stt *STT) Send(chunk []byte, end bool) error {
	if end {
		stt.client.sttInputStream.CloseStream() // 会触发 sessionStoppedHandler
		logger.Info("[stt] 停止推流", slog.Int64("sid", stt.SID))
		return nil
	}
	return stt.client.pumpChunkIntoStream(chunk)
}

func (stt *STT) GetResult() <-chan *common.Result {
	return stt.client.result
}
