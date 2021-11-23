package smbios

import (
	"fmt"

	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

// Cache represents a cache spec
type Cache struct {
	SocketDesignation   string
	OperationalMode     string
	Enabled             bool
	Location            string
	Socketed            bool
	Level               uint8
	MaximumCacheSize    uint32
	InstalledSize       uint32
	SupportedSRAMType   []string
	CurrentSRAMType     string
	CacheSpeed          uint8
	ErrorCorrectionType string
	SystemCacheType     string
	Associativity       string
}

func parseCache(s *gosmbios.Structure) (*Cache, error) {
	if s == nil {
		return nil, fmt.Errorf("nil given")
	}

	c := new(Cache)

	c.SocketDesignation = getStringsSet(s, 0x04)
	c.OperationalMode, c.Enabled, c.Location, c.Socketed, c.Level =
		parseCacheConfiguration(getWord(s, 0x05))
	c.MaximumCacheSize = parseCacheSize(getWord(s, 0x07), getDWord(s, 0x13))
	c.InstalledSize = parseCacheSize(getWord(s, 0x09), getDWord(s, 0x17))
	c.SupportedSRAMType = parseCacheSupportedSRAMType(getWord(s, 0x0B))
	c.CurrentSRAMType = parseCacheCurrentSRAMType(getWord(s, 0x0D))
	c.CacheSpeed = getByte(s, 0x0F)
	c.ErrorCorrectionType = parseCacheErrorCorrectionType(getByte(s, 0x10))
	c.SystemCacheType = parseSystemCacheType(getByte(s, 0x11))
	c.Associativity = parseCacheAssociativity(getByte(s, 0x12))

	log.Debugf("%+v", c)

	return c, nil
}

func parseCacheConfiguration(b uint16) (opmode string, enabled bool, location string, socketed bool, level uint8) {
	switch (b >> 8) & 0x03 {
	case 0:
		opmode = "Write Through"
	case 1:
		opmode = "Write Back"
	case 2:
		opmode = "Varies with Memory Address"
	default:
		opmode = "Unknown"
	}

	enabled = (b&(1<<7) != 0)

	switch (b >> 5) & 0x03 {
	case 0:
		location = "Internal"
	case 1:
		location = "External"
	case 2:
		location = "Reserved"
	default:
		location = "Unknown"
	}

	socketed = (b&(1<<3) != 0)

	level = uint8((b & 0x07) + 1)

	return
}

func parseCacheSize(b uint16, b2 uint32) uint32 {
	// For Cache sizes greater than 2047 MB,
	// the Maximum Cache Size field is set to 0xFFFF and
	// the Maximum Cache Size 2 field is present.
	if b == 0xFFFF && b2 != 0 {
		return parseCacheSize2(b2)
	}

	// bit 15 indicates granularity
	var gran uint32 = 1000
	if b&(1<<15) != 0 {
		gran = 64000
	}

	size := uint32(b&0x7FFF) * gran
	return size
}

func parseCacheSize2(b uint32) uint32 {
	// bit 31 indicates granularity
	var gran uint32 = 1000
	if b&(1<<31) != 0 {
		gran = 64000
	}

	size := (b & 0x7FFFFFFF) * gran
	return size
}

func parseCacheSupportedSRAMType(b uint16) []string {
	types := []string{
		"Other",
		"Unknown",
		"Non-Burst",
		"Burst",
		"Pipeline Burst",
		"Synchronous",
		"Asynchronous",
	}

	res := []string{}

	var i uint8
	for i = 0; i < uint8(len(types)); i++ {
		if b&(1<<i) != 0 {
			res = append(res, types[i])
		}
	}
	return res
}

func parseCacheCurrentSRAMType(b uint16) string {
	res := parseCacheSupportedSRAMType(b)
	if len(res) == 0 {
		return ""
	}
	return res[0]
}

func parseCacheErrorCorrectionType(b uint8) string {
	ec := []string{
		"Other", // 0x01
		"Unknown",
		"None",
		"Parity",
		"Single-bit ECC",
		"Multi-bit ECC", // 0x06
	}

	if b >= 0x01 && b <= 0x06 {
		return ec[b-0x01]
	}

	log.Warnf("CacheErrorCorrectionType: unsupported value %x was given.", b)
	return ""
}

func parseSystemCacheType(b uint8) string {
	ct := []string{
		"Other", // 0x01
		"Unknown",
		"Instruction",
		"Data",
		"Unified", // 0x05
	}

	if b >= 0x01 && b <= 0x05 {
		return ct[b-0x01]
	}

	log.Warnf("SystemCacheType: unsupported value %x was given.", b)
	return ""
}

func parseCacheAssociativity(b uint8) string {
	as := []string{
		"Other", // 0x01
		"Unknown",
		"Direct Mapped",
		"2-way Set-Associative",
		"4-way Set-Associative",
		"Fully Associative",
		"8-way Set-Associative",
		"16-way Set-Associative",
		"12-way Set-Associative",
		"24-way Set-Associative",
		"32-way Set-Associative",
		"48-way Set-Associative",
		"64-way Set-Associative",
		"20-way Set-Associative", // 0x0E
	}

	if b >= 0x01 && b <= 0x0E {
		return as[b-0x01]
	}

	log.Warnf("CacheAssociativity: unsupported value %x was given.", b)
	return ""
}
