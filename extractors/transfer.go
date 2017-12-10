package extractors

import (
	"git-me/extractors/netease"
)

type CommonDownLoad func(url, outputDir string)

var (
	TransferMap = make(map[string]CommonDownLoad)
)

func BeforeRun() {
	TransferMap["163"] = netease.DownloadByURL
}
