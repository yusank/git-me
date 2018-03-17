package common

type VideoExtractor interface {
	//DownloadByUrl(url string, params map[string]interface{}) error
	//DownLoadByVid(vid string, params map[string]interface{}) error
	Prepare(params map[string]interface{}) error
	//Extract(params map[string]interface{}) error
	Download(params map[string]interface{}) error
}

func DownloadByUrl(v VideoExtractor, params map[string]interface{}) error {
	if err := v.Prepare(params);err != nil {
		return err
	}

	return v.Download(params)
}