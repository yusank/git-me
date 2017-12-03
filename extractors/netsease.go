package extractors

import (
	"fmt"
	"sort"
	"strings"

	"git-me/common"
	"git-me/utils"
)

// NeteaseCloudMusicDownload -
func NeteaseCloudMusicDownload(url, outputDir string) error {
	rid := utils.Match(`\Wid=(.*)`, url)
	if len(rid) == 0 {
		rid = utils.Match(`/(\d+)/?`, url)
		fmt.Println(rid)
	} else {
		newRid := []string{}
		for _, v := range rid {
			s := strings.Split(v, "=")
			newRid = append(newRid, s[1])
		}
		rid = newRid
	}

	header := make(map[string]string)
	header["Referer"] = "http://music.163.com/"

	switch {
	case strings.Contains(url, "mv"):
		fmt.Println("it`s mv")
		reqUrl := fmt.Sprintf("http://music.163.com/api/mv/detail/?id=%s&ids=%s&csrf_token=", rid[0], rid)
		body, err := utils.GetContent(reqUrl, "GET", header)
		if err != nil {
			return err
		}
		j, err := utils.LoadJSON(body)
		if err != nil {
			return err
		}

		//for k, v := range j.Get("data").Get("brs").MustMap() {
		//	fmt.Println(k, ":", v)
		//}

		vinfo := j.Get("data").MustMap()
		NeteaseMvDownload(vinfo, outputDir, false)

	case strings.Contains(url, "album"):
		body, err := utils.GetContent(url, "GET", header)
		if err != nil {
			return err
		}

		j, err := utils.LoadJSON(body)
		if err != nil {
			return err
		}

		fmt.Printf("%+v \n", j)
	}

	return nil
}

// DownloadByURL -
func DownloadByURL(url string) {
	if strings.Contains(url, "163.fm") {
		fmt.Println(url)
		return
	}

	if strings.Contains(url, "music.163.com") {
		fmt.Println(url)
		if err := NeteaseCloudMusicDownload(url, ""); err != nil {
			fmt.Println(err)
		}
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

	urls := []string{}
	exc := ""
	size := 0
	if len(src) > 0 {
		urls = src
		exc = "mp4"
		size = 1
		// url_info func, return ext, size

	} else {
		urls = utils.Match(`["\\'](.+)-list.m3u8["\\']`, data)
		if len(urls) == 0 {
			urls = utils.Match(`["\\'](.+).m3u8["\\']`, data)
		}

		size = 2
		exc = "mp4"
	}

	// todo:Print url_info
	// download func
	fmt.Println(urls)
	fmt.Println(data, exc, size)
}

func NeteaseMvDownload(vinfo map[string]interface{}, outputDir string, infoOnly bool) {
	title := fmt.Sprintf("%s - %s", vinfo["name"], vinfo["artistName"])
	urlBest := vinfo["brs"].(map[string]interface{})
	keys := []string{}
	for k := range urlBest {
		keys = append(keys, k)
	}
	// sort and get the best quality mv
	sort.Strings(keys)

	NeteaseDownloadCommon(title, urlBest[keys[len(keys)-1]].(string), outputDir, infoOnly)
}

func NeteaseDownloadCommon(title string, urlBest string, outputDir string, infoOnly bool) {
	songType, ext, size, err := common.UrlInfo(urlBest, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !infoOnly {
		fmt.Println("info:", title, songType, ext, size)
		common.DownloadURL([]string{urlBest}, title, ext, outputDir, size, false, nil)
	}
}
