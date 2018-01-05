package netease

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
		newRid := []string{}
		for _, v := range rid {
			s := v[1:]
			newRid = append(newRid, s[:len(s)-1])
		}
		rid = newRid
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

		j, err := utils.LoadJSON(reqUrl, header)
		if err != nil {
			return err
		}

		vinfo := j.Get("data").MustMap()
		NeteaseMvDownload(vinfo, outputDir, false)

	case strings.Contains(url, "album"):
		j, err := utils.LoadJSON(url, header)
		if err != nil {
			return err
		}

		fmt.Printf("%+v \n", j)
	case strings.Contains(url, "program"):
		reqUrl := fmt.Sprintf("http://music.163.com/api/dj/program/detail/?id=%s&ids=%s&csrf_token=", rid[0], rid)
		common.FakeHeader["Referer"] = "http://music.163.com/"
		fmt.Println(common.FakeHeader)
		j, err := utils.LoadJSON(reqUrl, header)
		if err != nil {
			return err
		}

		songInfo, err := j.String()
		fmt.Println(songInfo)
		NeteaseSongDownload(nil, outputDir, "", false)

	}

	return nil
}

// DownloadByURL -
func DownloadByURL(url, outputDir string) {
	if strings.Contains(url, "163.fm") {
		fmt.Println(url)
		return
	}

	if strings.Contains(url, "music.163.com") {
		fmt.Println(url)
		if err := NeteaseCloudMusicDownload(url, outputDir); err != nil {
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
	ext := ""
	size := 0
	if len(src) > 0 {
		urls = src
		ext = "mp4"
		size = 1
		// url_info func, return ext, size

	} else {
		urls = utils.Match(`["\\'](.+)-list.m3u8["\\']`, data)
		if len(urls) == 0 {
			urls = utils.Match(`["\\'](.+).m3u8["\\']`, data)
		}

		size = 2
		ext = "mp4"
	}

	// todo:Print url_info
	common.DownloadURL(urls, title, ext, outputDir, size, false, nil)
	fmt.Println(urls)
	fmt.Println(data, ext, size)
}

func NeteaseSongDownload(song map[string]interface{}, outputDir, playListPrefix string, infoOnly bool) {
	title := fmt.Sprintf("%s%s. %s", playListPrefix, song["position"], song["name"])
	mp3Url := song["mp3Url"]
	if mp3Url == nil {
		return
	}
	nets := strings.Split(song["mp3Url"].(string), "/")
	songNet := ""
	if len(nets) > 2 && len(nets[2]) > 1 {
		songNet = nets[2][1:]
	}

	urlBest := ""
	if hm, found := song["hMusic"]; found && hm != nil {
		dfs := song["hMusic"].(map[string]interface{})
		urlBest = MakeUrl(songNet, dfs["dfsId"].(string))
	} else if mp, found := song["mp3Url"]; found {
		urlBest = mp.(string)
	} else if _, found := song["bMusic"]; found {
		dfs := song["bMusic"].(map[string]interface{})
		urlBest = MakeUrl(songNet, dfs["dfsId"].(string))
	} else {
		return
	}

	NeteaseDownloadCommon(title, urlBest, outputDir, infoOnly)
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
		common.DownloadURL([]string{urlBest}, []string{title}, ext, outputDir, size, false, nil)
	}
}

func MakeUrl(songNet, dfsId string) string {
	return fmt.Sprintf("http://%s/%s/%s.mp3", songNet, "", dfsId)
}
