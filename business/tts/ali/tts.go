package ali

import (
	"context"
	"go-aigc-agent-demo/business/tts/common"
	"go-aigc-agent-demo/clients/alitts"
)

type TTS struct {
	*common.Sender
	ctx context.Context
}

func NewTTS(ctx context.Context, con int) *TTS {
	return &TTS{
		Sender: common.NewHttpSender(ctx, alitts.Inst().StreamAsk, con),
		ctx:    ctx,
	}
}
