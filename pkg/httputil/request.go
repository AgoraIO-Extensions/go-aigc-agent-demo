package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func WarmUpConnectionPool(client *Client, url string) error {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return fmt.Errorf("[http.NewRequest]%w", err)
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return fmt.Errorf("[client.Do]%w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) JSONPost(ctx context.Context, url string, reqStruct interface{}, headers map[string]string) (*http.Response, error) {
	var err error

	reqBody := make([]byte, 0)
	if reqStruct != nil {
		reqBody, err = json.Marshal(reqStruct)
		if err != nil {
			err = fmt.Errorf("[Marshal]%w", err)
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		err = fmt.Errorf("[http.NewRequest]%w", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		err = fmt.Errorf("[client.Do]%w", err)
		return nil, err
	}
	return resp, err
}
