package main

import (
	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/smbios"
)

func shapeBaseboard(r *model.Report, sm *smbios.Baseboard) {
	if sm == nil {
		return
	}

	b := new(model.Baseboard)
	b.Manufacturer = sm.Manufacturer
	b.SerialNumber = sm.SerialNumber
	b.ProductName = sm.Product
	r.Baseboard = b
}
