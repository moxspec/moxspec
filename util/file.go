package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FilterPrefixedFiles returns filtered regular files which has given prefix
func FilterPrefixedFiles(path, prefix string) []string {
	return FilterFiles(path, func(f os.FileInfo) bool {
		return (strings.HasPrefix(f.Name(), prefix) && f.Mode().IsRegular())
	})
}

// FilterPrefixedLinks returns filtered symlinks which has given prefix
func FilterPrefixedLinks(path, prefix string) []string {
	return FilterFiles(path, func(f os.FileInfo) bool {
		return (strings.HasPrefix(f.Name(), prefix) && (f.Mode()&os.ModeSymlink) != 0)
	})
}

// FilterPrefixedDirs returns filtered directories which has given prefix
func FilterPrefixedDirs(path, prefix string) []string {
	return FilterFiles(path, func(f os.FileInfo) bool {
		return (strings.HasPrefix(f.Name(), prefix) && f.IsDir())
	})
}

// FilterFiles returns filtered items by filter function
func FilterFiles(path string, filter func(os.FileInfo) bool) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	var list []string
	for _, file := range files {
		if filter(file) {
			list = append(list, filepath.Join(path, file.Name()))
		}
	}

	return list
}
