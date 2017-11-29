package utils

import (
	"git-me/common"
	"net/http"
)

// Response - get http response
func Response(url string, isFake bool) (*http.Response, error) {
	httpClient := &http.Client{}
	header := make(map[string]string)
	if isFake {
		header = common.FakeHeader
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
