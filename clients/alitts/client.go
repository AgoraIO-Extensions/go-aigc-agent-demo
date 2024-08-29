package alitts

import (
	"go-aigc-agent-demo/pkg/httputil"
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
	urlPath string
	appkey  string
	token   string
}

func Init(url string, appkey string, token string) error {
	client = &Client{
		client: httputil.NewClient(nil),
		streamAsk: streamAsk{
			urlPath: url,
			appkey:  appkey,
			token:   token,
		},
	}

	return httputil.WarmUpConnectionPool(client.client, url)
}
