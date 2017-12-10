package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"time"
)

var (
	httpClient *http.Client
)

// GetDecodeHTML request url and read body
func GetDecodeHTML(url string, header map[string]string) []byte {
	response, err := Response(url, header)
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
func Response(url string, header map[string]string) (*http.Response, error) {
	//httpClient := &http.Client{}

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

func RequestWithRetry(url string, header map[string]string) (resp *http.Response, err error) {
	for i := 0; i < 3; i++ {
		resp, err = Response(url, header)
		if err != nil || resp == nil {
			time.Sleep(500 * time.Millisecond)
			err = nil
		}
	}
	return
}

// GetContent -
func GetContent(url string, header map[string]string) ([]byte, error) {
	fmt.Printf("GetContent:%s\n", url)

	resp, err := Response(url, header)
	if err != nil {
		return nil, err
	}

	return DecodeResp(resp)
}

// DecodeResp -
func DecodeResp(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(body))
	return body, nil
}
