package smbios

import (
	"fmt"

	gosmbios "github.com/digitalocean/go-smbios/smbios"
	"github.com/moxspec/moxspec/util"
)

const (
	typeDetailPersistent = "Non-volatile"
	typeDetailUDIMM      = "Unbuffered (Unregistered)"
	typeDetailRDIMM      = "Registered (Buffered)"
	typeDetailLRDIMM     = "LRDIMM"
)

// MemoryDevice represents a memory module spec
type MemoryDevice struct {
	TotalWidth        uint16
	DataWidth         uint16
	Size              uint32 // GiB
	DeviceLocator     string
	BankLocator       string
	Manufacturer      string
	SerialNumber      string
	AssetTag          string
	PartNumber        string
	Speed             uint16
	ConfiguredSpeed   uint16
	MinVoltage        float32
	MaxVoltage        float32
	ConfiguredVoltage float32
	FormFactor        string
	Type              string
	TypeDetail        []string
}

// IsPersistent returns whether MemoryDevice is persistent memory.
func (m *MemoryDevice) IsPersistent() bool {
	for _, deteil := range m.TypeDetail {
		if deteil == typeDetailPersistent {
			return true
		}
	}

	return false
}

func getModuleType(typeDetail []string) string {
	const unknown = -1
	const unbuffered = 0
	const registered = 1
	const loadReduced = 2

	moduleTypeString := map[int]string{
		unknown:     "DIMM",
		unbuffered:  "UDIMM",
		registered:  "RDIMM",
		loadReduced: "LRDIMM",
	}

	moduleType := unknown
	for _, detail := range typeDetail {
		switch detail {
		case typeDetailUDIMM:
			if unbuffered > moduleType {
				moduleType = unbuffered
			}
		case typeDetailRDIMM:
			if registered > moduleType {
				moduleType = registered
			}
		case typeDetailLRDIMM:
			if loadReduced > moduleType {
				moduleType = loadReduced
			}
		}

	}

	return moduleTypeString[moduleType]
}

func parseMemoryDevice(s *gosmbios.Structure) (*MemoryDevice, error) {
	mem := new(MemoryDevice)

	if len(s.Strings) == 0 {
		return nil, fmt.Errorf("no data")
	}

	mem.DeviceLocator = s.Strings[0]
	mem.TotalWidth = getWord(s, 0x08) // 0xffff == unknown
	mem.DataWidth = getWord(s, 0x0A)
	mem.Size = parseMemorySize(getWord(s, 0x0C), getDWord(s, 0x1C))
	mem.DeviceLocator = getStringsSet(s, 0x10)
	mem.BankLocator = getStringsSet(s, 0x11)
	mem.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x17))
	mem.SerialNumber = getStringsSet(s, 0x18)
	mem.AssetTag = getStringsSet(s, 0x19)
	mem.PartNumber = getStringsSet(s, 0x1A)
	mem.Speed = getWord(s, 0x15)
	mem.ConfiguredSpeed = getWord(s, 0x20)
	mem.MinVoltage = float32(getWord(s, 0x22)) / 1000
	mem.MaxVoltage = float32(getWord(s, 0x24)) / 1000
	mem.ConfiguredVoltage = float32(getWord(s, 0x26)) / 1000
	mem.Type = parseMemoryType(getByte(s, 0x12))
	mem.TypeDetail = parseMemoryTypeDetail(getWord(s, 0x13))

	formFactor := parseFormFactor(getWord(s, 0x0E))
	if formFactor == "DIMM" {
		formFactor = getModuleType(mem.TypeDetail)
	}
	mem.FormFactor = formFactor

	log.Debugf("%+v", mem)

	return mem, nil
}

func parseMemorySize(b uint16, b2 uint32) uint32 {
	var size = uint32(b)
	if b == 0x7FFF {
		size = b2
	}

	// If the bit is 0, the value is specified in megabyte units.
	// if the bit is 1, the value is specified in kilobyte units.
	if b&(1<<15) != 0 {
		size = uint32(size / 1024 / 1024)
	} else {
		size = uint32(size / 1024)
	}

	return size
}

func parseMemoryType(b uint8) string {
	mtype := []string{
		"Other", // 0x01
		"Unknown",
		"DRAM",
		"EDRAM",
		"VRAM",
		"SRAM",
		"RAM",
		"ROM",
		"FLASH",
		"EEPROM",
		"FEPROM",
		"EPROM",
		"CDRAM",
		"3DRAM",
		"SDRAM", // 0x0F
		"SGRAM",
		"RDRAM",
		"DDR",
		"DDR2",
		"DDR2 FB-DIMM",
		"Reserved",
		"Reserved",
		"Reserved",
		"DDR3",
		"FBD2",
		"DDR4",
		"LPDDR",
		"LPDDR2",
		"LPDDR3",
		"LPDDR4", // 0x1E
	}

	mt := ""
	if b >= 0x01 && b <= 0x1E {
		mt = mtype[b-0x01]
	} else {
		log.Debugf("MemoryType: unsupported value %x was given.", b)
	}

	return mt
}

func parseFormFactor(b uint16) string {
	formfactor := []string{
		"Other", // 0x01
		"Unknown",
		"SIMM",
		"SIP",
		"Chip",
		"DIP",
		"ZIP",
		"Proprietary Card",
		"DIMM",
		"TSOP",
		"Row of chips",
		"RIMM",
		"SODIMM",
		"SRIMM",
		"FB-DIMM", // 0x0F
	}

	ff := ""
	if b >= 0x01 && b <= 0x0F {
		ff = formfactor[b-0x01]
	} else {
		log.Debugf("FormFactor: unsupported value %x was given.", b)
	}

	return ff
}

func parseMemoryTypeDetail(b uint16) []string {
	details := []string{
		"Other", // bit 1
		"Unknown",
		"Fast-paged",
		"Static column",
		"Pseudo-static",
		"RAMBUS",
		"Synchronous",
		"CMOS",
		"EDO",
		"Window DRAM",
		"Cache DRAM",
		"Non-volatile",
		"Registered (Buffered)",
		"Unbuffered (Unregistered)",
		"LRDIMM", // bit 15
	}

	d := []string{}
	if (b & 0xFFFE) == 0 { // bit 0 is reserved
		return d
	}

	var i uint16
	for i = 1; i <= 15; i++ {
		if (b & (1 << i)) != 0 {
			d = append(d, details[i-1])
		}
	}

	return d
}
