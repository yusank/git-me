package utils

import (
	"runtime"
	"strings"
)

// FileName Converts a string to a valid filename
func FileName(name string) string {
	// FIXME(iawia002) file name can't have /
	name = strings.Replace(name, "/", " ", -1)
	name = strings.Replace(name, "|", "-", -1)
	name = strings.Replace(name, ": ", "：", -1)
	name = strings.Replace(name, ":", "：", -1)
	name = strings.Replace(name, "#", " ", -1)
	if runtime.GOOS == "windows" {
		winSymbols := []string{
			"\"", "?", "*", "\\", "<", ">",
		}
		for _, symbol := range winSymbols {
			name = strings.Replace(name, symbol, " ", -1)
		}
	}
	return name
}

