package youku

import (
	"encoding/json"
	"fmt"
	"git-me/common"
	"git-me/utils"
	"log"
	"strings"
	"time"
)

func (yk BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

type errorData struct {
	Note string `json:"note"`
	Code int    `json:"code"`
}

type segs struct {
	Size int64  `json:"size"`
	URL  string `json:"cdn_url"`
}

type stream struct {
	Size   int64  `json:"size"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Segs   []segs `json:"segs"`
	Type   string `json:"stream_type"`
}

type data struct {
	Error  errorData `json:"error"`
	Stream []stream  `json:"stream"`
}

type BasicInfo struct {
	Data data `json:"data"`
}

var ccodes = []string{"0507", "0508", "0512", "0513", "0514", "0503", "0502", "0590"}
var referer = "https://v.youku.com"

func (yk BasicInfo) ups(vid string) {
	var url string
	var utid string
	var html string
	headers := utils.Headers("http://log.mmstat.com/eg.js", referer)
	setCookie := headers.Get("Set-Cookie")
	utid = utils.MatchOneOf(setCookie, `cna=(.+?);`)[1]
	for _, ccode := range ccodes {
		url = fmt.Sprintf(
			"https://ups.youku.com/ups/get.json?vid=%s&ccode=%s&client_ip=192.168.1.1&client_ts=%d&utid=%s",
			vid, ccode, time.Now().Unix(), utid,
		)
		html = string(utils.GetDecodeHTML(url, nil))
		// data must be emptied before reassignment, otherwise it will contain the previous value(the 'error' data)
		json.Unmarshal([]byte(html), &yk)
		if yk.Data.Error.Code != -6004 {
			return
		}
	}

	return
}

func (yk BasicInfo) genData() ([]common.URLData, int64, string) {
	var (
		urls  []common.URLData
		size  int64
		index int
	)
	// get the best quality
	for i, s := range yk.Data.Stream {
		if s.Size > size {
			size = s.Size
			index = i
		}
	}
	stream := yk.Data.Stream[index]
	ext := strings.Split(
		strings.Split(stream.Segs[0].URL, "?")[0],
		".",
	)
	for _, data := range stream.Segs {
		url := common.URLData{
			URL:  data.URL,
			Size: data.Size,
			Ext:  ext[len(ext)-1],
		}
		urls = append(urls, url)
	}
	quality := fmt.Sprintf("%s %dx%d", stream.Type, stream.Width, stream.Height)
	return urls, stream.Size, quality
}

// Download implement common.VideoExtractor
func (yk BasicInfo) Download(url string) (data common.VideoData, err error) {
	html := string(utils.GetDecodeHTML(url, nil))
	// get the title
	doc := utils.GetDoc(html)
	title := utils.Title(doc)
	vid := utils.MatchOneOf(url, `id_(.+?).html`)[1]
	yk.ups(vid)
	if yk.Data.Error.Code != 0 {
		log.Fatal(yk.Data.Error.Note)
	}
	urls, size, quality := yk.genData()
	data = common.VideoData{
		Site:    "优酷 youku.com",
		Title:   title,
		Type:    "video",
		URLs:    urls,
		Size:    size,
		Quality: quality,
	}

	return
}
