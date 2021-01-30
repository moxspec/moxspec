package main

import (
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/smbios"
)

func shapeSystem(r *model.Report, sm *smbios.System) {
	if sm == nil {
		return
	}

	s := new(model.System)
	s.ProductName = sm.ProductName
	s.Manufacturer = sm.Manufacturer
	s.SerialNumber = sm.SerialNumber
	r.System = s
}
