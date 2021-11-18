package util

import (
	"fmt"
)

// ScanPathList returns path which was found by firtst from the list
func ScanPathList(list []string) (string, error) {
	for _, l := range list {
		if Exists(l) {
			return l, nil
		}
	}
	return "", fmt.Errorf("not found")
}
