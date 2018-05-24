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

package common

import (
	"fmt"
	"io"
	"net/http"
	httpurl "net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"sync"

	"github.com/yusank/git-me/utils"
)

/*
	此文件主要放置通用配置、通用函数以及全局变量
*/
var (
	// FakeHeader for when web sites checke request header
	FakeHeader = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Charset":  "UTF-8,*;q=0.5",
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": "en-US,en;q=0.8",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:51.0) Gecko/20100101 Firefox/51.0",
	}

	// MediaTypeMap -
	MediaTypeMap = map[string]string{
		"video/3gpp":      "3gp",
		"video/f4v":       "flv",
		"video/mp4":       "mp4",
		"video/MP2T":      "ts",
		"video/quicktime": "mov",
		"video/webm":      "webm",
		"video/x-flv":     "flv",
		"video/x-ms-asf":  "asf",
		"audio/mp4":       "mp4",
		"audio/mpeg":      "mp3",
		"audio/wav":       "wav",
		"audio/x-wav":     "wav",
		"audio/wave":      "wav",
		"image/jpeg":      "jpg",
		"image/png":       "png",
		"image/gif":       "gif",
		"application/pdf": "pdf",
	}

	// SitesMap -
	SitesMap = map[string]string{
		"163":              "netease",
		"56":               "w56",
		"acfun":            "acfun",
		"archive":          "archive",
		"baidu":            "baidu",
		"bandcamp":         "bandcamp",
		"baomihua":         "baomihua",
		"bigthink":         "bigthink",
		"bilibili":         "bilibili",
		"cctv":             "cntv",
		"cntv":             "cntv",
		"cbs":              "cbs",
		"coub":             "coub",
		"dailymotion":      "dailymotion",
		"dilidili":         "dilidili",
		"douban":           "douban",
		"douyu":            "douyutv",
		"ehow":             "ehow",
		"facebook":         "facebook",
		"fantasy":          "fantasy",
		"fc2":              "fc2video",
		"flickr":           "flickr",
		"freesound":        "freesound",
		"fun":              "funshion",
		"google":           "google",
		"giphy":            "giphy",
		"heavy-music":      "heavymusic",
		"huaban":           "huaban",
		"huomao":           "huomaotv",
		"iask":             "sina",
		"icourses":         "icourses",
		"ifeng":            "ifeng",
		"imgur":            "imgur",
		"in":               "alive",
		"infoq":            "infoq",
		"instagram":        "instagram",
		"interest":         "interest",
		"iqilu":            "iqilu",
		"iqiyi":            "iqiyi",
		"isuntv":           "suntv",
		"joy":              "joy",
		"kankanews":        "bilibili",
		"khanacademy":      "khan",
		"ku6":              "ku6",
		"kugou":            "kugou",
		"kuwo":             "kuwo",
		"le":               "le",
		"letv":             "le",
		"lizhi":            "lizhi",
		"magisto":          "magisto",
		"metacafe":         "metacafe",
		"mgtv":             "mgtv",
		"miomio":           "miomio",
		"mixcloud":         "mixcloud",
		"mtv81":            "mtv81",
		"musicplayon":      "musicplayon",
		"naver":            "naver",
		"7gogo":            "nanagogo",
		"nicovideo":        "nicovideo",
		"panda":            "panda",
		"pinterest":        "pinterest",
		"pixnet":           "pixnet",
		"pptv":             "pptv",
		"qingting":         "qingting",
		"qq":               "qq",
		"quanmin":          "quanmin",
		"showroom-live":    "showroom",
		"sina":             "sina",
		"smgbb":            "bilibili",
		"sohu":             "sohu",
		"soundcloud":       "soundcloud",
		"ted":              "ted",
		"theplatform":      "theplatform",
		"tucao":            "tucao",
		"tudou":            "tudou",
		"tumblr":           "tumblr",
		"twimg":            "twitter",
		"twitter":          "twitter",
		"ucas":             "ucas",
		"videomega":        "videomega",
		"vidto":            "vidto",
		"vimeo":            "vimeo",
		"wanmen":           "wanmen",
		"weibo":            "miaopai",
		"veoh":             "veoh",
		"vine":             "vine",
		"vk":               "vk",
		"xiami":            "xiami",
		"xiaokaxiu":        "yixia",
		"xiaojiadianvideo": "fc2video",
		"ximalaya":         "ximalaya",
		"yinyuetai":        "yinyuetai",
		"miaopai":          "yixia",
		"yizhibo":          "yizhibo",
		"youku":            "youku",
		"iwara":            "iwara",
		"youtu":            "youtube",
		"youtube":          "youtube",
		"zhanqi":           "zhanqi",
		"365yg":            "toutiao",
	}
)

const (
	isForce = false
)

func UrlInfo(url string, fake bool, header map[string]string) (songType, ext string, size int, err error) {
	fmt.Println("url:", url)
	response := &http.Response{}
	if fake {
		return
	}

	response, err = utils.HttpGet(url, header)
	if err != nil {
		return
	}

	tp := response.Header.Get("content-type")
	if tp == "image/jpg; charset=UTF-8" || tp == "image/jpg" {
		tp = "audio/mpeg" // fix for netease
	}

	if ex, found := MediaTypeMap[tp]; found {
		ext = ex
	} else {
		tp = ""
		cd := response.Header.Get("content-disposition")
		if cd != "" {
			escape, err := httpurl.QueryUnescape(cd)
			if err != nil {
				fmt.Println(err)
			} else {
				m := utils.Match(`filename="?([^"]+)"?`, escape)
				if len(m) != 0 {
					ss := strings.Split(m[0], ".")
					if len(ss) > 1 {
						ext = ss[len(ss)-1]
					}
				}
			}
		}
	}

	if response.Header.Get("transfer-encoding") != "chunked" {
		size, _ = strconv.Atoi(response.Header.Get("content-length"))
	}

	return
}

