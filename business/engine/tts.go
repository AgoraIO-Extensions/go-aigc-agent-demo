package engine

import (
	"context"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type ttsResult struct {
	ctx   context.Context
	audio <-chan []byte
}

func (e *Engine) ProcessTTS(input <-chan *llmResult, output chan<- *ttsResult) {
	for {
		r := <-input
		ctx := r.ctx

		/* create a TTS client connection */
		ttsClient, err := e.ttsFactory.CreateTTS(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "[tts] Failed to create TTS client instance.", slog.Any("err", err))
			continue
		}

		/* send the LLM result to TTS */
		go func() {
			for i := 0; ; i++ {
				select {
				case seg, ok := <-r.segChan:
					if !ok {
						ttsClient.Send(ctx, i, "")
						return
					}
					ttsClient.Send(ctx, i, seg)
				case <-ctx.Done(): // interrupted
					logger.InfoContext(ctx, "[tts] The process of sending a segment to TTS was interrupted.")
					return
				}
			}
		}()

		output <- &ttsResult{
			ctx:   ctx,
			audio: ttsClient.GetResult(),
		}
	}
}
