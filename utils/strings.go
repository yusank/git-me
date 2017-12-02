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

// MatchSlice match multi regex in text
func MatchSlice(text string, patterns []string) []string {
	if len(patterns) == 0 {
		return nil
	} else if len(patterns) == 1 {
		return Match(text, patterns[0])
	} else {
		result := []string{}
		for _, v := range patterns {
			match := Match(text, v)
			result = append(result, match...)
		}
		return result
	}
}
