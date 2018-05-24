/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/04/01        Yusan Kurban
 */

package utils

import (
	"net"
	"net/http"
	netURL "net/url"
	"time"

	"golang.org/x/net/proxy"
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
		Timeout: time.Second * 100,
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
