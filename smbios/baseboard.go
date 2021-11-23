package smbios

import (
	"fmt"

	gosmbios "github.com/digitalocean/go-smbios/smbios"
	"github.com/moxspec/moxspec/util"
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

func parseBaseboard(s *gosmbios.Structure) (*Baseboard, error) {
	if s == nil {
		return nil, fmt.Errorf("nil given")
	}

	bboard := new(Baseboard)

	bboard.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x04))
	bboard.Product = getStringsSet(s, 0x05)
	bboard.Version = getStringsSet(s, 0x06)
	bboard.SerialNumber = getStringsSet(s, 0x07)
	bboard.AssetTag = getStringsSet(s, 0x08)

	log.Debugf("%+v", bboard)
	return bboard, nil
}
