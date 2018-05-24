package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	httpClient *http.Client

	// Cookie http cookies
	Cookie string

	// FakeHeader for when web sites checke request header
	FakeHeader = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Charset":  "UTF-8,*;q=0.5",
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": "en-US,en;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36",
	}
)

// GetDecodeHTML request url and read body
func GetDecodeHTML(url string, header map[string]string) []byte {
	response, err := HttpGet(url, header)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	if response.Header.Get("Content-Encoding") == "gzip" {
		response.Body, _ = gzip.NewReader(response.Body)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	return body
}

// Headers return the HTTP Headers of the url
func Headers(url, refer string) http.Header {
	headers := map[string]string{
		"Referer": refer,
	}
	res, _ := HttpGet(url, headers)
	defer res.Body.Close()
	return res.Header
}

// Size get size of the url
func DownloadFileSize(url, refer string) int64 {
	h := Headers(url, refer)
	s := h.Get("Content-Length")
	size, _ := strconv.ParseInt(s, 10, 64)
	return size
}

// FilePath gen valid file path
func FilePath(name, ext, output string, escape bool) string {
	var outputPath string
	if output != "" {
		_, err := os.Stat(output)
		if err != nil && os.IsNotExist(err) {
			log.Println("found err", output)
			log.Fatal(err)
		}
	}
	fileName := fmt.Sprintf("%s.%s", name, ext)
	if escape {
		fileName = FileName(fileName)
	}
	outputPath = filepath.Join(output, fileName)
	return outputPath
}

// FileSize return the file size of the specified path file
func FileSize(filePath string) (int64, bool) {
	file, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return 0, false
	}
	return file.Size(), true
}

// HttpGetByte - get http response
func HttpGet(url string, header map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if Cookie != "" {
		var cookie string
		if _, fileErr := os.Stat(Cookie); fileErr == nil {
			// Cookie is a file
			data, _ := ioutil.ReadFile(Cookie)
			cookie = string(data)
		} else {
			// Just strings
			cookie = Cookie
		}
		req.Header.Set("Cookie", cookie)
	}

	for k, v := range FakeHeader {
		req.Header.Set(k, v)
	}
	req.Header.Set("Referer", url)

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
		resp, err = HttpGet(url, header)
		if err != nil || resp == nil {
			time.Sleep(500 * time.Millisecond)
			err = nil
		}
	}
	return
}

func GetRequestStr(url string, refer string) string {
	headers := map[string]string{}
	if refer != "" {
		headers["Referer"] = refer
	}
	resp, err := HttpGet(url, headers)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := DecodeResp(resp)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

// HttpGetByte -
func HttpGetByte(url string, header map[string]string) ([]byte, error) {
	resp, err := HttpGet(url, header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return DecodeResp(resp)
}

// DecodeResp -
func DecodeResp(resp *http.Response) ([]byte, error) {
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	} else {
		reader = resp.Body
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpPost(url string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err = httpClient.Do(req)
	return
}
