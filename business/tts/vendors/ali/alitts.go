package ali

import (
	"context"
	"go-aigc-agent-demo/clients/alitts"
	"io"
)

type AliTTS struct{}

func NewAliTTS() *AliTTS {
	return new(AliTTS)
}

func (a *AliTTS) StreamAsk(ctx context.Context, text string) (io.ReadCloser, error) {
	return alitts.Inst().StreamAsk(ctx, text)
}
