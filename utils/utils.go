package utils

import (
	"runtime"
	"strings"
	"crypto/md5"
	"encoding/hex"
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


// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)[1]
	return domain
}

func StringMd5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))

	return hex.EncodeToString(hash.Sum(nil))
}
