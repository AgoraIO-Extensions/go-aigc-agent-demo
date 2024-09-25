package qwen

import (
	"context"
	"fmt"
	"go-aigc-agent-demo/pkg/logger"
	"io"
	"log/slog"
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

type Choice struct {
	Message      Msg    `json:"message"`
	FinishReason string `json:"finish_reason"`
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

func (cli *Client) StreamAsk(ctx context.Context, model string, msgs []Msg) (io.ReadCloser, error) {
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
			logger.ErrorContext(ctx, "[io.ReadAll]", slog.Any("err", err))
		}
		return nil, fmt.Errorf("server returned an error, statuscode:%d, resp.body:%s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}
