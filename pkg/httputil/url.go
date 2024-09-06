package httputil

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

func ParseUrl(rawURL string) (scheme string, hostname string, port string, err error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	scheme = parsedURL.Scheme
	if scheme != "http" && scheme != "https" {
		err = fmt.Errorf("wrong scheme:%s", scheme)
		return
	}
	host := parsedURL.Host
	hostname, port, err = net.SplitHostPort(host)
	if err == nil || !strings.Contains(err.Error(), "missing port") {
		return
	}
	err = nil
	hostname = host
	if scheme == "http" {
		port = "80"
	}
	if scheme == "https" {
		port = "443"
	}
	return
}
