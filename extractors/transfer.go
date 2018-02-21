package extractors

import (
	"git-me/extractors/netease"
	"git-me/extractors/youku"
	"git-me/extractors/youtube"
)

type CommonDownLoad func(url, outputDir string)

var (
	TransferMap = make(map[string]CommonDownLoad)
)

func BeforeRun() {
	TransferMap["163"] = netease.DownloadByURL
	TransferMap["youku"] = youku.DownLoadByURL
	TransferMap["youtube"] = youtube.DownLoadByURL
}
