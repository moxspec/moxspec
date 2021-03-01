package main

import (
	"strings"

	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/smbios"
)

func shapePowerSupply(r *model.Report, sm []*smbios.PowerSupply) {
	if sm == nil {
		return
	}

	for _, ps := range sm {
		if ps.DeviceName == "" || strings.Contains(ps.DeviceName, "OEM") {
			continue
		}
		if ps.Manufacturer == "" || strings.Contains(ps.Manufacturer, "OEM") {
			continue
		}

		p := new(model.PowerSupply)
		p.ProductName = ps.DeviceName
		p.Manufacturer = ps.Manufacturer
		p.SerialNumber = ps.SerialNumber
		p.ModelPartNumber = ps.ModelPartNumber
		p.Capacity = ps.MaxPowerCapacity
		p.Present = ps.Present
		p.Plugged = ps.Plugged
		p.HotSwappable = ps.HotReplaceable

		r.PowerSupply = append(r.PowerSupply, p)
	}
}
