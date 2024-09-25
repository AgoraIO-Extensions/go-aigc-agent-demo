package llm

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"go-aigc-agent-demo/business/llm/common/clause"
	"go-aigc-agent-demo/business/llm/common/dialogctx"
	ms_chat_gpt "go-aigc-agent-demo/pkg/azureopenai/chat-gpt"
	"go-aigc-agent-demo/pkg/logger"
	"io"
	"log/slog"
	"time"
)

type ChatGPT struct {
	modelName string
}

func NewChatGPT(modelName string) *ChatGPT {
	return &ChatGPT{
		modelName: modelName,
	}
}

func (gpt *ChatGPT) StreamAsk(ctx context.Context, llmMsgs []dialogctx.Message) (segChan <-chan string, err error) {
	var chatGptMsgs []ms_chat_gpt.Msg
	for _, m := range llmMsgs {
		chatGptMsgs = append(chatGptMsgs, ms_chat_gpt.Msg{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}
	startTime := time.Now()
	resp, err := ms_chat_gpt.Inst().StreamAsk(chatGptMsgs, gpt.modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to request ms_chat_gpt in a streaming manner.%w", err)
	}

	segmentChan := make(chan string, 1000)
	gpt.streamRead(ctx, resp, segmentChan, startTime)

	return segmentChan, nil

}

// streamRead 流式读取数据put到 segChan，退出时close segChan
func (gpt *ChatGPT) streamRead(ctx context.Context, resp azopenai.GetChatCompletionsStreamResponse, segChan chan<- string, startTime time.Time) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx, "[panic]", slog.Any("panic msg", r))
				return
			}
			close(segChan)
			resp.ChatCompletionsStream.Close()
		}()

		segment := ""
		var isFirstContent = true
		var isFirstSegment = true
		for {
			chatCompletions, err := resp.ChatCompletionsStream.Read()
			if errors.Is(err, io.EOF) {
				if segment != "" {
					segChan <- segment
				}
				break
			}
			if err != nil {
				logger.ErrorContext(ctx, "Failed to read data from chat-gpt in a streaming manner", slog.Any("err", err))
				return
			}

			if isFirstContent {
				logger.InfoContext(ctx, "[llm] Time taken to receive the first content", slog.Int64("dur", time.Since(startTime).Milliseconds()))
				isFirstContent = false
			}

			for _, choice := range chatCompletions.Choices {
				if choice.Delta.Content == nil {
					continue
				}
				for _, char := range *choice.Delta.Content {
					segment = segment + string(char)
					if clause.CharMap[char] {
						segChan <- segment
						if isFirstSegment {
							logger.InfoContext(ctx, "[llm] Time taken to receive the first segment", slog.Int64("dur", time.Since(startTime).Milliseconds()))
							isFirstSegment = false
						}
						segment = ""
					}
				}
			}
		}
	}()
}
