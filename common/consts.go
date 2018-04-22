package common

import "errors"

var (
	ErrUrlIsEmpty = errors.New("url is empty")
)

var (
	// Debug debug mode
	Debug bool
	// Version show version
	Version bool
	// InfoOnly Information only mode
	InfoOnly bool
	// Playlist download playlist
	Playlist bool
	// Refer use specified Referrer
	Refer string
	// Cookie http cookies
	Cookie string
	// Proxy HTTP proxy
	Proxy string
	// Socks5Proxy SOCKS5 proxy
	Socks5Proxy string
	// Format select specified format to download
	Format string
	// OutputPath output file path
	OutputPath string
	// OutputName output file name
	OutputName string
	// ExtractedData print extracted data
	ExtractedData bool
	// The number of download thread
	ThreadNumber int
)
