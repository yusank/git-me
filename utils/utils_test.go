package utils

import (
	"strings"
	"testing"
)

func TestFileName(t *testing.T) {
	symbol := "/#|:"
	suspectStrings := []string{"a/b", "a#b", "a|b", "a:b", "a：b"}

	for _, s := range suspectStrings {
		result := FileName(s)
		if strings.ContainsAny(result, symbol) {
			t.Fail()
		}
	}
}
