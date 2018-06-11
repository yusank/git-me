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

package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/yusank/git-me/consts"
	"github.com/yusank/git-me/extractors"
	"github.com/yusank/git-me/models"
	"github.com/yusank/git-me/utils"
)

type DownloadController struct {
	BasicController
}

type DownloadInfo struct {
	URL  string `json:"url" valid:"Required"`
	Type int    `json:"type" valid:"Required"`
}

type CheckResp struct {
	IsUseful bool `json:"isUseful"`
}

func (dc *DownloadController) General() {
	var (
		di      DownloadInfo
		resp    *http.Response
		ext, ip string
		finish  bool
		now     = time.Now().Unix()
	)

	di.URL = dc.GetString("url")
	di.Type, _ = dc.GetInt("type")

	//di.URL = "https://www.bilibili.com/video/av23606332/?spm_id_from=333.334.bili_game.4"
	uid := dc.GetString("userId")
	if uid == "" {
		ip = dc.GetIp()
		if ip == "" {
			dc.OnCustomError(consts.ErrIpCannotFound)
			return
		}

		isBlock := models.IsBlocked(ip)
		if isBlock {
			dc.OnCustomError(consts.ErrIpHasBlocked)
			return
		}
	}

	vid, err := extractors.MatchUrl(di.URL)
	if err != nil {
		dc.OnCustomError(consts.ErrInvalidUrl)
		return
	}

	if di.Type == consts.DownloadTypeCheck {
		var resp CheckResp
		resp.IsUseful = len(vid.Formats) > 0 && len(vid.Formats["default"].URLs) > 0
		dc.JSON(resp)
		return
	}

	newFileName := fmt.Sprintf("temp-%d-%s.down", now, ip)
	file, _ := os.Create(newFileName)
	header := map[string]string{
		"Referer": di.URL,
	}

	go func() {
		writer := io.MultiWriter(file)
		if _, found := vid.Formats["default"]; found {
			for _, v := range vid.Formats["default"].URLs {
				ext = v.Ext
				fmt.Println(ext)
				resp, err = utils.HttpGet(v.URL, header)
				if err != nil {
					fmt.Println(err)
					return
				}
				io.Copy(writer, resp.Body)
				finish = true
			}
		}
	}()

	dc.Ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")

	for {
		info, _ := file.Stat()
		if info.Size() >= 1024*1024*2 || finish {
			resp.Body.Close()
			file.Close()

			newF, err := os.Open(newFileName)
			fmt.Println(err)
			b := make([]byte, info.Size())
			_, err = newF.Read(b)
			newF.Close()
			fmt.Println(err)

			dc.Ctx.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
			dc.Ctx.ResponseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=temp-%d.%s", now, ext))
			dc.Ctx.ResponseWriter.Write(b)
			os.Remove(newFileName)
			return
		}
	}

	dc.JSON(vid)
}
