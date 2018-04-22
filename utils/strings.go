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

// MatchAll return all matching results
func MatchAll(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	value := re.FindAllStringSubmatch(text, -1)
	return value
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

// MatchOneOf match one of the patterns
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// Min return min
func Min(a, b int) int {
	if a > b {
		return b
	}

	return a
}
