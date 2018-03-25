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
	URLs    []URLData
	Size    int64
	Type    string
	Quality string
}
