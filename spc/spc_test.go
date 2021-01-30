package spc

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func loadTestData(path string) ([]byte, error) {
	bytes := []byte{}

	data, err := ioutil.ReadFile(filepath.Join("testdata", path))
	if err != nil {
		return nil, err
	}

	replaced := strings.ReplaceAll(string(data), "\n", " ")
	splited := strings.Split(replaced, " ")

	for _, s := range splited {
		if s == "" {
			continue
		}

		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, byte(b))
	}

	return bytes, nil
}
