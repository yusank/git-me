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
	Title   string `json:"title"`
	Stream  string `json:"url_encoded_fmt_stream_map"`
	Stream2 string `json:"url_encoded_fmt_stream_map"`
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

func genSignedURL(streamURL string, stream url.Values, js string) string {
	var realURL, sig string
	if strings.Contains(streamURL, "signature=") {
		// URL itself already has a signature parameter
		realURL = streamURL
	} else {
		// URL has no signature parameter
		sig = stream.Get("sig")
		if sig == "" {
			// Signature need decrypt
			sig = getSig(stream.Get("s"), js)
		}
		realURL = fmt.Sprintf("%s&signature=%s", streamURL, sig)
	}
	return realURL
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
		log.Fatal("Can't get list ID from URL")
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
		log.Fatal("Can't find vid")
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
	format := extractVideoURLS(youtube, uri)
	streams := strings.Split(youtube.Args.Stream, ",")
	stream, _ := url.ParseQuery(streams[0]) // Best quality
	quality := stream.Get("quality")
	fmt.Println("[quality]", quality)

	result.Site = "YouTube youtube.com"
	result.Title = utils.FileName(title)
	fmt.Println("[youtube] title:", result.Title)
	result.Formats = format
	return
}

func extractVideoURLS(data youtubeData, referer string) map[string]common.FormatData {
	streams := strings.Split(data.Args.Stream, ",")
	if data.Args.Stream == "" {
		streams = strings.Split(data.Args.Stream2, ",")
	}
	var ext string
	var audio common.URLData
	format := map[string]common.FormatData{}

	bestQualityURL, _ := url.ParseQuery(streams[0])
	bestQualityItag := bestQualityURL.Get("itag")

	for _, s := range streams {
		stream, _ := url.ParseQuery(s)
		itag := stream.Get("itag")
		streamType := stream.Get("type")
		isAudio := strings.HasPrefix(streamType, "audio/mp4")

		if !isAudio {
			continue
		}

		quality := stream.Get("quality_label")
		if quality == "" {
			quality = stream.Get("quality") // for url_encoded_fmt_stream_map
		}
		if quality != "" {
			quality = fmt.Sprintf("%s %s", quality, streamType)
		} else {
			quality = streamType
		}
		if isAudio {
			// audio file use m4a extension
			ext = "m4a"
		} else {
			ext = utils.MatchOneOf(streamType, `(\w+)/(\w+);`)[2]
		}
		realURL := genSignedURL(stream.Get("url"), stream, data.Assets.JS)
		if !strings.Contains(realURL, "ratebypass=yes") {
			realURL += "&ratebypass=yes"
		}
		size := utils.DownloadFileSize(realURL, referer)
		urlData := common.URLData{
			URL:  realURL,
			Size: size,
			Ext:  ext,
		}
		if ext == "m4a" {
			// Audio data for merging with video
			audio = urlData
		}
		format[itag] = common.FormatData{
			URLs:    []common.URLData{urlData},
			Size:    size,
			Quality: quality,
		}
	}

	format["default"] = format[bestQualityItag]
	delete(format, bestQualityItag)

	// `url_encoded_fmt_stream_map`
	if data.Args.Stream == "" {
		return format
	}

	// Unlike `url_encoded_fmt_stream_map`, all videos in `adaptive_fmts` have no sound,
	// we need download video and audio both and then merge them.
	// Another problem is that even if we add `ratebypass=yes`, the download speed still slow sometimes.

	// All videos here have no sound and need to be added separately
	for itag, f := range format {
		if strings.Contains(f.Quality, "video/") {
			f.Size += audio.Size
			f.URLs = append(f.URLs, audio)
			format[itag] = f
		}
	}
	return format
}
