package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	client      *http.Client
	serviceName string
	clientName  string
}

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

func NewClient(transport *http.Transport) (c *Client) {
	if transport == nil {
		transport = &http.Transport{
			MaxIdleConns:        100,              // 最大空闲连接数
			MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
			IdleConnTimeout:     90 * time.Second, // 空闲连接的超时时间
		}
	}

	return &Client{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second, // 客户端的请求超时时间
		},
	}
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
