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

func Foo(key string) {
	val := TransferMap[key]
	param := map[string]interface{}{
		"url": "http://www.xiami.com/mv/K6YvR7?spm=a1z1s.2943549.6862561.2.CoVfLo",
	}
	common.DownloadByUrl(val.(common.VideoExtractor), param)
}
