package util

import "fmt"

// FindMSB returns most significant bit position in given value
func FindMSB(b interface{}) (int, error) {
	var size int
	var val uint64

	switch b.(type) {
	case int8:
		size = 8
		val = uint64(b.(int8))
	case int16:
		size = 16
		val = uint64(b.(int16))
	case int32:
		size = 16
		val = uint64(b.(int32))
	case int64:
		size = 64
		val = uint64(b.(int64))
	case int:
		size = 64
		val = uint64(b.(int))
	case byte:
		size = 8
		val = uint64(b.(byte))
	case uint16:
		size = 16
		val = uint64(b.(uint16))
	case uint32:
		size = 16
		val = uint64(b.(uint32))
	case uint64:
		size = 64
		val = uint64(b.(uint64))
	case uint:
		size = 64
		val = uint64(b.(uint))
	default:
		return -1, fmt.Errorf("unsupported type")
	}

	var i int
	for i = size - 1; i >= 0; i-- {
		if val&(1<<uint(i)) != 0 {
			return int(i), nil
		}
	}
	return -1, fmt.Errorf("no bits")
}
