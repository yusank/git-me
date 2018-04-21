package common

type URLData struct {
	URL  string
	Size int64
	Ext  string
}

// VideoData data struct of video info
type VideoData struct {
	Site      string
	Title     string
	OutputDir string
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once
	Formats []FormatData
	Type    string
}


// FormatData data struct of every format
type FormatData struct {
	// [URLData: {URL, Size, Ext}, ...]
	// Some video files have multiple fragments
	// and support for downloading multiple image files at once
	URLs    []URLData
	Quality string
	Size    int64 // total size of all urls
}
