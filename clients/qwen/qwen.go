package qwen

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/pkg/logger"
	"go.uber.org/zap"
	"io"
)

type Input struct {
	Messages []Msg `json:"messages"`
}

type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Parameters struct {
	ResultFormat      string `json:"result_format"`
	IncrementalOutput bool   `json:"incremental_output"`
}

type Output struct {
	Choices []Choice `json:"choices"`
}

type FinishReason string

const (
	FinishNull   FinishReason = "null" // 生成过程中
	FinishStop   FinishReason = "stop" // stop token导致结束
	FinishLength FinishReason = "null" // 生成长度导致结束
)

type Choice struct {
	Message      Msg          `json:"message"`
	FinishReason FinishReason `json:"finish_reason"`
}

type Usage struct {
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

/* --------------------------------------------------------------------------------------------------------------------- */

type SSEReq struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type SSEResp struct {
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestID string `json:"request_id"`
}

func (cli *Client) StreamAsk(ctx context.Context, model string, msgs []Msg, sid int64) (io.ReadCloser, error) {
	header := map[string]string{
		"Authorization":   "Bearer " + cli.streamAsk.apiKey,
		"Content-Type":    "application/json",
		"X-DashScope-SSE": "enable",
	}
	req := SSEReq{
		Model: model,
		Input: Input{Messages: msgs},
		Parameters: Parameters{
			ResultFormat:      "message",
			IncrementalOutput: true,
		},
	}

	resp, err := cli.client.JSONPost(ctx, cli.streamAsk.url, req, header)
	if err != nil {
		return nil, fmt.Errorf("[cli.client.JSONPost]%w", err)
	}
	if resp.StatusCode/100 != 2 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Inst().Error("io.ReadAll报错", zap.Error(err), zap.Int64("sid", sid))
		}
		return nil, fmt.Errorf("服务端返回错误响应, statuscode:%d, resp.body:%s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}
