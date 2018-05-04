package netease

import (
	"fmt"
	"strings"

	"git-me/common"
	"git-me/utils"
	"sort"
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
	if strings.Contains(url, "music.163.com") {
		fmt.Println(url)
		vid, err = CloudMusicDownload(url)
		if err != nil {
			return
		}
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

	if len(src) > 0 {
		urlData := []common.URLData{}
		for _, v := range src {
			u := common.URLData{
				URL: v,
				Ext: "mp4",
			}
			urlData = append(urlData, u)
			vid.Type = "video"
		}
		format := common.FormatData{
			URLs: urlData,
			Size: int64(len(urlData)),
		}
		vid.Formats = []common.FormatData{format}
	} else {
		urls := utils.Match(`["\\'](.+)-list.m3u8["\\']`, data)
		if len(urls) == 0 {
			urls = utils.Match(`["\\'](.+).m3u8["\\']`, data)
		}
		urlData := []common.URLData{}
		for _, v := range urls {
			u := common.URLData{
				URL: v,
				Ext: "mp4",
			}
			urlData = append(urlData, u)
			vid.Type = "video"
		}

		format := common.FormatData{
			URLs: urlData,
			Size: int64(len(urls)),
		}
		vid.Formats = []common.FormatData{format}
	}

	return
}

// CloudMusicDownload -
func CloudMusicDownload(url string) (common.VideoData, error) {
	var v common.VideoData
	rid := utils.Match(`\Wid=(.*)`, url)
	if len(rid) == 0 {
		rid = utils.Match(`/(\d+)/?`, url)
		newRid := []string{}
		for _, v := range rid {
			s := v[1:]
			newRid = append(newRid, s[:len(s)-1])
		}
		rid = newRid
		fmt.Println(rid)
	} else {
		newRid := []string{}
		for _, v := range rid {
			s := strings.Split(v, "=")
			newRid = append(newRid, s[1])
		}
		rid = newRid
	}

	header := make(map[string]string)
	header["Referer"] = "http://music.163.com/"

	switch {
	case strings.Contains(url, "mv"):
		fmt.Println("it`s mv")
		reqUrl := fmt.Sprintf("http://music.163.com/api/mv/detail/?id=%s&ids=%s&csrf_token=", rid[0], rid)

		j, err := utils.LoadJSON(reqUrl, header)
		if err != nil {
			return v, err
		}

		vinfo := j.Get("data").MustMap()
		return MvDownload(vinfo)

	case strings.Contains(url, "album"):
		j, err := utils.LoadJSON(url, header)
		if err != nil {
			return v, err
		}

		fmt.Printf("%+v \n", j)
	case strings.Contains(url, "program"):
		reqUrl := fmt.Sprintf("http://music.163.com/api/dj/program/detail/?id=%s&ids=%s&csrf_token=", rid[0], rid)
		common.FakeHeader["Referer"] = "http://music.163.com/"
		fmt.Println(common.FakeHeader)
		j, err := utils.LoadJSON(reqUrl, header)
		if err != nil {
			return v, err
		}

		songInfo, err := j.String()
		fmt.Println(songInfo)
	}
	return v, fmt.Errorf("not found")
}

func MvDownload(vinfo map[string]interface{}) (vid common.VideoData, err error) {
	//title := fmt.Sprintf("%s - %s", vinfo["name"], vinfo["artistName"])
	urlBest := vinfo["brs"].(map[string]interface{})
	var keys []string
	for k := range urlBest {
		keys = append(keys, k)
	}
	// sort and get the best quality mv
	sort.Strings(keys)

	urlData := []common.URLData{}
	for _, v := range keys {
		u := common.URLData{
			URL: urlBest[v].(string),
			Ext: "mp4",
		}
		urlData = append(urlData, u)
		vid.Type = "video"
	}

	format := common.FormatData{
		URLs: urlData,
		Size: int64(len(keys)),
	}
	vid.Formats = []common.FormatData{format}

	return
}
