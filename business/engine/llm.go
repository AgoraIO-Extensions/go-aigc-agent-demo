package engine

import (
	"context"
	"errors"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type llmResult struct {
	ctx     context.Context
	sid     int64
	sgid    int64
	segChan <-chan string
}

func (e *Engine) ProcessLLM(input <-chan *sentenceGroupText, output chan<- *llmResult) {
	for {
		sGroupText := <-input
		ctx := sGroupText.ctx
		sid := sGroupText.sid
		sgid := sGroupText.sgid
		groupText := sGroupText.text
		segChan, err := e.llm.Ask(ctx, sid, groupText)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("[llm] Interrupted while requesting LLM.", slog.Any("msg", err), sentencelifecycle.Tag(sid, sgid))
				continue
			}
			logger.Error("[llm.Ask]fail", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
			continue
		}
		output <- &llmResult{
			ctx:     ctx,
			sid:     sid,
			sgid:    sgid,
			segChan: segChan,
		}
	}
}
