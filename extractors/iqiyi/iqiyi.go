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

package iqiyi

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/utils"
)

type vidl struct {
	M3utx      string `json:"m3utx"`
	Vd         int    `json:"vd"` // quality number
	ScreenSize string `json:"screenSize"`
}

type iqiyiData struct {
	Vidl []vidl `json:"vidl"`
}

type iqiyi struct {
	Code string    `json:"code"`
	Data iqiyiData `json:"data"`
}

type BasicInfo struct {
}

var iqiyiFormats = []int{
	18, // 1080p
	5,  // 1072p, 1080p
	17, // 720p
	4,  // 720p
	21, // 504p
	2,  // 480p, 504p
	1,  // 336p, 360p
	96, // 216p, 240p
}

const iqiyiReferer = "https://www.iqiyi.com"

func getIqiyiData(tvid, vid string) iqiyi {
	t := time.Now().Unix() * 1000
	src := "76f90cbd92f94a2e925d83e8ccd22cb7"
	key := "d5fb4bd9d50c4be6948c97edd7254b0e"
	sc := utils.StringMd5(strconv.FormatInt(t, 10) + key + vid)
	info := utils.GetRequestStr(
		fmt.Sprintf(
			"http://cache.m.iqiyi.com/jp/tmts/%s/%s/?t=%d&sc=%s&src=%s",
			tvid, vid, t, sc, src,
		),
		iqiyiReferer,
	)
	var data iqiyi
	json.Unmarshal([]byte(info[len("var tvInfoJs="):]), &data)
	return data
}

// Iqiyi download function
func (iq BasicInfo) Download(url string) (data common.VideoData, err error) {
	html := utils.GetRequestStr(url, iqiyiReferer)
	tvid := utils.MatchOneOf(
		url,
		`#curid=(.+)_`,
		`tvid=([^&]+)`,
	)
	if tvid == nil {
		tvid = utils.MatchOneOf(
			html,
			`data-player-tvid="([^"]+)"`,
			`param\['tvid'\]\s*=\s*"(.+?)"`,
		)
	}
	vid := utils.MatchOneOf(
		url,
		`#curid=.+_(.*)$`,
		`vid=([^&]+)`,
	)
	if vid == nil {
		vid = utils.MatchOneOf(
			html,
			`data-player-videoid="([^"]+)"`,
			`param\['vid'\]\s*=\s*"(.+?)"`,
		)
	}
	doc, err := utils.GetDoc(html)
	if err != nil {
		fmt.Println(err)
		return

	}
	title := strings.TrimSpace(doc.Find("h1 a").Text()) +
		strings.TrimSpace(doc.Find("h1 span").Text())
	if title == "" {
		title = doc.Find("title").Text()
	}
	videoDatas := getIqiyiData(tvid[1], vid[1])
	if videoDatas.Code != "A00000" {
		log.Fatal("Can't play this video")
	}
	var format []common.FormatData
	var urlData common.URLData
	var size, totalSize int64
	for _, video := range videoDatas.Data.Vidl {
		var urls []common.URLData
		totalSize = 0
		for _, ts := range utils.M3u8URLs(video.M3utx) {
			size, _ = strconv.ParseInt(
				utils.MatchOneOf(ts, `contentlength=(\d+)`)[1], 10, 64,
			)
			// http://dx.data.video.qiyi.com -> http://data.video.qiyi.com
			urlData = common.URLData{
				URL:  strings.Replace(ts, "dx.data.video.qiyi.com", "data.video.qiyi.com", 1),
				Size: size,
				Ext:  "ts",
			}
			totalSize += size
			urls = append(urls, urlData)
		}
		format = append(format, common.FormatData{
			URLs:    urls,
			Size:    totalSize,
			Quality: video.ScreenSize,
		})
	}
	// get best quality
	var videoData vidl
	for _, quality := range iqiyiFormats {
		for index, video := range videoDatas.Data.Vidl {
			if video.Vd == quality {
				videoData = videoDatas.Data.Vidl[index]
				break
			}
		}
		if videoData.M3utx != "" {
			break
		}
	}

	data = common.VideoData{
		Site:    "爱奇艺 iqiyi.com",
		Title:   title,
		Type:    "video",
		Formats: format,
	}

	return
}
