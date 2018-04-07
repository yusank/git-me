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

)