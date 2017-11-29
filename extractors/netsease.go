package extractors

import (
	"fmt"
	"strings"

	"git-me/utils"
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

	data := string(utils.GetDecodeHTML(url))
	if len(data) == 0 {
		fmt.Println("data is nil")
		return
	}

	fmt.Println(data)
}
