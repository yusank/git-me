package extractors

type CommonDownLoad func(url, outputDir string)

var (
	TransferMap = make(map[string]CommonDownLoad)
)

func BeforeRun() {
	TransferMap["163"] = DownloadByURL
}