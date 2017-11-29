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

	title := utils.Match(`movieDescription=\'([^\']+)\'`, data)
	if len(title) == 0 {
		title = utils.Match(`<title>(.+)</title>`, data)
	}

	if len(title) > 1 && title[0] == " " {
		title = title[1:]
	}

	src := utils.Match(`<source src="([^"]+)"`, data)
	if len(src) == 0 {
		src = utils.Match(`<source type="[^"]+" src="([^"]+)"`, data)
	}

	urls := []string{}
	exc := ""
	size := 0
	if len(src) > 0 {
		urls = src
		exc = "mp4"
		size = 1
		// url_info func, return ext, size

	} else {
		urls = utils.Match(`["\\'](.+)-list.m3u8["\\']`, data)
		if len(urls) == 0 {
			urls = utils.Match(`["\\'](.+).m3u8["\\']`, data)
		}

		size = 2
		exc = "mp4"
	}

	// todo:Print url_info
	// download func

	fmt.Println(data, exc, size)
}
