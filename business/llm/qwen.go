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

func (qw *QWen) StreamAsk(ctx context.Context, sid int64, llmMsgs []dialogctx.Message) (res <-chan string, err error) {
	var qwenMsgs []qwenCli.Msg
	for _, m := range llmMsgs {
		qwenMsgs = append(qwenMsgs, qwenCli.Msg{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}

	startTime := time.Now()
	readCloser, err := qwenCli.Inst().StreamAsk(ctx, qw.Model, qwenMsgs, sid)
	if err != nil {
		return nil, fmt.Errorf("failed to request qwen in a streaming manner.%w", err)
	}

	result := make(chan string, 1000)
	qw.streamRead(sid, readCloser, result, startTime)

	return result, nil
}

func (qw *QWen) streamRead(questionID int64, readCloser io.ReadCloser, result chan<- string, startTime time.Time) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("[panic]", slog.Any("panic msg", r), slog.Int64("sid", questionID))
				return
			}
			close(result)
			readCloser.Close()
		}()

		scanner := bufio.NewScanner(readCloser)
		defer func() {
			err := scanner.Err()
			if errors.Is(err, context.Canceled) {
				logger.Info("[llm] Interrupted while reading the qwen result in a streaming manner.", slog.String("msg", err.Error()), slog.Int64("sid", questionID))
				return
			}
			if err != nil {
				logger.Error("scanner.Err()", slog.Any("err", err), slog.Int64("sid", questionID))
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
					logger.Error("[llm] Failed to parse the data from the SSE event", slog.Any("err", err), slog.Int64("sid", questionID))
					return
				}
				choices := respData.Output.Choices

				if len(choices) == 0 {
					continue
				}
				content := choices[0].Message.Content

				logger.Info("[llm] Received content returned by the LLM", slog.Int64("sid", questionID), slog.String("content", choices[0].Message.Content))
				if content == "" {
					logger.Info("[llm] Returned empty content", slog.Int64("sid", questionID))
					continue
				}
				if isFirstContent {
					logger.Info("[llm] Time taken to receive the first content", slog.Int64("dur", time.Since(startTime).Milliseconds()), slog.Int64("sid", questionID))
					isFirstContent = false
				}

				switch qw.ClauseMode {
				case config.NoClause:
					if isFirstSegment {
						logger.Info("[llm] Time taken to receive the first segment", slog.Int64("dur", time.Since(startTime).Milliseconds()), slog.Int64("sid", questionID))
					}
					result <- content
				case config.PunctuationClause:
					segment, send, interrut := qw.SetSegmentByPunctuation(questionID, content, seg, result, false)
					if interrut {
						return
					}
					if send && isFirstSegment {
						logger.Info("[llm] Time taken to receive the first segment", slog.Int64("dur", time.Since(startTime).Milliseconds()), slog.Int64("sid", questionID))
						isFirstSegment = false
					}
					seg = segment
				}
			}
		}
		if qw.ClauseMode == config.PunctuationClause {
			qw.SetSegmentByPunctuation(questionID, "", seg, result, true)
		}

	}()
}

func (qw *QWen) SetSegmentByPunctuation(sid int64, streamText, seg string, result chan<- string, end bool) (segment string, send bool, interrupt bool) {
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
			logger.Info("[llm] Generate segment", slog.Int64("sid", sid), slog.String("seg", seg))
			seg = ""
		}
	}
	return seg, send, false
}
