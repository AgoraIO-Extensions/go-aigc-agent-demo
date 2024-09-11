package qwen

import (
	"fmt"
	"go-aigc-agent-demo/pkg/httputil"
)

var client *Client

func Inst() *Client {
	return client
}

type Client struct {
	scheme     string
	serverHost string
	serverPort string
	client     *httputil.Client
	streamAsk  streamAsk
}

type streamAsk struct {
	url    string
	apiKey string
}

func Init(url, apikey string) error {
	scheme, hostName, port, err := httputil.ParseUrl(url)
	if err != nil {
		return fmt.Errorf("[httputil.ParseUrl]%w", err)
	}

	client = &Client{
		scheme:     scheme,
		serverHost: hostName,
		serverPort: port,
		client:     httputil.NewClient(scheme, hostName, port),
		streamAsk: streamAsk{
			url:    url,
			apiKey: apikey,
		},
	}
	return nil
}
