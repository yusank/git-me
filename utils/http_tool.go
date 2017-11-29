package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"git-me/common"
)

// GetDecodeHTML request url and read body
func GetDecodeHTML(url string) []byte {
	response, err := Response(url, false)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	return body
}

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
