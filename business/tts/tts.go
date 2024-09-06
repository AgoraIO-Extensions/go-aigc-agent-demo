package tts

import (
	"context"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/business/tts/ali"
	"go-aigc-agent-demo/business/tts/ms"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

type TTS interface {
	Send(ctx context.Context, segmentID int, segmentContent string)
	GetResult() <-chan []byte
}

type Factory struct {
	Vendor      config.TTSSelect
	concurrence int
}

func NewFactory(vendor config.TTSSelect, concurrence int) (*Factory, error) {
	switch vendor {
	case config.AliTTS, config.MsTTS:
		return &Factory{Vendor: vendor, concurrence: concurrence}, nil
	default:
		return nil, fmt.Errorf("不支持vendor参数:%s", vendor)
	}
}

func (f *Factory) CreateTTS(sid int64) (TTS, error) {
	switch f.Vendor {
	case config.AliTTS:
		return ali.NewTTS(sid, f.concurrence), nil
	case config.MsTTS:
		c := config.Inst().TTS.MS
		msConfig := ms.NewTTSConfig(c.SetLog, c.SpeechKey, c.SpeechRegion, c.LanguageCheckMode, c.SpecifyLanguage, c.OutputVoice, common.Riff16Khz16BitMonoPcm)
		start := time.Now()
		msTTS, err := ms.NewTTS(sid, f.concurrence, msConfig)
		if err != nil {
			return nil, fmt.Errorf("[ms.NewTTS]%v", err)
		}
		logger.Debug("[tts]<duration> ms.NewTTS", slog.Int64("dur", time.Since(start).Milliseconds()), sentencelifecycle.Tag(sid))
		return msTTS, nil
	default:
		return nil, fmt.Errorf("不支持vendor参数:%s", f.Vendor)
	}
}
