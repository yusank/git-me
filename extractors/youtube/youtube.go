package youtube

type BasicInfo struct {
	Url string
}

func (yt BasicInfo) Prepare(params map[string]interface{}) error {
	return nil
}

func (yt BasicInfo) Download(param map[string]interface{}) error {
	return nil
}
