package util

import (
	"context"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Exists returns if the file named by path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Executable returns if the file named by path is executable
func Executable(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}

	// 0x0555 == 0101 0101 0101
	if f.Mode()|0x0555 == 0 {
		return false
	}

	return true
}

// Exec runs specified command and returns output as string
func Exec(c string, arg ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, c, arg...)

	out, err := cmd.Output()

	if err != nil {
		return string(out), err
	}

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return string(out), nil
}

// DumpBinary returns the string that represents the value as binary array
func DumpBinary(in interface{}) string {
	var mlen byte
	var tgt uint64
	switch in.(type) {
	case byte:
		mlen = 8
		tgt = uint64(in.(byte))
	case uint16:
		mlen = 16
		tgt = uint64(in.(uint16))
	case uint32:
		mlen = 32
		tgt = uint64(in.(uint32))
	case uint64:
		mlen = 64
		tgt = in.(uint64)
	default:
		return ""
	}

	amlen := mlen + (mlen/4 - 1) // actual maximum length
	res := make([]string, amlen, amlen)
	var i byte
	var j = amlen - 1
	for i < mlen {
		r := "1"
		if tgt&(1<<i) == 0 {
			r = "0"
		}
		res[j] = r
		if j != 0 && j != amlen-1 && (i+1)%4 == 0 {
			j--
			res[j] = " "
		}
		j--
		i++
	}

	return strings.Join(res, "")
}

// BlkLabelAscSorter is the sorter for sort.Slice() to sort block device labels by asc
func BlkLabelAscSorter(a, b string) bool {
	if len(a) == len(b) {
		return (a < b)
	}

	return len(a) < len(b)
}

// IPv4MaskSize returns mask size of given string representation
func IPv4MaskSize(mask string) int {
	i := net.ParseIP(mask)
	if i == nil {
		return 0
	}
	i4 := i.To4()

	m := net.IPv4Mask(i4[0], i4[1], i4[2], i4[3])
	if m == nil {
		return 0
	}
	size, _ := m.Size()
	return size
}
