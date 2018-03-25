package utils

import (
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	netURL "net/url"
	"time"
)

var (
	// HttpProxy HTTP proxy
	HttpProxy string
	// Socks5Proxy SOCKS5 proxy
	Socks5Proxy string
)

// init http-client
func InitHttpClient() {
	httpClient = &http.Client{
		Timeout:   time.Second * 100,
		}
	initProxy()
}

func initProxy() {
	transport := &http.Transport{
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if HttpProxy != "" {
		var httpProxy, err = netURL.Parse(HttpProxy)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(httpProxy)
	}

	if Socks5Proxy != "" {
		dialer, err := proxy.SOCKS5(
			"tcp",
			Socks5Proxy,
			nil,
			&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			},
		)
		if err != nil {
			panic(err)
		}
		transport.Dial = dialer.Dial
	}

	httpClient.Transport = transport
}
