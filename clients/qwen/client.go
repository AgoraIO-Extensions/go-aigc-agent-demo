package qwen

import (
	"go-aigc-agent-demo/pkg/httputil"
	"net/http"
	"time"
)

var client *Client

func Inst() *Client {
	return client
}

type Client struct {
	client    *httputil.Client
	streamAsk streamAsk
}

type streamAsk struct {
	url    string
	apiKey string
}

func Init(url, apikey string) error {
	transport := &http.Transport{
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 3,                // 每个主机的最大空闲连接数
		IdleConnTimeout:     90 * time.Second, // 空闲连接的超时时间
	}

	client = &Client{
		client: httputil.NewClient(transport),
		streamAsk: streamAsk{
			url:    url,
			apiKey: apikey,
		},
	}
	return httputil.WarmUpConnectionPool(client.client, url)
}
