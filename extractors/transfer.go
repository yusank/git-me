package extractors

import (
	"git-me/common"
	"git-me/extractors/netease"
	"git-me/extractors/xiami"
	"git-me/extractors/youku"
	"git-me/extractors/youtube"
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
}

func Foo(uri,output string, implement interface{}) {
	param := map[string]interface{}{
		"url": uri,
		"output":output,
	}
	common.DownloadByUrl(implement.(common.VideoExtractor), param)
}
