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

package extractors

import (
	"fmt"
	"log"
	"net/url"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/extractors/bilibili"
	"github.com/yusank/git-me/extractors/general"
	"github.com/yusank/git-me/extractors/iqiyi"
	"github.com/yusank/git-me/extractors/xiami"
	"github.com/yusank/git-me/extractors/youku"
	"github.com/yusank/git-me/extractors/youtube"
	"github.com/yusank/git-me/utils"
)

type CommonDownLoad func(url, outputDir string)

var (
	TransferMap = make(map[string]interface{})
)

func BeforeRun() {
	TransferMap["youku"] = youku.BasicInfo{}
	TransferMap["youtube"] = youtube.BasicInfo{}
	TransferMap["xiami"] = xiami.BasicInfo{}
	TransferMap["general"] = general.BasicInfo{}
	TransferMap["bilibili"] = bilibili.BasicInfo{}
	TransferMap["iqiyi"] = iqiyi.BasicInfo{}
}

func Foo(uri, output string, implement interface{}) {
	param := map[string]interface{}{
		"url":    uri,
		"output": output,
	}

	upload := common.UploadInfo{
		URL:      uri,
		Status:   common.TaskStatusFinish,
		Schedule: 100.0,
	}
	if err := common.DownloadByUrl(implement.(common.VideoExtractor), param); err != nil {
		upload.Status = common.TaskStatusFail
		upload.Schedule = 0.0
	}

	common.ProcessChan <- upload
}

func MatchUrl(videoURL, outputPath string) {
	var (
		domain     string
		downloader interface{}
		found      bool
	)

	bilibiliShortLink := utils.MatchOneOf(videoURL, `^(av|ep)\d+`)
	if bilibiliShortLink != nil {
		bilibiliURL := map[string]string{
			"av": "https://www.bilibili.com/video/",
			"ep": "https://www.bilibili.com/bangumi/play/",
		}

		domain = "bilibili"
		videoURL = bilibiliURL[bilibiliShortLink[1]] + videoURL
	} else {
		u, err := url.ParseRequestURI(videoURL)
		if err != nil {
			log.Fatal(err)
		}
		domain = utils.Domain(u.Host)
	}

	downloader, found = TransferMap[domain]
	if !found {
		fmt.Println("I am very sorry.I can't parese this kind of url yet. but I still try to download it.")
		downloader = TransferMap["general"]

	}

	Foo(videoURL, outputPath, downloader)
}
