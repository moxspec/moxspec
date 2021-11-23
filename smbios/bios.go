package smbios

import (
	"fmt"

	gosmbios "github.com/digitalocean/go-smbios/smbios"
	"github.com/moxspec/moxspec/util"
)

// BIOS represents a BIOS spec
type BIOS struct {
	Vendor          string
	Version         string
	ReleaseDate     string
	MajorRelease    uint8
	MinorRelease    uint8
	Characteristics []string
}

func parseBIOS(s *gosmbios.Structure) (*BIOS, error) {
	if s == nil {
		return nil, fmt.Errorf("nil given")
	}

	bios := new(BIOS)

	bios.Vendor = util.ShortenVendorName(getStringsSet(s, 0x04))
	bios.Version = getStringsSet(s, 0x05)
	bios.ReleaseDate = getStringsSet(s, 0x08)
	bios.MajorRelease = getByte(s, 0x14)
	bios.MinorRelease = getByte(s, 0x15)

	log.Debugf("%+v", bios)

	return bios, nil
}
