package tts

import (
	"context"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
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
		return nil, fmt.Errorf("[tts] incorrect value for the vendor parameter:%s", vendor)
	}
}

func (f *Factory) CreateTTS(ctx context.Context) (TTS, error) {
	switch f.Vendor {
	case config.AliTTS:
		return ali.NewTTS(ctx, f.concurrence), nil
	case config.MsTTS:
		c := config.Inst().TTS.MS
		msConfig := ms.NewTTSConfig(c.SetLog, c.SpeechKey, c.SpeechRegion, c.LanguageCheckMode, c.SpecifyLanguage, c.OutputVoice, common.Riff16Khz16BitMonoPcm)
		start := time.Now()
		msTTS, err := ms.NewTTS(ctx, f.concurrence, msConfig)
		if err != nil {
			return nil, fmt.Errorf("[ms.NewTTS]%v", err)
		}
		logger.DebugContext(ctx, "[tts]<duration> ms.NewTTS", slog.Int64("dur", time.Since(start).Milliseconds()))
		return msTTS, nil
	default:
		return nil, fmt.Errorf("incorrect value for the vendor parameter:%s", f.Vendor)
	}
}