func (vid VideoData) DownloadURL(refer string) error {
	if len(vid.Formats) < 1 {
		return fmt.Errorf("empty url")
	}
	if len(vid.Formats[0].URLs) < 1 {
		return fmt.Errorf("empty url")
	}
	var wg sync.WaitGroup

	for i, format := range vid.Formats {
		for _, url := range format.URLs {
			fmt.Println("start downloading...", url.URL)
			outPath := path.Join(vid.OutputDir, fmt.Sprintf("%s[%d].%s", vid.Title, i, url.Ext))
			if strings.Contains(refer, "mgtv") {
				// Too many threads cause mgtv to return HTTP 403 error
				URLSave(url, outPath, refer)
			} else {
				wg.Add(1)
				go func(url URLData, output, referer string) {
					URLSave(url, output, referer)
					wg.Done()
				}(url, outPath, refer)
			}
		}
		wg.Wait()
	}

	return nil
}

//todo:完善各类检查：1.文件已存在
func URLSave(url URLData, path, refer string) error {
	var (
		err     error
		file    *os.File
		timeOut int
		resp    *http.Response
	)

	tmpHeaders := make(map[string]string)

	if refer != "" {
		tmpHeaders["Referer"] = url.URL
	}

	fileSize, err := URLSize(url.URL, false, tmpHeaders)
	fmt.Println("size", fileSize)
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("not exist")
		} else {
			return err
		}
	}

	tmpFilePath := path + ".download"

	var received int64
	open_mode := ""
	if !isForce {
		open_mode = "ab"
		_, err := os.Stat(tmpFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				file, err = os.Create(tmpFilePath)
				if err != nil {
					fmt.Printf("create %s error: %v \n", tmpFilePath, err)
					return err
				}
				defer file.Close()
			} else {
				fmt.Println("get file stat error")
				return err
			}
		} else {
			file, err = os.OpenFile(tmpFilePath, os.O_RDWR|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Println("file open error")
				return err
			}
			defer file.Close()
			size, err := file.Stat()
			if err != nil {
				fmt.Println("get size got error")
				return err
			}

			fs := size.Size()
			received += fs
		}
	} else {
		open_mode = "web"
	}

	if received < fileSize {
		//if received != 0 {
		tmpHeaders["Range"] = "bytes=" + fmt.Sprint(received) + "-"
		//}

		if refer != "" {
			tmpHeaders["Referer"] = url.URL
		}

		if timeOut != 0 {
			resp, err = utils.RequestWithRetry(url.URL, tmpHeaders)
			if err != nil || resp == nil {
				timeOut++
			}
		} else {
			resp, err = utils.HttpGet(url.URL, tmpHeaders)
			if resp != nil {
				fmt.Println("may be i got it")
			}
			if err != nil {
				fmt.Println("1", err)
			}
		}

		fmt.Println(open_mode)
		var rangeLength int64
		leng := resp.Header.Get("content-range")
		if leng != "" {
			leng := leng[6:]
			lengStart := strings.Split(leng, "/")[0]
			lengStart = strings.Split(lengStart, "-")[0]

			lengEnd := strings.Split(leng, "/")[1]

			rangeStart, err := strconv.ParseInt(lengStart, 10, 64)
			if err != nil {
				return err
			}
			rangeEnd, err := strconv.ParseInt(lengEnd, 10, 64)
			if err != nil {
				return err
			}
			rangeLength = rangeEnd - rangeStart
		} else {
			leng = resp.Header.Get("content-length")

			rangeLength, err = strconv.ParseInt(leng, 10, 64)
			if err != nil {
				return err
			}
		}

		if fileSize != received+rangeLength {
			received = 0
			open_mode = "wb"
		}

		start := time.Now()
		for {
			buffer := make([]byte, 1024*256)
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			if n == 0 {
				if received == fileSize {
					break // finish!
				}

				tmpHeaders["Range"] = "bytes=" + fmt.Sprint(received) + "-"
				resp, err = utils.RequestWithRetry(url.URL, tmpHeaders)
				if err != nil {
					fmt.Println(err)
					// todo: 如何处理好？
				}
				continue
			}

			buffer = buffer[:n]
			wn, err := file.Write(buffer)
			if err != nil {
				fmt.Println(err)
				return err
				//todo: when actually use it should be continue but return
			}

			if wn != (n) {
				fmt.Println("write error:", wn, n)
			}
			received += int64(wn)

			fmt.Printf("正在下载: \r%.2f ...", (float64(received)/float64(fileSize))*100)
		}

		fmt.Println("用时：", time.Now().Sub(start))
		if err = os.Rename(tmpFilePath, path); err != nil {
			fmt.Println("rename fiald:", err)
		}

	}

	return nil
}

func URLSize(url string, fake bool, header map[string]string) (int64, error) {
	if fake {
		header = FakeHeader
	}

	resp, err := utils.HttpGet(url, header)
	if err != nil {
		return 0, err
	}

	s := resp.Header.Get("content-length")
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return size, err
}
