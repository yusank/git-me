/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/04/01        Yusan Kurban
 */

package youku

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/utils"
)

//
//type BasicInfo struct {
//	Title string
//	Url        string
//	Referer    string
//	Page       []byte
//	VideoList  interface{}
//	VideoNext  interface{}
//	Password   string
//	ApiData    *simplejson.Json
//	ApiErrCode int
//	ApiErrMsg  string
//
//	Streams map[string]*StreamStruct
//
//	CCode string
//	Vid   string
//	Utid  interface{}
//	Ua    string
//
//	PassProtected bool
//}
//
//type StreamStruct struct {
//	Id           string
//	Container    string
//	VideoProfile string
//	Size         int
//	Pieces       []Piece
//	M3u8Url      string
//}
//
//type Piece struct {
//	Segs string
//}
//
//var StreamTypes map[string]*StreamStruct
//
//func InitStream() {
//	StreamTypes = map[string]*StreamStruct{
//		"hd3":      &StreamStruct{Id: "hd3", Container: "flv", VideoProfile: "1080p"},
//		"hd3v2":    &StreamStruct{Id: "hd3v2", Container: "flv", VideoProfile: "1080p"},
//		"mp4hd3":   &StreamStruct{Id: "hd3v2", Container: "mp4", VideoProfile: "1080p"},
//		"mp4hd3v2": &StreamStruct{Id: "hd3v2", Container: "mp4", VideoProfile: "1080p"},
//
//		"hd2":   &StreamStruct{Id: "hd3", Container: "flv", VideoProfile: "超清"},
//		"hd2v2": &StreamStruct{Id: "hd3v2", Container: "flv", VideoProfile: "超清"},
//
//		// todo: 完善
//	}
//}

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

	if len(yk.Data.Stream) == 0 {
		log.Fatal("用户异常，请重新登录对应网站账号")
	}

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

	fmt.Printf("%+v \n", yk)
	urls, size, quality := yk.genData()
	format := common.FormatData{
		URLs:    urls,
		Size:    size,
		Quality: quality,
	}
	data = common.VideoData{
		Site:    "优酷 youku.com",
		Title:   title,
		Type:    "video",
		Formats: []common.FormatData{format},
	}

	return
}
