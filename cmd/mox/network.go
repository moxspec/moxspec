package main

import (
	"net"

	"github.com/moxspec/moxspec/bonding"
	"github.com/moxspec/moxspec/eth"
	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/netlink"
	"github.com/moxspec/moxspec/nw"
	"github.com/moxspec/moxspec/pci"
)

func shapeNetwork(r *model.Report, pcidevs *pci.Devices) {
	r.Network = new(model.NetworkReport)

	var bondings []*model.BondInterface
	var controllers []*model.EthController

	for _, bond := range bonding.GetBondDevices() {
		c := new(model.BondInterface)
		c.Name = bond

		bd := bonding.NewDecoder(bond)
		err := bd.Decode()
		if err != nil {
			log.Debug(err)
			bondings = append(bondings, c)
			continue
		}

		c.Slaves = bd.Slaves

		c.LinkAttrs.State = bd.LinkAttrs.OperState.String()
		c.LinkAttrs.HWAddr = bd.LinkAttrs.HardwareAddr.String()
		c.LinkAttrs.MTU = bd.LinkAttrs.MTU
		c.LinkAttrs.TxQLen = bd.LinkAttrs.TxQLen

		c.BondAttrs.Mode = bd.BondAttrs.Mode.String()
		c.BondAttrs.ActiveSlave = bd.BondAttrs.ActiveSlave
		c.BondAttrs.Miimon = bd.BondAttrs.Miimon
		c.BondAttrs.UpDelay = bd.BondAttrs.UpDelay
		c.BondAttrs.DownDelay = bd.BondAttrs.DownDelay
		c.BondAttrs.UseCarrier = bd.BondAttrs.UseCarrier
		c.BondAttrs.ArpInterval = bd.BondAttrs.ArpInterval
		c.BondAttrs.ArpIpTargets = bd.BondAttrs.ArpIpTargets
		c.BondAttrs.ArpValidate = bd.BondAttrs.ArpValidate.String()
		c.BondAttrs.ArpAllTargets = bd.BondAttrs.ArpAllTargets.String()
		c.BondAttrs.Primary = bd.BondAttrs.Primary
		c.BondAttrs.PrimaryReselect = bd.BondAttrs.PrimaryReselect.String()
		c.BondAttrs.FailOverMac = bd.BondAttrs.FailOverMac.String()
		c.BondAttrs.XmitHashPolicy = bd.BondAttrs.XmitHashPolicy.String()
		c.BondAttrs.LacpRate = bd.BondAttrs.LacpRate.String()

		for _, ipaddress := range bd.AddrList {
			ipaddr := new(model.IPAddress)
			ipaddr.Addr = ipaddress.IP.String()

			if ipaddress.IP.To4() == nil {
				ipaddr.Version = 6
			} else {
				ipaddr.Version = 4
			}

			ipaddr.Netmask = net.IP(ipaddress.Mask).String()

			ipaddr.MaskSize, _ = ipaddress.Mask.Size()

			ipaddr.Broadcast = ipaddress.Broadcast.String()

			ipaddr.Network = ipaddress.IPNet.String()

			c.LinkAttrs.IPAddrs = append(c.LinkAttrs.IPAddrs, ipaddr)
		}

		bondings = append(bondings, c)
	}

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
		} else {
			log.Debug(err)
		}

		nl := netlink.NewDecoder(intf.Name)
		err = nl.Decode()
		if err == nil {
			intf.RxErrors = nl.Stats.RxErrors
			intf.TxErrors = nl.Stats.TxErrors
			intf.RxDropped = nl.Stats.RxDropped
			intf.TxDropped = nl.Stats.TxDropped
		} else {
			log.Debug(err)
		}

		c.Interfaces = append(c.Interfaces, intf)
		controllers = append(controllers, c)
	}

	r.Network.EthControllers = controllers
	r.Network.BondInterfaces = bondings
}
