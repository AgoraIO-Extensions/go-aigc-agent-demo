package llm

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/business/llm/common/dialogctx"
	"go-aigc-agent-demo/business/sentencelifecycle"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

type streamClient interface {
	StreamAsk(ctx context.Context, sid int64, llmMsgs []dialogctx.Message) (segChan <-chan string, err error)
}

type LLM struct {
	prompt       string
	dCTX         *dialogctx.DialogCTX
	streamClient streamClient
}

func NewLLM(vendorName config.ModelSelect, prompt string, cfg *config.LLM) (*LLM, error) {
	var client streamClient
	var dCTX *dialogctx.DialogCTX
	switch vendorName {
	case config.LLMQwen:
		dCTX = dialogctx.NewDialogCTX(cfg.QWen.DialogNums, cfg.WithHistory)
		client = NewQWen(cfg.QWen.Model, cfg.ClauseMode)
	case config.LLMChatGPT4o:
		dCTX = dialogctx.NewDialogCTX(cfg.ChatGPT4o.DialogNums, cfg.WithHistory)
		client = NewChatGPT(cfg.ChatGPT4o.Model)
	default:
		return nil, fmt.Errorf("vendorName parameter is incorrect:%s", vendorName)
	}
	return &LLM{prompt: prompt, dCTX: dCTX, streamClient: client}, nil
}

func (l *LLM) Ask(ctx context.Context, sid int64, question string) (<-chan string, error) {
	sgid := sentencelifecycle.GroupInst().GetSgidBySid(sid)
	msgs := l.dCTX.AddQuestion(question, sgid) // Build and return the current context information chain.
	if l.prompt != "" {
		msgs = append([]dialogctx.Message{{Role: dialogctx.SYSTEM, Content: l.prompt}}, msgs...)
	}
	logger.Info("[llm] question with dialog context", slog.Any("dialog_ctx", msgs), sentencelifecycle.Tag(sid))
	segChan, err := l.streamClient.StreamAsk(ctx, sid, msgs)
	if err != nil {
		return nil, fmt.Errorf("[streamAsk]%w", err)
	}
	logger.Info("[llm] Return response headers for LLM requests", sentencelifecycle.Tag(sid))

	segChanCopy := make(chan string, 1000)
	go func() {
		for {
			select {
			case seg, ok := <-segChan:
				if !ok {
					close(segChanCopy)
					return
				}
				segChanCopy <- seg
				/*
					The reason for appending the answer to dCTX in a streaming manner (rather than adding it all at once
					after the entire response is returned) is that when a response is 'completed' (i.e., the sentence is
					finished), the content already returned might have been heard by the user. The user's next question
					might be based on the content that has already been returned, so it is necessary to add this part of
					the answer to dCTX in a timely manner.
				*/
				if err = l.dCTX.StreamAddAnswer(seg, sgid); err != nil {
					logger.Error("[llm] Streaming append of the answer failed.", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
					return
				}
			case <-ctx.Done():
				logger.Info("[llm] Interrupted while reading the LLM's returned text in a streaming manner.", sentencelifecycle.Tag(sid))
				return
			}
		}
	}()

	return segChanCopy, nil
}
