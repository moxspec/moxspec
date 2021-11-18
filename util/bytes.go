package util

// BytesToUint32 converts []byte to uint32
func BytesToUint32(b []byte) uint32 {
	if len(b) > 4 {
		return 0
	}

	padding := make([]byte, 4-len(b))
	buf := append(b, padding...)

	var ret uint32
	ret = uint32(buf[3])<<24 | uint32(buf[2])<<16 | uint32(buf[1])<<8 | uint32(buf[0])
	return ret
}

// BytesToUint64 converts []byte to uint64
func BytesToUint64(b []byte) uint64 {
	if len(b) > 8 {
		return 0
	}

	padding := make([]byte, 8-len(b))
	buf := append(b, padding...)

	var ret uint64
	lower := BytesToUint32(buf[:4])
	upper := BytesToUint32(buf[4:])
	ret = (uint64(upper) << 32) | uint64(lower)
	return ret
}
