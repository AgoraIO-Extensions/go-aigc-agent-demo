package engine

import (
	"context"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type ttsResult struct {
	ctx   context.Context
	sid   int64
	sgid  int64
	audio <-chan []byte
}

func (e *Engine) ProcessTTS(input <-chan *llmResult, output chan<- *ttsResult) {
	for {
		r := <-input
		sid, sgid := r.sid, r.sgid
		ctx := r.ctx

		/* create a TTS client connection */
		ttsClient, err := e.ttsFactory.CreateTTS(sid)
		if err != nil {
			logger.Error("[tts] Failed to create TTS client instance.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
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
					logger.Info("[tts] The process of sending a segment to TTS was interrupted.", sentencelifecycle.Tag(sid))
					return
				}
			}
		}()

		output <- &ttsResult{
			ctx:   ctx,
			sid:   sid,
			sgid:  sgid,
			audio: ttsClient.GetResult(),
		}
	}
}
