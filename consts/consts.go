package consts

type VideoCommon struct {
	Url               string
	Title             string
	Vid               string
	M3u8Url           string
	Streams           map[string]interface{}
	StreamsSort       []interface{} // 排序了的 stream
	AudioLang         string
	PasswordProtected bool
	DashStreams       map[string]interface{}
	CaptionTracks     []string
	Out               bool
	UA                string
	Referer           string
	Danmuku           string
}

const (
	DefaultSize = 1024 * 64
	DefaultText = "default"
	// session
	SessionUserID = "userId"

	ErrCodeSuccess = iota
)
