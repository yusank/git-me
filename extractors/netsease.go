package extractors

import (
	"fmt"
	"git-me/utils"
	"io/ioutil"
	"strings"
)

// DownloadByURL -
func DownloadByURL(url string) {
	if strings.Contains(url, "163.fm") {
		fmt.Println(url)
		return
	}

	if strings.Contains(url, "music.163.com") {
		fmt.Println(url)
		return
	}

	data := string(GetDecodeHTML(url))
	if len(data) == 0 {
		fmt.Println("data is nil")
		return
	}

	fmt.Println(data)
}

// GetDecodeHTML request url and read body
func GetDecodeHTML(url string) []byte {
	response, err := utils.Response(url, false)
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
