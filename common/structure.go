package common

type URLData struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

// VideoData data struct of video info
type VideoData struct {
	Site      string `json:"site"`
	Title     string `json:"title"`
	OutputDir string `json:"-"`
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once
	Formats []FormatData `json:"formats"`
	Type    string       `json:"type"`
}

// FormatData data struct of every format
type FormatData struct {
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once
	URLs    []URLData `json:"urls"`
	Quality string    `json:"quality"`
	Size    int64     `json:"size"` // total size of all urls
}
