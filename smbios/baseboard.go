package smbios

import (
	"github.com/actapio/moxspec/util"
	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

// Baseboard represents a baseboard spec
type Baseboard struct {
	Manufacturer string
	Product      string
	Version      string
	SerialNumber string
	AssetTag     string
	BoardType    string
}

func parseBaseboard(s *gosmbios.Structure) *Baseboard {
	bboard := new(Baseboard)

	bboard.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x04))
	bboard.Product = getStringsSet(s, 0x05)
	bboard.Version = getStringsSet(s, 0x06)
	bboard.SerialNumber = getStringsSet(s, 0x07)
	bboard.AssetTag = getStringsSet(s, 0x08)

	log.Debugf("%+v", bboard)
	return bboard
}
