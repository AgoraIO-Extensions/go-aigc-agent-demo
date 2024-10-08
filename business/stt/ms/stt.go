package ms

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/pkg/logger"
)

type LanguageCheckMode int

const (
	AutoCheck LanguageCheckMode = 0 // Automatic language detection mode
	Specify   LanguageCheckMode = 1 // Specify language mode
)

// Config stt initialization configuration
type Config struct {
	speechKey    string
	speechRegion string
	//inputFormat
	setLog                 bool // Whether to enable Microsoft SDK logging
	languageCheckMode      LanguageCheckMode
	autoAudioCheckLanguage []string // Automatic audio language detection range. For example：{"zh-CN", "en-US", "ja-JP"}
	specifyLanguage        string   // Language of the specified input audio. For example,Chinese: "zh-CN"
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
	ctx    context.Context
	cfg    *Config
	client *client
}

func NewSTT(ctx context.Context, cfg *Config) (*STT, error) {
	stt := &STT{
		cfg: cfg,
	}
	c, err := newClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("[newClient]%v", err)
	}
	stt.client = c
	return stt, nil
}

func (stt *STT) Send(chunk []byte, end bool) error {
	if end {
		stt.client.sttInputStream.CloseStream() // will trigger sessionStoppedHandler
		logger.InfoContext(stt.ctx, "[stt] Send stop pushing stream signal")
		return nil
	}
	return stt.client.pumpChunkIntoStream(chunk)
}

func (stt *STT) GetResult() <-chan *common.Result {
	return stt.client.result
}
