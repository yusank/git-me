package extractors

import (
	"fmt"
	"net/url"

	"git-me/extractors/general"
	"git-me/extractors/netease"
	"git-me/extractors/xiami"
	"git-me/extractors/youku"
	"git-me/extractors/youtube"

	"git-me/common"
	"git-me/extractors/bilibili"
	"git-me/utils"
)

type CommonDownLoad func(url, outputDir string)

var (
	TransferMap = make(map[string]interface{})
)

func BeforeRun() {
	TransferMap["163"] = netease.BasicInfo{}
	TransferMap["youku"] = youku.BasicInfo{}
	TransferMap["youtube"] = youtube.BasicInfo{}
	TransferMap["xiami"] = xiami.BasicInfo{}
	TransferMap["bilibili"] = bilibili.BasicInfo{}
	TransferMap["general"] = general.BasicInfo{}
}

func Foo(uri string, implement interface{}) (vi *common.VideoData, err error) {
	param := map[string]interface{}{
		"url": uri,
	}
	return common.DownloadByUrl(implement.(common.VideoExtractor), param)
}

func MatchUrl(videoURL string) (vi *common.VideoData, err error) {
	if videoURL == "" {
		return nil, fmt.Errorf("nil url")
	}

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
			return nil, err
		}
		domain = utils.Domain(u.Host)
		fmt.Println(domain)
	}

	downloader, found = TransferMap[domain]
	if !found {
		fmt.Println("I am very sorry.I can't parese this kind of url yet. but I still try to download it.")
		downloader = TransferMap["general"]

	}

	return Foo(videoURL, downloader)
}
