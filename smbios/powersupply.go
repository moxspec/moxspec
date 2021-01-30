package smbios

import (
	"github.com/actapio/moxspec/util"
	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

// PowerSupply represents a power supply spec
type PowerSupply struct {
	Location         string
	DeviceName       string
	Manufacturer     string
	SerialNumber     string
	AssetTagNumber   string
	ModelPartNumber  string
	RevisionLevel    string
	MaxPowerCapacity uint16
	Plugged          bool
	Present          bool
	HotReplaceable   bool
}

func parsePowerSupply(s *gosmbios.Structure) *PowerSupply {
	ps := new(PowerSupply)

	ps.Location = getStringsSet(s, 0x05)
	ps.DeviceName = getStringsSet(s, 0x06)
	ps.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x07))
	ps.SerialNumber = getStringsSet(s, 0x08)
	ps.AssetTagNumber = getStringsSet(s, 0x09)
	ps.ModelPartNumber = getStringsSet(s, 0x0A)
	ps.RevisionLevel = getStringsSet(s, 0x0B)
	ps.MaxPowerCapacity = parsePowerSupplyCap(getWord(s, 0x0C))

	ps.Plugged, ps.Present, ps.HotReplaceable = parsePowerSupplyStats(getWord(s, 0x0E))

	log.Debugf("%+v", ps)

	return ps
}

func parsePowerSupplyCap(w uint16) uint16 {
	// Set to 0x8000 if unknown.
	if w == 0x8000 {
		return 0
	}
	return uint16(w)
}

func parsePowerSupplyStats(w uint16) (plugged, present, hotswap bool) {
	bits := w & 0x07
	log.Debugf("parsing power supply stat: %d", bits)

	if (bits & 0x01) != 0 {
		hotswap = true
	}

	if (bits & 0x02) != 0 {
		present = true
	}

	// bit 2: 1b power supply is ***unplugged*** from the wall
	if (bits & 0x04) == 0 {
		plugged = true
	}

	return
}
