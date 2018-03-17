package common

type VideoInfo struct {
	Url           string
	Title         string
	Vid           string
	M3u8Url       string
	Streams       interface{}
	StreamsSorted []interface{}
	PassProtected bool
	DashStreams   interface{}
	CaptionTracks interface{}
	Out           bool
	Ua            string
	Referer       string
	Danmuku       string
}
