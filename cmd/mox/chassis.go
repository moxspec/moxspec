package main

import (
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/smbios"
)

func shapeChassis(r *model.Report, sm *smbios.Chassis) {
	if sm == nil {
		return
	}

	c := new(model.Chassis)
	c.Manufacturer = sm.Manufacturer
	c.SerialNumber = sm.SerialNumber
	r.Chassis = c
}
