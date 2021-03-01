package main

import (
	"github.com/moxspec/moxspec/ipmi"
	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/util"
)

func shapeBMC(r *model.Report) {
	d := ipmi.NewDecoder()
	err := d.Decode()
	if err != nil {
		log.Debug(err) // it's not fatal (eg. virtual machine)
		return
	}

	b := new(model.BMC)
	b.Type = "IPMI"
	b.Firmware = d.Firmware
	b.MAC = d.MAC
	b.IPAddr = d.IPAddr
	b.Netmask = d.Netmask
	b.MaskSize = util.IPv4MaskSize(b.Netmask)
	b.Gateway = d.Gateway

	r.BMC = b
}
