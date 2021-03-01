package smbios

import (
	"encoding/binary"

	"github.com/moxspec/moxspec/util"
	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

// Structure represents a smbios table
type Structure struct {
	Header *Header
	Data   Data
}

// Header represents a smbios header
type Header struct {
	Type   uint8
	Handle uint16
}

// Data represents a smbios content as generic type
type Data interface{}

const headerSize = 4

// Table Type
const (
	biosInformation       = 0
	systemInformation     = 1
	baseboardInformation  = 2
	systemEnclosure       = 3
	processorInformation  = 4
	cacheInformation      = 7
	memoryDevice          = 17
	systemBootInformation = 32
	systemPowerSupply     = 39
)

func getByte(s *gosmbios.Structure, offset int) uint8 {
	o := offset - headerSize
	if o < 0 || o >= len(s.Formatted) {
		return 0
	}
	return uint8(s.Formatted[o])
}

func getWord(s *gosmbios.Structure, offset int) uint16 {
	o := offset - headerSize
	if o < 0 || o > len(s.Formatted)-2 {
		return 0
	}
	return binary.LittleEndian.Uint16(s.Formatted[o : o+2])
}

func getDWord(s *gosmbios.Structure, offset int) uint32 {
	o := offset - headerSize
	if o < 0 || o > len(s.Formatted)-4 {
		return 0
	}
	return binary.LittleEndian.Uint32(s.Formatted[o : o+4])
}

func getQWord(s *gosmbios.Structure, offset int) uint64 {
	o := offset - headerSize
	if o < 0 || o > len(s.Formatted)-8 {
		return 0
	}
	return binary.LittleEndian.Uint64(s.Formatted[o : o+8])
}

func getStringsSet(s *gosmbios.Structure, offset int) string {
	o := offset - headerSize
	if o >= len(s.Formatted) {
		return ""
	}

	ssLen := len(s.Strings)
	ptr := int(s.Formatted[o])
	if ptr > ssLen || ptr == 0 {
		return ""
	}

	return util.SanitizeString(s.Strings[ptr-1])
}
