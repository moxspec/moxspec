package util

import "strings"

// HasPrefixIn returns true if given string has one of prefix
func HasPrefixIn(l string, ptns ...string) bool {
	if len(ptns) == 0 {
		return false
	}

	for _, p := range ptns {
		if strings.HasPrefix(l, p) {
			return true
		}
	}
	return false
}
