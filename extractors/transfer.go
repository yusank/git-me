package extractors

import (
	"fmt"
	"git-me/common"
	"git-me/extractors/bilibili"
	"git-me/extractors/general"
	"git-me/extractors/iqiyi"
	"git-me/extractors/netease"
	"git-me/extractors/xiami"
	"git-me/extractors/youku"
	"git-me/extractors/youtube"
	"git-me/utils"
	"log"
	"net/url"
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
	TransferMap["general"] = general.BasicInfo{}
	TransferMap["bilibili"] = bilibili.BasicInfo{}
	TransferMap["iqiyi"] = iqiyi.BasicInfo{}
}

func Foo(uri, output string, implement interface{}) {
	param := map[string]interface{}{
		"url":    uri,
		"output": output,
	}

	upload := common.UploadInfo{
		URL:    uri,
		Status: common.TaskStatusFinish,
	}
	if err := common.DownloadByUrl(implement.(common.VideoExtractor), param); err != nil {
		upload.Status = common.TaskStatusFail
	}

	common.FinishChan <- upload
}

func MatchUrl(videoURL, outputPath string) {
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
			log.Fatal(err)
		}
		domain = utils.Domain(u.Host)
	}

	downloader, found = TransferMap[domain]
	if !found {
		fmt.Println("I am very sorry.I can't parese this kind of url yet. but I still try to download it.")
		downloader = TransferMap["general"]

	}

	Foo(videoURL, outputPath, downloader)
}
