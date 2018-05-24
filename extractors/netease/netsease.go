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

/*
	放弃
*/
package netease

import (
	"fmt"
	"strings"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/utils"
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
