package alitts

import (
	"fmt"
	"go-aigc-agent-demo/pkg/httputil"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
)

var client *Client

func Inst() *Client {
	return client
}

type Client struct {
	scheme       string
	serverHost   string
	serverPort   string
	client       *httputil.Client
	streamAskAPI streamAskAPI
}

type streamAskAPI struct {
	urlPath string
	appkey  string
	token   string
}

func Init(url string, appkey string, token string) error {
	scheme, hostName, port, err := httputil.ParseUrl(url)
	if err != nil {
		return fmt.Errorf("[httputil.ParseUrl]%w", err)
	}

	client = &Client{
		scheme:     scheme,
		serverHost: hostName,
		serverPort: port,
		client:     httputil.NewClient(scheme, hostName, port),
		streamAskAPI: streamAskAPI{
			urlPath: url,
			appkey:  appkey,
			token:   token,
		},
	}

	if err := httputil.WarmUpConnectionPool(client.client, url); err != nil {
		logger.Error("[httputil.WarmUpConnectionPool] failed.", slog.Any("err", err))
	}
	return nil
}
