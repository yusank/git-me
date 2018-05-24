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

package general

/*
	通用下载,用于尝试直接下载无法解析的 URL
*/

import (
	"strings"

	"git-me/common"
	"git-me/utils"
)

type BasicInfo struct{}

func (gn BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

func (gn BasicInfo) Download(url string) (vid common.VideoData, err error) {
	exts := strings.Split(url, ".")
	ext := ""
	if len(exts) > 1 {
		ext = exts[len(exts)-1]
	}

	urlData := common.URLData{
		URL:  url,
		Size: common.DefaultSize,
		Ext:  ext,
	}

	format := common.FormatData{URLs: []common.URLData{urlData}}
	vid.Site = ""
	vid.Title = utils.FileName(exts[0])
	vid.Formats = []common.FormatData{format}

	return
}
