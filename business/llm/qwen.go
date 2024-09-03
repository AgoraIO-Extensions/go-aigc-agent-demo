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
	"go.uber.org/zap"
	"io"
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
		return nil, fmt.Errorf("流式请求qwen失败.%w", err)
	}

	result := make(chan string, 1000)
	qw.streamRead(sid, readCloser, result, startTime)

	return result, nil
}

func (qw *QWen) streamRead(questionID int64, readCloser io.ReadCloser, result chan<- string, startTime time.Time) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Inst().Error("[panic]", zap.Any("panic msg", r), zap.Int64("sid", questionID))
				return
			}
			close(result)
			readCloser.Close()
		}()

		scanner := bufio.NewScanner(readCloser)
		defer func() {
			err := scanner.Err()
			if errors.Is(err, context.Canceled) {
				logger.Inst().Info("[llm] 流式读取返回结果是被打断", zap.String("msg", err.Error()), zap.Int64("sid", questionID))
				return
			}
			if err != nil {
				logger.Inst().Error("scanner.Err()报错", zap.Error(err), zap.Int64("sid", questionID))
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
					logger.Inst().Error("解析sse data事件的数据失败", zap.Error(err), zap.Int64("sid", questionID))
					return
				}
				choices := respData.Output.Choices

				if len(choices) == 0 {
					continue
				}
				content := choices[0].Message.Content

				logger.Inst().Info("[llm] 收到llm返回的content", zap.Int64("sid", questionID), zap.String("content", choices[0].Message.Content))
				if content == "" {
					logger.Inst().Info("返回空content", zap.Int64("sid", questionID))
					continue
				}
				if isFirstContent {
					logger.Inst().Info("[llm] 收到首个content的耗时", zap.Int64("dur", time.Since(startTime).Milliseconds()), zap.Int64("sid", questionID))
					isFirstContent = false
				}

				switch qw.ClauseMode {
				case config.NoClause:
					if isFirstContent {
						logger.Inst().Info("[llm] 收到首个segment的耗时", zap.Int64("dur", time.Since(startTime).Milliseconds()), zap.Int64("sid", questionID))
					}
					result <- content
				case config.PunctuationClause:
					segment, send, interrut := qw.SetSegmentByPunctuation(questionID, content, seg, result, false)
					if interrut {
						return
					}
					if send && isFirstSegment {
						logger.Inst().Info("[llm] 收到首个segment的耗时", zap.Int64("dur", time.Since(startTime).Milliseconds()), zap.Int64("sid", questionID))
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
			logger.Inst().Debug("[llm] 生成segment", zap.Int64("sid", sid), zap.String("seg", seg))
			seg = ""
		}
	}
	return seg, send, false
}
