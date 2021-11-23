package smbios

import (
	"fmt"

	gosmbios "github.com/digitalocean/go-smbios/smbios"
	"github.com/moxspec/moxspec/util"
)

// Processor represents a processor spec
//
// NOTE:
// this package ignores 0x08 Processor ID field
// bacause DMI v3.1.1 does NOT support extra features id.
type Processor struct {
	SocketDesignation string
	Manufacturer      string
	Version           string
	Voltage           float32
	ExternalClock     uint16
	MaxSpeed          uint16 // This field identifies a capability for the system, not the processor itself.
	CurrentSpeed      uint16 // This field identifies the processor's speed at system boot; the processor may support more than one speed.
	Status            []string
	L1CacheHandle     uint16
	L1Cache           *Cache
	L2CacheHandle     uint16
	L2Cache           *Cache
	L3CacheHandle     uint16
	L3Cache           *Cache
	SerialNumber      string
	AssetTag          string
	PartNumber        string
	CoreCount         uint16
	CoreEnabled       uint16
	ThreadCount       uint16
}

func parseProcessor(s *gosmbios.Structure) (*Processor, error) {
	if s == nil {
		return nil, fmt.Errorf("nil given")
	}

	p := new(Processor)

	p.SocketDesignation = getStringsSet(s, 0x04)
	p.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x07))
	p.Version = util.ShortenProcName(getStringsSet(s, 0x10))
	p.Voltage = parseProcessorVoltage(getByte(s, 0x11))
	p.ExternalClock = getWord(s, 0x12)
	p.MaxSpeed = getWord(s, 0x14)
	p.CurrentSpeed = getWord(s, 0x16)
	p.Status = parseProcessorStatus(getByte(s, 0x18))
	p.L1CacheHandle = uint16(getWord(s, 0x1A))
	p.L2CacheHandle = uint16(getWord(s, 0x1C))
	p.L3CacheHandle = uint16(getWord(s, 0x1E))
	p.SerialNumber = getStringsSet(s, 0x20)
	p.AssetTag = getStringsSet(s, 0x21)
	p.PartNumber = getStringsSet(s, 0x22)
	p.CoreCount = parseProcessorCount(getByte(s, 0x23), getWord(s, 0x2A))
	p.CoreEnabled = parseProcessorCount(getByte(s, 0x24), getWord(s, 0x2C))
	p.ThreadCount = parseProcessorCount(getByte(s, 0x25), getWord(s, 0x2E))

	log.Debugf("%+v", p)

	return p, nil
}

func parseProcessorStatus(b uint8) []string {
	list := []string{}

	if b&(1<<6) == 0 {
		list = append(list, "CPU Socket Unpopulated")
		return list
	}

	list = append(list, "CPU Socket Populated")

	stats := []string{
		"Unknown",
		"CPU Enabled",
		"CPU Disabled by User through BIOS Setup",
		"CPU Disabled By BIOS (POST Error)",
		"CPU is Idle, waiting to be enabled.",
		"Reserved",
		"Reserved",
		"Other",
	}

	statPtr := b & 0x07 // 0x07 = 0111
	if statPtr > uint8(len(stats)-1) {
		return list
	}

	return append(list, stats[statPtr])
}

func parseProcessorCount(cc uint8, cce uint16) uint16 {
	if cc == 0xFF && cce > 0xFF {
		return cce
	}
	return uint16(cc)
}

func parseProcessorVoltage(b uint8) float32 {
	// If bit 7 is set to 1, the remaining seven bits of the field are
	// set to contain the processor’s current voltage times 10.
	if b&(1<<7) == 0 {
		if b&1 != 0 {
			return 5.0
		}
		if b&(1<<1) != 0 {
			return 3.3
		}
		if b&(1<<2) != 0 {
			return 2.9
		}
		return 0.0
	}

	return float32((b & 0x7F)) / 10.0
}
