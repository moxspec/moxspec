package util

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// LoadBytes reads the file named by path and returns the contents as []byte
func LoadBytes(path string) ([]byte, error) {
	if !Exists(path) {
		return nil, fmt.Errorf("not found %s", path)
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// LoadString reads the file named by path and returns the contens as string
func LoadString(path string) (string, error) {
	bs, err := LoadBytes(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// LoadUint64 reads the file named by path and returns the contents as uint64
func LoadUint64(path string) (uint64, error) {
	str, err := LoadString(path)
	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// LoadUint32 reads the file named by path and returns the contents as uint32
func LoadUint32(path string) (uint32, error) {
	u, err := LoadUint64(path)
	return uint32(u), err
}

// LoadUint16 reads the file named by path and returns the contents as uint16
func LoadUint16(path string) (uint16, error) {
	u, err := LoadUint64(path)
	return uint16(u), err
}

// LoadByte reads the file named by path and returns the contents as byte(uint8)
func LoadByte(path string) (byte, error) {
	u, err := LoadUint64(path)
	return byte(u), err
}
