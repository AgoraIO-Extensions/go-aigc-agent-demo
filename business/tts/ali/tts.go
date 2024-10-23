package ali

import (
	"context"
	"go-aigc-agent-demo/business/tts/common"
	"go-aigc-agent-demo/clients/alitts"
)

type TTS struct {
	*common.HttpSender
	ctx context.Context
}

func NewTTS(ctx context.Context) *TTS {
	return &TTS{
		HttpSender: common.NewHttpSender(ctx, alitts.Inst().StreamAsk),
		ctx:        ctx,
	}
}

func (tts *TTS) GetResult() <-chan []byte {
	return tts.HttpSender.Result()
}
