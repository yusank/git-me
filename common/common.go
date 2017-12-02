package common

import (
	"fmt"
	"net/http"
	httpurl "net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"git-me/utils"
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
	forse = false
)

func UrlInfo(url string, fake bool, header map[string]string) (songType, ext string, size int, err error) {
	fmt.Println("url:", url)
	response := &http.Response{}
	if fake {
		return
	}

	response, err = utils.Response(url, header)
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

func DownloadURL(urls []string, title, ext, outputDir string, size int, fake bool, header map[string]string) {
	if len(urls) == 0 {
		return
	}

	if len(urls) == 1 {
		url := urls[0]
		fmt.Println("start downloading...", url)
		outPath := path.Join(outputDir, title)
		outPath = outPath + "." + ext
		if err := URLSave(url, outPath, "", fake, header); err != nil {
			fmt.Println(err)
		}

	}
}

func URLSave(url, path, refer string, fake bool, header map[string]string) error {
	var (
		err     error
		file    *os.File
		timeOut int
		resp    *http.Response
	)
	tmpHeaders := make(map[string]string)
	if refer != "" {
		tmpHeaders["Referer"] = refer
	}

	fileSize, err := URLSize(url, false, header)
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
	if !forse {
		open_mode = "ab"
		_, err := os.Stat(tmpFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				file, err = os.Create(tmpFilePath)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			size, err := file.Stat()
			if err != nil {
				return err
			}

			fs := size.Size()
			received += fs
		}
	} else {
		open_mode = "web"
	}

	if received < fileSize {

		if fake {
			tmpHeaders = FakeHeader
		}

		if received != 0 {
			tmpHeaders["Range"] = "bytes=" + fmt.Sprint(received) + "-"
		}

		if refer != "" {
			tmpHeaders["Referer"] = refer
		}

		if timeOut != 0 {
			resp, err = utils.RequestWithRetry(url, header)
			if err != nil || resp == nil {
				timeOut++
			}
		} else {
			resp, err = utils.Response(url, header)
			if resp != nil {
				fmt.Println("may be i got it")
			}
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println(open_mode)
		fmt.Println(resp.Header.Get("content-range"))

	}

	return nil
}

func URLSize(url string, fake bool, header map[string]string) (int64, error) {
	if fake {
		header = FakeHeader
	}

	resp, err := utils.Response(url, header)
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
