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

func NewTTS(ctx context.Context, con int) *TTS {
	return &TTS{
		HttpSender: common.NewHttpSender(ctx, con, alitts.Inst().StreamAsk),
		ctx:        ctx,
	}
}

func (tts *TTS) Send(ctx context.Context, segmentID int, segmentContent string) {
	tts.HttpSender.Send(ctx, segmentID, segmentContent)
}

func (tts *TTS) GetResult() <-chan []byte {
	return tts.HttpSender.Result()
}
