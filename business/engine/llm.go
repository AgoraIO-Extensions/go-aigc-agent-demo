package engine

import (
	"context"
	"errors"
	"go-aigc-agent-demo/business/aigcCtx"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type llmResult struct {
	ctx     *aigcCtx.AIGCContext
	segChan <-chan string
}

func (e *Engine) ProcessLLM(input <-chan *sentenceGroupText, output chan<- *llmResult) {
	for {
		sGroupText := <-input
		ctx := sGroupText.ctx
		groupText := sGroupText.text
		segChan, err := e.llm.Ask(ctx, groupText)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.InfoContext(ctx, "[llm] Interrupted while requesting LLM.", slog.Any("msg", err))
				continue
			}
			logger.ErrorContext(ctx, "[llm.Ask]fail", slog.Any("err", err))
			continue
		}
		output <- &llmResult{
			ctx:     ctx,
			segChan: segChan,
		}
	}
}
