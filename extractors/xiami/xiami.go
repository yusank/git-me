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

package xiami

import (
	"encoding/xml"
	"fmt"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/utils"

	//"github.com/beevik/etree"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BasicInfo struct {
	Url    string
	Title  string
	Name   string
	LrcUrl string
}

type StringResources struct {
	XMLName        xml.Name         `xml:"resources"`
	ResourceString []ResourceString `xml:"string"`
}

type ResourceString struct {
	XMLName    xml.Name `xml:"string"`
	StringName string   `xml:"name,attr"`
	InnerText  string   `xml:",innerxml"`
}

func (xm BasicInfo) ParseVideo(url string) (common.VideoData, error) {
	// albums

	// collections

	// single track

	// mv
	return downloadMv(url)
}

func downloadMv(url string) (data common.VideoData, err error) {
	page, err := utils.HttpGetByte(url, nil)
	if err != nil {
		return
	}
	//fmt.Println(string(page))

	title := "xiami.flv"
	match := utils.Match(`<title>(^<]+)`, string(page))
	if len(match) > 0 {
		title = match[0]
	}

	data.Title = title
	data.Type = "video"

	vid, uid := "", ""
	match = utils.Match(`vid:"(\d+)"`, string(page))
	if len(match) > 0 {
		vid = strings.Split(match[0], `"`)[1]
		fmt.Println(vid)
	}

	match = utils.Match(`uid:"(\d+)"`, string(page))
	if len(match) > 0 {
		uid = strings.Split(match[0], `"`)[1]
		fmt.Println(uid)
	}

	apiUrl := fmt.Sprintf("http://cloud.video.taobao.com/videoapi/info.php?vid=%s&uid=%s", vid, uid)
	_, err = utils.HttpGetByte(apiUrl, nil)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocument(apiUrl)
	if err != nil {
		return
	}

	str := doc.Find("video_url").Eq(-1).Text()
	end := doc.Find("length").Eq(-1).Text()
	str += fmt.Sprintf("/start_%d/end_%s/1.flv", 0, end)
	u := common.URLData{
		URL: str,
	}
	format := common.FormatData{
		URLs: []common.URLData{u},
	}
	data.Formats = map[string]common.FormatData{
		"default": format,
	}

	return
}
