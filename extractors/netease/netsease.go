package netease

import (
	"fmt"
	"strings"

	"git-me/common"
	"git-me/utils"
)

type BasicInfo struct {
	Url string
}

func (wy BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

func (wy BasicInfo) Download(url string) (common.VideoData, error) {
	return DownloadByURL(url)
}

// DownloadByURL -
func DownloadByURL(url string) (vid common.VideoData, err error) {
	if strings.Contains(url, "163.fm") {
		fmt.Println(url)
		return
	}

	if strings.Contains(url, "music.163.com") {
		fmt.Println(url)
		return
	}

	data := string(utils.GetDecodeHTML(url, nil))
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
	ext := ""
	size := 0

	if len(src) > 0 {
		for _, v := range src {
			u := common.URLData{
				URL: v,
				Ext: "mp4",
			}
			vid.Type = "video"
			vid.URLs = append(vid.URLs, u)
			vid.Size++
		}
	} else {
		urls = utils.Match(`["\\'](.+)-list.m3u8["\\']`, data)
		if len(urls) == 0 {
			urls = utils.Match(`["\\'](.+).m3u8["\\']`, data)
		}
		for _, v := range urls {
			u := common.URLData{
				URL: v,
				Ext: "mp4",
			}

			vid.Type = "video"
			vid.URLs = append(vid.URLs, u)
			vid.Size++
		}
	}
	fmt.Println(urls)
	fmt.Println(data, ext, size)

	return
}
