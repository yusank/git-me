package common

type VideoExtractor interface {
	//DownloadByUrl(url string, params map[string]interface{}) error
	//DownLoadByVid(vid string, params map[string]interface{}) error
	//Extract(params map[string]interface{}) error
	Download(url string) (VideoData, error)
}

func DownloadByUrl(v VideoExtractor, params map[string]interface{}) error {
	vid, err := v.Download(params["url"].(string))
	if err != nil {
		return err
	}

	vid.OutputDir = params["output"].(string)

	vid.Download(params["url"].(string))
	return nil
}
