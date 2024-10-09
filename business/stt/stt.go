package stt

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/stt/ali"
	"go-aigc-agent-demo/business/stt/common"
	"go-aigc-agent-demo/business/stt/ms"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/alibaba/speech"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"time"
)

type STT interface {
	Send(chunk []byte, end bool) error
	GetResult() <-chan *common.Result
}

type Factory struct {
	vendorName config.SttSelect
	aliConfig  *ali.Config
	msConfig   *ms.Config
}

func NewFactory(vendorName config.SttSelect, sttConfig config.STT) (*Factory, error) {
	factory := &Factory{vendorName: vendorName}
	var err error
	switch vendorName {
	case config.AliSTT:
		c := sttConfig.Ali
		factory.aliConfig, err = ali.Init(c.URL, c.AppKey, speech.TOKEN)
		if err != nil {
			return nil, err
		}
	case config.MsSTT:
		c := sttConfig.MS
		factory.msConfig = ms.NewConfig(c.SetLog, c.LanguageCheckMode, c.AutoAudioCheckLanguage, c.SpecifyLanguage, c.SpeechKey, c.SpeechRegion)
	default:
		return nil, fmt.Errorf("incorrect vendorname parameter")
	}
	return factory, nil
}

func (f *Factory) CreateSTT(ctx context.Context) (STT, error) {
	switch f.vendorName {
	case config.AliSTT:
		aliSTT, err := ali.NewSTT(ctx, f.aliConfig)
		if err != nil {
			return nil, fmt.Errorf("[ali.NewSTT]%w", err)
		}
		return aliSTT, nil
	case config.MsSTT:
		start := time.Now()
		msSTT, err := ms.NewSTT(ctx, f.msConfig)
		if err != nil {
			return nil, fmt.Errorf("[ms.NewSTT]%w", err)
		}
		logger.InfoContext(ctx, "[stt]<duration> ms.NewSTT", slog.Int64("dur", time.Since(start).Milliseconds()))
		return msSTT, nil
	default:
		return nil, fmt.Errorf("incorrect vendorname parameter")
	}
}
