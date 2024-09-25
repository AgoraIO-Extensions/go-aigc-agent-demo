package llm

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-aigc-agent-demo/business/llm/common/clause"
	"go-aigc-agent-demo/business/llm/common/dialogctx"
	qwenCli "go-aigc-agent-demo/clients/qwen"
	"go-aigc-agent-demo/config"
	"go-aigc-agent-demo/pkg/logger"
	"io"
	"log/slog"
	"strings"
	"time"
)

type QWen struct {
	Model      string
	ClauseMode config.ClauseMode
}

func NewQWen(modelName string, clauseMode config.ClauseMode) *QWen {
	return &QWen{
		Model:      modelName,
		ClauseMode: clauseMode,
	}
}

func (qw *QWen) StreamAsk(ctx context.Context, llmMsgs []dialogctx.Message) (res <-chan string, err error) {
	var qwenMsgs []qwenCli.Msg
	for _, m := range llmMsgs {
		qwenMsgs = append(qwenMsgs, qwenCli.Msg{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}

	startTime := time.Now()
	readCloser, err := qwenCli.Inst().StreamAsk(ctx, qw.Model, qwenMsgs)
	if err != nil {
		return nil, fmt.Errorf("failed to request qwen in a streaming manner.%w", err)
	}

	result := make(chan string, 1000)
	qw.streamRead(ctx, readCloser, result, startTime)

	return result, nil
}

func (qw *QWen) streamRead(ctx context.Context, readCloser io.ReadCloser, result chan<- string, startTime time.Time) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx, "[panic]", slog.Any("panic msg", r))
				return
			}
			close(result)
			readCloser.Close()
		}()

		scanner := bufio.NewScanner(readCloser)
		defer func() {
			err := scanner.Err()
			if errors.Is(err, context.Canceled) {
				logger.InfoContext(ctx, "[llm] Interrupted while reading the qwen result in a streaming manner.", slog.String("msg", err.Error()))
				return
			}
			if err != nil {
				logger.ErrorContext(ctx, "scanner.Err()", slog.Any("err", err))
			}
		}()

		var isFirstContent = true
		var isFirstSegment = true
		var seg string

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimPrefix(line, "data:")
				data = strings.TrimSpace(data)
				var respData qwenCli.SSEResp
				if err := json.Unmarshal([]byte(data), &respData); err != nil {
					logger.ErrorContext(ctx, "[llm] Failed to parse the data from the SSE event", slog.Any("err", err))
					return
				}
				choices := respData.Output.Choices

				if len(choices) == 0 {
					continue
				}
				content := choices[0].Message.Content

				logger.InfoContext(ctx, "[llm] Received content returned by the LLM", slog.String("content", choices[0].Message.Content))
				if content == "" {
					logger.InfoContext(ctx, "[llm] Returned empty content")
					continue
				}
				if isFirstContent {
					logger.InfoContext(ctx, "[llm] Time taken to receive the first content", slog.Int64("dur", time.Since(startTime).Milliseconds()))
					isFirstContent = false
				}

				switch qw.ClauseMode {
				case config.NoClause:
					if isFirstSegment {
						logger.InfoContext(ctx, "[llm] Time taken to receive the first segment", slog.Int64("dur", time.Since(startTime).Milliseconds()))
					}
					result <- content
				case config.PunctuationClause:
					segment, send, interrut := qw.SetSegmentByPunctuation(ctx, content, seg, result, false)
					if interrut {
						return
					}
					if send && isFirstSegment {
						logger.Info("[llm] Time taken to receive the first segment", slog.Int64("dur", time.Since(startTime).Milliseconds()))
						isFirstSegment = false
					}
					seg = segment
				}
			}
		}
		if qw.ClauseMode == config.PunctuationClause {
			qw.SetSegmentByPunctuation(ctx, "", seg, result, true)
		}

	}()
}

func (qw *QWen) SetSegmentByPunctuation(ctx context.Context, streamText, seg string, result chan<- string, end bool) (segment string, send bool, interrupt bool) {
	if end && seg != "" {
		result <- seg
		return "", true, false
	}
	// 遇到标点就分句
	for _, char := range streamText {
		seg = seg + string(char)
		if clause.CharMap[char] {
			result <- seg
			send = true
			logger.InfoContext(ctx, "[llm] Generate segment", slog.String("seg", seg))
			seg = ""
		}
	}
	return seg, send, false
}
