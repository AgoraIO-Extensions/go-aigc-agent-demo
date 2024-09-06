package ali

import (
	"context"
	"go-aigc-agent-demo/business/tts/common"
	"go-aigc-agent-demo/clients/alitts"
)

type TTS struct {
	*common.HttpSender
	sid int64
}

func NewTTS(sid int64, con int) *TTS {
	return &TTS{
		HttpSender: common.NewHttpSender(sid, con, alitts.Inst().StreamAsk),
		sid:        sid,
	}
}

func (tts *TTS) Send(ctx context.Context, segmentID int, segmentContent string) {
	tts.HttpSender.Send(ctx, segmentID, segmentContent)
}

func (tts *TTS) GetResult() <-chan []byte {
	return tts.HttpSender.Result()
}
