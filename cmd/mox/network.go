package main

import (
	"github.com/actapio/moxspec/eth"
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/netlink"
	"github.com/actapio/moxspec/nw"
	"github.com/actapio/moxspec/pci"
)

func shapeNetwork(r *model.Report, pcidevs *pci.Devices) {
	r.Network = new(model.NetworkReport)

	var controllers []*model.EthController
	for _, ctl := range pcidevs.FilterByClass(pci.NetworkController) {
		c := new(model.EthController)
		c.PCIBaseSpec = *shapePCIDevice(ctl)

		nwd := nw.NewDecoder(c.Path, c.Driver)
		err := nwd.Decode()
		if err != nil {
			log.Debug(err)
			controllers = append(controllers, c)
			continue
		}

		intf := new(model.NetInterface)
		intf.State = nwd.Port.State
		intf.Name = nwd.Port.Name
		intf.HWAddr = nwd.Port.HWAddr
		intf.Speed = nwd.Port.Speed
		intf.MTU = nwd.Port.MTU

		for _, a := range nwd.Port.IPAddrs {
			addr := new(model.IPAddress)
			addr.Version = a.Ver
			addr.Addr = a.Addr
			addr.Netmask = a.Netmask
			addr.MaskSize = a.MaskSize
			addr.Network = a.Network
			addr.Broadcast = a.Broadcast

			intf.IPAddrs = append(intf.IPAddrs, addr)
		}

		ed := eth.NewDecoder(intf.Name)
		err = ed.Decode()
		if err == nil {
			intf.Speed = ed.Speed
			intf.SupportedSpeed = ed.SupportedSpeed
			intf.AdvertisingSpeed = ed.AdvertisingSpeed
			intf.FirmwareVersion = ed.FirmwareVersion

			if ed.Module != nil {
				mod := new(model.Module)
				mod.FormFactor = ed.Module.FormFactor
				mod.Connector = ed.Module.Connector
				mod.VendorName = ed.Module.VendorName
				mod.ProductName = ed.Module.ProductName
				mod.SerialNumber = ed.Module.SerialNumber
				mod.CableLength = ed.Module.CableLength
				intf.Module = mod
			}
		}

		nl := netlink.NewDecoder(intf.Name)
		err = nl.Decode()
		if err == nil {
			intf.RxErrors = nl.Stats.RxErrors
			intf.TxErrors = nl.Stats.TxErrors
			intf.RxDropped = nl.Stats.RxDropped
			intf.TxDropped = nl.Stats.TxDropped
		}

		c.Interfaces = append(c.Interfaces, intf)
		controllers = append(controllers, c)
	}

	r.Network.EthControllers = controllers
}
