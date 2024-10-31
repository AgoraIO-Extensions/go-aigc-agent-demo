package tts

import (
	"context"
	"fmt"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"go-aigc-agent-demo/business/aigcCtx"
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
	concurrency int
}

func NewFactory(vendor config.TTSSelect, concurrency int) (*Factory, error) {
	fac := &Factory{Vendor: vendor, concurrency: concurrency}
	switch vendor {
	case config.AliTTS:
		return fac, nil
	case config.MsTTS:
		c := config.Inst().TTS.MS
		if err := ms.Init(c.SetLog, c.SpeechKey, c.SpeechRegion, c.LanguageCheckMode, c.SpecifyLanguage, c.OutputVoice, common.Riff16Khz16BitMonoPcm); err != nil {
			return nil, fmt.Errorf("[ms.Init]%w", err)
		}
		ms.PreConn(concurrency)
		return fac, nil
	default:
		return nil, fmt.Errorf("[tts] incorrect value for the vendor parameter:%s", vendor)
	}
}

func (f *Factory) CreateTTS(ctx *aigcCtx.AIGCContext) (TTS, error) {
	switch f.Vendor {
	case config.AliTTS:
		return ali.NewTTS(ctx, f.concurrency), nil
	case config.MsTTS:
		start := time.Now()
		msTTS := ms.NewTTS(ctx, f.concurrency)
		logger.DebugContext(ctx, "[tts]<duration> ms.NewTTS", slog.Int64("dur", time.Since(start).Milliseconds()))
		return msTTS, nil
	default:
		return nil, fmt.Errorf("incorrect value for the vendor parameter:%s", f.Vendor)
	}
}
