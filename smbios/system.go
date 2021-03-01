package smbios

import (
	gosmbios "github.com/digitalocean/go-smbios/smbios"
	"github.com/moxspec/moxspec/util"
)

// System represents a system spec
type System struct {
	Manufacturer string
	ProductName  string
	Version      string
	SerialNumber string
	SKUNumber    string
	Family       string
}

func parseSystem(s *gosmbios.Structure) *System {
	sinfo := new(System)

	sinfo.Manufacturer = util.ShortenVendorName(getStringsSet(s, 0x04))
	sinfo.ProductName = getStringsSet(s, 0x05)
	sinfo.Version = getStringsSet(s, 0x06)
	sinfo.SerialNumber = getStringsSet(s, 0x07)
	sinfo.SKUNumber = getStringsSet(s, 0x19)
	sinfo.Family = getStringsSet(s, 0x1A)

	log.Debugf("%+v", sinfo)

	return sinfo
}
