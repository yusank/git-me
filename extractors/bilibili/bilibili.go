package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"git-me/common"
	"git-me/utils"
)

const (
	bilibiliAPI        = "https://interface.bilibili.com/v2/playurl?"
	bilibiliBangumiAPI = "https://bangumi.bilibili.com/player/web_api/v2/playurl?"
	bilibiliTokenAPI   = "https://api.bilibili.com/x/player/playurl/token?"
)

const (
	// BiliBili blocks keys from time to time.
	// You can extract from the Android client or bilibiliPlayer.min.js
	appKey = "84956560bc028eb7"
	secKey = "94aba54af9065f71de72f5508f1cd42e"
)

const referer = "https://www.bilibili.com"

type BasicInfo struct {
	URL string
}

func genAPI(aid, cid string, bangumi bool, seasonType string) string {
	var (
		baseAPIURL string
		params     string
	)
	utoken := ""
	if common.Cookie != "" {
		utoken = utils.GetRequestStr(
			fmt.Sprintf("%said=%s&cid=%s", bilibiliTokenAPI, aid, cid),
			referer,
		)
		var t token
		json.Unmarshal([]byte(utoken), &t)
		if t.Code != 0 {
			log.Println(common.Cookie)
			log.Fatal("Cookie error: ", t.Message)
		}
		utoken = t.Data.Token
	}
	if bangumi {
		// The parameters need to be sorted by name
		// qn=0 flag makes the CDN address different every time
		// quality=116(1080P 60) is the highest quality so far
		params = fmt.Sprintf(
			"appkey=%s&cid=%s&module=bangumi&otype=json&qn=116&quality=116&season_type=%s&type=",
			appKey, cid, seasonType,
		)
		baseAPIURL = bilibiliBangumiAPI
	} else {
		params = fmt.Sprintf(
			"appkey=%s&cid=%s&otype=json&qn=116&quality=116&type=",
			appKey, cid,
		)
		baseAPIURL = bilibiliAPI
	}
	// bangumi utoken also need to put in params to sign, but the ordinary video doesn't need
	api := fmt.Sprintf(
		"%s%s&sign=%s", baseAPIURL, params, utils.StringMd5(params+secKey),
	)
	if !bangumi && utoken != "" {
		api = fmt.Sprintf("%s&utoken=%s", api, utoken)
	}
	return api
}

func genURL(durl []dURLData) ([]common.URLData, int64) {
	var (
		urls []common.URLData
		size int64
	)
	for _, data := range durl {
		size += data.Size
		url := common.URLData{
			URL:  data.URL,
			Size: data.Size,
			Ext:  "flv",
		}
		urls = append(urls, url)
	}
	return urls, size
}

type bilibiliOptions struct {
	Bangumi  bool
	Subtitle string
	Aid      string
	Cid      string
	HTML     string
}

func getMultiPageData(html string) (multiPage, error) {
	var data multiPage
	multiPageDataString := utils.MatchOneOf(
		html, `window.__INITIAL_STATE__=(.+?);\(function`,
	)
	if multiPageDataString == nil {
		return data, errors.New("This page has no playlist")
	}
	json.Unmarshal([]byte(multiPageDataString[1]), &data)
	return data, nil
}

func (bl BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

// Download bilibili main download function
func (bl BasicInfo) Download(url string) (vid common.VideoData, err error) {
	result := new(common.VideoData)
	var options bilibiliOptions
	if strings.Contains(url, "bangumi") {
		options.Bangumi = true
	}
	html := utils.GetRequestStr(url, referer)
	if !common.Playlist {
		options.HTML = html
		data, err := getMultiPageData(html)
		if err == nil && !options.Bangumi {
			// handle URL that has a playlist, mainly for unified titles
			// <h1> tag does not include subtitles
			// bangumi doesn't need this
			pageString := utils.MatchOneOf(url, `\?p=(\d+)`)
			var p int
			if pageString == nil {
				// https://www.bilibili.com/video/av20827366/
				p = 1
			} else {
				// https://www.bilibili.com/video/av20827366/?p=2
				p, _ = strconv.Atoi(pageString[1])
			}
			page := data.VideoData.Pages[p-1]
			options.Aid = data.Aid
			options.Cid = strconv.Itoa(page.Cid)
			// "part":"" or "part":"Untitled"
			if page.Part == "Untitled" {
				options.Subtitle = ""
			} else {
				options.Subtitle = page.Part
			}
		}
		bilibiliDownload(url, options, result)
		vid = *result
		return
	}
	if options.Bangumi {
		dataString := utils.MatchOneOf(html, `window.__INITIAL_STATE__=(.+?);`)[1]
		var data bangumiData
		json.Unmarshal([]byte(dataString), &data)
		for _, u := range data.EpList {
			bilibiliDownload(
				fmt.Sprintf("https://www.bilibili.com/bangumi/play/ep%d", u.EpID), options, result,
			)
		}
	} else {
		data, err := getMultiPageData(html)
		if err != nil {
			// this page has no playlist
			options.HTML = html
			bilibiliDownload(url, options, result)
			return
		}
		// https://www.bilibili.com/video/av20827366/?p=1
		for _, u := range data.VideoData.Pages {
			options.Aid = data.Aid
			options.Cid = strconv.Itoa(u.Cid)
			options.Subtitle = u.Part
			bilibiliDownload(url, options, result)
		}
	}

	vid = *result
	return
}

func bilibiliDownload(url string, options bilibiliOptions, result *common.VideoData) *common.VideoData {
	var (
		aid, cid, html string
	)
	if options.HTML != "" {
		// reuse html string, but this can't be reused in case of playlist
		html = options.HTML
	} else {
		html = utils.GetRequestStr(url, referer)
	}
	if options.Aid != "" && options.Cid != "" {
		aid = options.Aid
		cid = options.Cid
	} else {
		if options.Bangumi {
			cid = utils.MatchOneOf(html, `"cid":(\d+)`)[1]
			aid = utils.MatchOneOf(html, `"aid":(\d+)`)[1]
		} else {
			cid = utils.MatchOneOf(html, `cid=(\d+)`)[1]
			aid = utils.MatchOneOf(url, `av(\d+)`)[1]
		}
	}
	var seasonType string
	if options.Bangumi {
		seasonType = utils.MatchOneOf(html, `"season_type":(\d+)`)[1]
	}
	api := genAPI(aid, cid, options.Bangumi, seasonType)
	apiData := utils.GetRequestStr(api, referer)
	var dataDict bilibiliData
	json.Unmarshal([]byte(apiData), &dataDict)

	// get the title
	doc := utils.GetDoc(html)
	title := utils.Title(doc)
	if options.Subtitle != "" {
		title = fmt.Sprintf("%s %s", title, options.Subtitle)
	}

	urls, size := genURL(dataDict.DURL)
	format := common.FormatData{
		URLs:    urls,
		Size:    size,
		Quality: quality[dataDict.Quality],
	}
	result.Site = "哔哩哔哩 bilibili.com"
	result.Title = title
	result.Type = "video"
	result.Formats = append(result.Formats, format)
	return result
}
