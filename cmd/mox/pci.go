package main

import (
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/pci"
)

func shapeAllPCIDevices(r *model.Report, devs *pci.Devices) {
	for _, d := range devs.AllDevices() {
		spec := shapePCIDevice(d)
		r.PCIDevice = append(r.PCIDevice, spec)
	}
}

func shapePCIDevice(dev *pci.Device) *model.PCIBaseSpec {
	p := new(model.PCIBaseSpec)

	p.Path = dev.Path
	p.Location.Domain = dev.Domain
	p.Location.Bus = dev.Bus
	p.Location.Device = dev.Device
	p.Location.Function = dev.Function
	p.Driver = dev.Driver
	p.VendorID = dev.VendorID
	p.VendorName = dev.VendorName
	p.DeviceID = dev.DeviceID
	p.DeviceName = dev.DeviceName
	p.SubSystemVendorID = dev.SubSystemVendorID
	p.SubSystemDeviceID = dev.SubSystemDeviceID
	p.SubSystemName = dev.SubSystemName
	p.ClassID = dev.ClassID
	p.ClassName = dev.ClassName
	p.SubClassID = dev.SubClassID
	p.SubClassName = dev.SubClassName
	p.InterfaceID = dev.InterfaceID
	p.InterfaceName = dev.InterfaceName

	p.Numa = byte(dev.Numa)
	p.SerialNumber = dev.SerialNumber
	if dev.LinkGen != 0 && dev.LinkSpeed != 0 && dev.LinkWidth != 0 {
		p.CurLink = &model.PCIeLink{
			Gen:   dev.LinkGen,
			Speed: dev.LinkSpeed,
			Width: dev.LinkWidth,
		}
	}
	if dev.MaxGen != 0 && dev.MaxSpeed != 0 && dev.MaxWidth != 0 {
		p.MaxLink = &model.PCIeLink{
			Gen:   dev.MaxGen,
			Speed: dev.MaxSpeed,
			Width: dev.MaxWidth,
		}
	}
	p.PowerLimit = dev.SlotPowetLimit
	p.UEList = dev.UncorrectableErrs
	p.CEList = dev.CorrectableErrs

	return p
}
