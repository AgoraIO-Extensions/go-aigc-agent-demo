package httputil

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-aigc-agent-demo/pkg/logger"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(scheme, hostName, port string) *Client {
	transport, err := preConnect(scheme, hostName, port)
	if err != nil {
		logger.Error("failed to exec [preConnect], will use default http client", slog.Any("err", err))
		return &Client{&http.Client{}}
	}
	return &Client{client: &http.Client{Transport: transport}}
}

// preConnect Establish an HTTP/HTTPS/HTTP2 connection to the specified host and port, and add it to the connection pool in advance.
func preConnect(scheme, hostName, port string) (*http.Transport, error) {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       nil,
	}

	once := sync.Once{}

	switch scheme {
	case "http":
		tcpConn, err := newTCPConn(context.Background(), hostName, port)
		if err != nil {
			return nil, fmt.Errorf("[newTCPConn]%v", err)
		}
		transport.DialContext = func(ctx context.Context, network, addr string) (c net.Conn, err error) {
			once.Do(func() {
				c = tcpConn
			})
			if c != nil {
				logger.Info("use pre http conn")
				return c, nil
			}
			return newTCPConn(ctx, hostName, port)
		}
	case "https":
		tlsConn, err := newTlsConn(context.Background(), hostName, port)
		if err != nil {
			return nil, fmt.Errorf("[newTlsConn]%v", err)
		}
		transport.DialTLSContext = func(ctx context.Context, network, addr string) (c net.Conn, err error) {
			once.Do(func() {
				c = tlsConn
			})
			if c != nil {
				logger.Info("use pre https conn")
				return c, nil
			}
			return newTlsConn(ctx, hostName, port)
		}
	default:
		return nil, fmt.Errorf("wrong scheme")
	}

	return transport, nil
}

func newTCPConn(ctx context.Context, hostName, port string) (net.Conn, error) {
	begin := time.Now()
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	logger.Info(fmt.Sprintf("[http]<duration> establish connection to %s:%s cost:%dms\n", hostName, port, time.Since(begin).Milliseconds()))
	return dialer.DialContext(ctx, "tcp", hostName+":"+port)
}

func newTlsConn(ctx context.Context, hostName, port string) (net.Conn, error) {
	begin := time.Now()
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	tcpConn, err := dialer.DialContext(ctx, "tcp", hostName+":"+port)
	if err != nil {
		return nil, fmt.Errorf("[newTCPConn]%v", err)
	}
	tlsConfig := &tls.Config{
		ServerName: hostName,
		NextProtos: []string{"h2", "http/1.1"},
	}
	tlsConn := tls.Client(tcpConn, tlsConfig)
	if err = tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("[tlsConn.Handshake]%v", err)
	}
	logger.Info(fmt.Sprintf("[https]<duration> establish connection to %s:%s cost:%dms\n", hostName, port, time.Since(begin).Milliseconds()))
	return tlsConn, nil
}
