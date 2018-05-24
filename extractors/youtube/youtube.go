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

package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/utils"
)

type BasicInfo struct {
	Url string
}

func (yt BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

type args struct {
	Title  string `json:"title"`
	Stream string `json:"url_encoded_fmt_stream_map"`
}

type assets struct {
	JS string `json:"js"`
}

type youtubeData struct {
	Args   args   `json:"args"`
	Assets assets `json:"assets"`
}

func getSig(sig, js string) string {
	fmt.Printf("sig:%s, js:%s \n", sig, js)
	html := utils.GetDecodeHTML(fmt.Sprintf("https://www.youtube.com%s", js), nil)
	return decipherTokens(getSigTokens(string(html)), sig)
}

// Youtube download function
func (yt BasicInfo) Download(url string) (vid common.VideoData, err error) {
	result := new(common.VideoData)

	if !common.Playlist {
		youtubeDownload(url, result)
		vid = *result
		return
	}
	listID := utils.MatchOneOf(url, `(list|p)=([^/&]+)`)[2]
	if listID == "" {
		log.Println("Can't get list ID from URL")
		return
	}
	html := utils.GetDecodeHTML("https://www.youtube.com/playlist?list="+listID, nil)
	// "videoId":"OQxX8zgyzuM","thumbnail"
	videoIDs := utils.MatchAll(string(html), `"videoId":"([^,]+?)","thumbnail"`)
	for _, videoID := range videoIDs {
		u := fmt.Sprintf(
			"https://www.youtube.com/watch?v=%s&list=%s", videoID[1], listID,
		)

		youtubeDownload(u, result)
	}

	vid = *result

	return
}

func youtubeDownload(uri string, result *common.VideoData) {
	fmt.Println("[uri]", uri)
	vid := utils.MatchOneOf(
		uri,
		`watch\?v=([^/&]+)`,
		`youtu\.be/([^?/]+)`,
		`embed/([^/?]+)`,
		`v/([^/?]+)`,
	)
	if vid == nil {
		log.Println("Can't find vid")
		return
	}
	fmt.Println("[vid]", vid)
	videoURL := fmt.Sprintf(
		"https://www.youtube.com/watch?v=%s&gl=US&hl=en&has_verified=1&bpctr=9999999999",
		vid[1],
	)
	fmt.Println("[vidURL]", videoURL)
	html := string(utils.GetDecodeHTML(videoURL, nil))
	//fmt.Println("[html]",html)
	yp := utils.MatchOneOf(html, `;ytplayer\.config\s*=\s*({.+?});`)
	//fmt.Println("[yp]", yp)
	ytplayer := ""
	if len(yp) > 1 {
		ytplayer = yp[1]
	}
	var youtube youtubeData
	json.Unmarshal([]byte(ytplayer), &youtube)
	title := youtube.Args.Title
	streams := strings.Split(youtube.Args.Stream, ",")
	stream, _ := url.ParseQuery(streams[0]) // Best quality
	quality := stream.Get("quality")
	fmt.Println("[quality]", quality)
	e := utils.MatchOneOf(stream.Get("type"), `video/(\w+);`)
	ext := ""
	if len(e) > 1 {
		ext = e[1]
	}
	streamURL := stream.Get("url")
	var realURL string
	if strings.Contains(streamURL, "signature=") {
		// URL itself already has a signature parameter
		realURL = streamURL
	} else {
		// URL has no signature parameter
		sig := stream.Get("sig")
		if sig == "" {
			// Signature need decrypt
			sig = getSig(stream.Get("s"), youtube.Assets.JS)
		}
		realURL = fmt.Sprintf("%s&signature=%s", streamURL, sig)
	}
	size := utils.DownloadFileSize(realURL, uri)
	urlData := common.URLData{
		URL:  realURL,
		Size: size,
		Ext:  ext,
	}

	format := common.FormatData{
		URLs:    []common.URLData{urlData},
		Size:    urlData.Size,
		Quality: quality,
	}

	result.Site = "YouTube youtube.com"
	result.Title = utils.FileName(title)
	fmt.Println("[youtube] title:", result.Title)
	result.Formats = append(result.Formats, format)
	return
}
