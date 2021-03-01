package main

import (
	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/smbios"
)

func shapeFirmware(r *model.Report, sm *smbios.BIOS) {
	if sm == nil {
		return
	}

	f := new(model.Firmware)
	f.Type = model.BIOS
	f.Vendor = sm.Vendor
	f.Version = sm.Version
	f.ReleaseDate = sm.ReleaseDate
	r.Firmware = f
}
