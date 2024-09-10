package alitts

import (
	"context"
	"fmt"
	"io"
)

func (cli *Client) StreamAsk(ctx context.Context, text string) (io.ReadCloser, error) {
	sa := cli.streamAskAPI
	bodyContent := make(map[string]interface{})
	bodyContent["appkey"] = sa.appkey
	bodyContent["text"] = text
	bodyContent["token"] = sa.token
	bodyContent["format"] = "pcm"
	bodyContent["sample_rate"] = 16000
	bodyContent["voice"] = "zhitian_emo"
	// volume: 0~100，default:50。
	bodyContent["volume"] = 50
	// speech_rate: -500~500，default:0。
	bodyContent["speech_rate"] = 0
	// pitch_rate: -500~500，default:0。
	bodyContent["pitch_rate"] = 0

	resp, err := cli.client.JSONPost(ctx, sa.urlPath, bodyContent, nil)
	if err != nil {
		return nil, fmt.Errorf("[cli.client.JSONPost]%w", err)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ali-tts server returned an error, statuscode:%d, err:%s", resp.StatusCode, string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "audio/mpeg" {
		return resp.Body, nil
	}
	// ContentType is "null" or "application/json"
	statusCode := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("ali-tts server returned an error, statuscode:%d, err:%s", statusCode, string(body))
}
