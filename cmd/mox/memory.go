package main

import (
	"github.com/moxspec/moxspec/edac"
	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/smbios"
)

func shapeMemory(r *model.Report, sm []*smbios.MemoryDevice) {
	if sm == nil {
		return
	}

	r.Memory = new(model.MemoryReport)

	var memories []*model.MemoryModule
	var empty byte
	var total uint64

	for _, mem := range sm {
		if mem.Size == 0 {
			empty++
		} else {

			m := new(model.MemoryModule)
			m.Locator = mem.DeviceLocator
			m.Manufacturer = mem.Manufacturer
			m.PartNumber = mem.PartNumber
			m.SerialNumber = mem.SerialNumber
			m.FormFactor = mem.FormFactor
			m.Type = mem.Type
			m.Speed = mem.Speed
			m.ConfiguredSpeed = mem.ConfiguredSpeed
			m.Size = uint64(mem.Size) * 1000 * 1000 * 1000 // Size should be bytes
			m.Voltage = mem.ConfiguredVoltage

			total += uint64(m.Size)
			memories = append(memories, m)
		}
	}

	r.Memory.Total = total // should be bytes
	r.Memory.Empty = empty
	r.Memory.Modules = memories

	edacd := edac.NewDecoder()
	err := edacd.Decode()
	if err != nil {
		log.Debug(err)
		return
	}

	var ctls []*model.MemoryController
	for _, mc := range edacd.Controllers {
		ctl := new(model.MemoryController)
		ctl.Name = mc.Name
		ctl.Size = mc.Size
		ctl.CECount = mc.CECount
		ctl.CENoInfoCount = mc.CENoInfoCount
		ctl.UECount = mc.UECount
		ctl.UENoInfoCount = mc.UENoInfoCount

		for _, csrow := range mc.CSRows {
			cs := new(model.ChipSelectRow)
			cs.Name = csrow.Name
			cs.Size = csrow.Size
			cs.CECount = csrow.CECount
			cs.UECount = csrow.UECount

			for _, c := range csrow.Channels {
				ch := new(model.MemoryChannel)
				ch.Name = c.Name
				ch.Label = c.Label
				ch.CECount = c.CECount

				cs.Channels = append(cs.Channels, ch)
			}

			ctl.CSRows = append(ctl.CSRows, cs)
		}

		ctls = append(ctls, ctl)
	}

	r.Memory.Controllers = ctls
}
