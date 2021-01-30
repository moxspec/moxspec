package smbios

import (
	"github.com/actapio/moxspec/util"
	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

// Chassis represents a chassis spec
type Chassis struct {
	Manufacturer       string
	Version            string
	SerialNumber       string
	AssetTagNumber     string
	BootUpState        string
	PowerSupplyState   string
	ThermalState       string
	Height             uint8 // OCP server does NOT support this field
	NumberOfPowerCords uint8
	SKUNumber          string
}

func parseChassis(s *gosmbios.Structure) *Chassis {
	c := new(Chassis)

	c.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x04))
	c.Version = getStringsSet(s, 0x06)
	c.SerialNumber = getStringsSet(s, 0x07)
	c.AssetTagNumber = getStringsSet(s, 0x08)
	c.BootUpState = parseChassisState(getByte(s, 0x09))
	c.PowerSupplyState = parseChassisState(getByte(s, 0x0A))
	c.ThermalState = parseChassisState(getByte(s, 0x0B))
	c.Height = getByte(s, 0x11)
	c.NumberOfPowerCords = getByte(s, 0x12)
	c.SKUNumber = getStringsSet(s, 0x15)

	log.Debugf("%+v", c)

	return c
}

func parseChassisState(b uint8) string {
	status := []string{
		"Other", // 0x01
		"Unknown",
		"Safe",
		"Warning",
		"Critical",
		"Non-recorverable", // 0x06
	}

	if b >= 0x01 && b <= 0x06 {
		return status[b-0x01]
	}
	return ""
}
