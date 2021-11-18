package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const loaderDefaultTimeout = time.Second * 5

// LoadBytesWithContext reads the file named by path and returns the contents as []byte
func LoadBytesWithContext(ctx context.Context, path string) ([]byte, error) {
	if !Exists(path) {
		return nil, fmt.Errorf("not found %s", path)
	}

	var (
		bs  []byte
		err error
	)

	done := make(chan bool)
	go func() {
		bs, err = ioutil.ReadFile(path)
		done <- true
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout to read %s", path)
	case <-done:
	}

	return bs, err
}

// LoadBytes reads the file named by path and returns the contents as []byte
func LoadBytes(path string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), loaderDefaultTimeout)
	defer cancel()

	return LoadBytesWithContext(ctx, path)
}

// LoadStringWithContext reads the file named by path and returns the contens as string
func LoadStringWithContext(ctx context.Context, path string) (string, error) {
	bs, err := LoadBytesWithContext(ctx, path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// LoadString reads the file named by path and returns the contens as string
func LoadString(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), loaderDefaultTimeout)
	defer cancel()

	return LoadStringWithContext(ctx, path)
}

// LoadUint64WithContext reads the file named by path and returns the contents as uint64
func LoadUint64WithContext(ctx context.Context, path string) (uint64, error) {
	str, err := LoadStringWithContext(ctx, path)
	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// LoadUint64 reads the file named by path and returns the contents as uint64
func LoadUint64(path string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), loaderDefaultTimeout)
	defer cancel()

	return LoadUint64WithContext(ctx, path)
}

// LoadUint32WithContext reads the file named by path and returns the contents as uint32
func LoadUint32WithContext(ctx context.Context, path string) (uint32, error) {
	u, err := LoadUint64WithContext(ctx, path)
	return uint32(u), err
}

// LoadUint32 reads the file named by path and returns the contents as uint32
func LoadUint32(path string) (uint32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), loaderDefaultTimeout)
	defer cancel()

	return LoadUint32WithContext(ctx, path)
}

// LoadUint16WithContext reads the file named by path and returns the contents as uint16
func LoadUint16WithContext(ctx context.Context, path string) (uint16, error) {
	u, err := LoadUint64WithContext(ctx, path)
	return uint16(u), err
}

// LoadUint16 reads the file named by path and returns the contents as uint16
func LoadUint16(path string) (uint16, error) {
	ctx, cancel := context.WithTimeout(context.Background(), loaderDefaultTimeout)
	defer cancel()

	return LoadUint16WithContext(ctx, path)
}

// LoadByteWithContext reads the file named by path and returns the contents as byte(uint8)
func LoadByteWithContext(ctx context.Context, path string) (byte, error) {
	u, err := LoadUint64WithContext(ctx, path)
	return byte(u), err
}

// LoadByte reads the file named by path and returns the contents as byte(uint8)
func LoadByte(path string) (byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), loaderDefaultTimeout)
	defer cancel()

	return LoadByteWithContext(ctx, path)
}
