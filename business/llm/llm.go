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
		return nil, fmt.Errorf("vendorName传参错误:%s", vendorName)
	}
	return &LLM{prompt: prompt, dCTX: dCTX, streamClient: client}, nil
}

func (l *LLM) Ask(ctx context.Context, sid int64, question string) (<-chan string, error) {
	sgid := sentencelifecycle.GroupInst().GetSgidBySid(sid)
	msgs := l.dCTX.AddQuestion(question, sgid) // 构建并返回 当前的 上下文信息链
	if l.prompt != "" {
		msgs = append([]dialogctx.Message{{Role: dialogctx.SYSTEM, Content: l.prompt}}, msgs...)
	}
	logger.Info("[llm] 带上下文的提问", slog.Any("dialog_ctx", msgs), sentencelifecycle.Tag(sid))
	segChan, err := l.streamClient.StreamAsk(ctx, sid, msgs)
	if err != nil {
		return nil, fmt.Errorf("[streamAsk]%w", err)
	}
	logger.Info("[llm] llm请求返回响应头", sentencelifecycle.Tag(sid))

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
					这里之所以流式地将answer追加到dCTX（而不是等全部返回后一次性添加到dCTX），是因为在被「并句」的时候，当前已经返
					回的内容可能已经被用户听到了，用户的下一句话可能是基于这部分已返回的内容进行提问的，所以必须将这部分answer及时地添加到dCtx中
				*/
				if err = l.dCTX.StreamAddAnswer(seg, sgid); err != nil {
					logger.Error("[llm] 流式追加回答失败", slog.Any("err", err), sentencelifecycle.Tag(sid, sgid))
					return
				}
			case <-ctx.Done():
				logger.Info("[llm] 流式读取llm返回文本时被打断", sentencelifecycle.Tag(sid))
				return
			}
		}
	}()

	return segChanCopy, nil
}
