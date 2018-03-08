package utils

import (
	"fmt"
	"net/http"
	"net/url"
)

// init http-client
func InitHttpClient() {
	httpClient = &http.Client{}

}

func SetProxy(port int) {
	if port == 0 {
		return
	}

	localUrl := fmt.Sprintf("http://127.0.0.1:%d", port)
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(localUrl)
	}

	// set agency
	transport := &http.Transport{
		Proxy: proxy,
	}

	httpClient.Transport = transport
}
