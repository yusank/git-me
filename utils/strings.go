package utils

import (
	"fmt"
	"regexp"
)

// Match match pattern in text
func Match(pattern, text string) []string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return reg.FindAllString(text, -1)
}
