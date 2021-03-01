package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/model"
)

func lssn() error {
	loglet.SetLevel(loglet.INFO)

	cli := newAppWithoutCmd(os.Args)
	cli.appendFlag("h", "", "hostname")
	err := cli.parse()
	if err != nil {
		return err
	}

	var r *model.Report
	r, err = decode(cli)

	if err != nil {
		return err
	}

	tbl := newTable("Catagory", "Model", "Serial Number", "Location", "Spec")

	if r.Processor != nil {
		writeDownProcessorSerialNumber(tbl, r.Processor)
	}

	if r.Memory != nil {
		writeDownMemorySerialNumber(tbl, r.Memory)
	}

	if r.Storage != nil {
		writeDownStorageSerialNumber(tbl, r.Storage)
	}

	if r.Network != nil {
		writeDownNetworkSerialNumber(tbl, r.Network)
	}

	writeDownSystemSerialNumber(tbl, r)
	tbl.print()

	return nil
}

func writeDownSystemSerialNumber(tbl *table, r *model.Report) {
	if r.System != nil {
		model := fmt.Sprintf("%s %s", r.System.Manufacturer, r.System.ProductName)
		tbl.append("System", model, r.System.SerialNumber, "", "")
	}

	if r.Chassis != nil {
		tbl.append("Chassis", "", r.Chassis.SerialNumber, "", "")
	}

	if r.Baseboard != nil {
		model := fmt.Sprintf("%s %s", r.Baseboard.Manufacturer, r.Baseboard.ProductName)
		tbl.append("Baseboard", model, r.Baseboard.SerialNumber, "", "")
	}

	for _, p := range r.PowerSupply {
		model := fmt.Sprintf("%s %s", p.Manufacturer, p.ProductName)
		spec := fmt.Sprintf("%dW", p.Capacity)
		tbl.append("PSU", model, p.SerialNumber, "", spec)
	}
}

func writeDownProcessorSerialNumber(tbl *table, r *model.ProcessorReport) {
	for _, p := range r.Packages {
		spec := fmt.Sprintf("%dcores, %dthreads", p.CoreCount, p.ThreadCount)
		tbl.append("Processor", p.ProductName, p.SerialNumber, p.Socket, spec)
	}
}

func writeDownMemorySerialNumber(tbl *table, r *model.MemoryReport) {
	for _, m := range r.Modules {
		spec := fmt.Sprintf("%s-%d %s", m.Type, m.Speed, m.SizeString())
		model := fmt.Sprintf("%s %s", m.Manufacturer, m.PartNumber)
		tbl.append("Memory", model, m.SerialNumber, m.Locator, spec)
	}
}

func writeDownNVMeSerialNumber(tbl *table, c *model.NVMeController) {
	spec := fmt.Sprintf("NVMe SSD %s", c.SizeString())
	model := fmt.Sprintf("%s %s", c.VendorName, c.Model)
	tbl.append("Storage", model, c.SerialNumber, "", spec)
}

func writeDownPhyDriveSerialNumber(tbl *table, d *model.PhyDrive) {
	media := "HDD"
	if d.SolidStateDrive {
		media = "SSD"
	}
	spec := fmt.Sprintf("%s %s %s", d.Transport, media, d.SizeString())
	location := fmt.Sprintf("Enc %s - Slot %s", d.Enclosure, d.Slot)
	tbl.append("Storage", d.Model, d.SerialNumber, location, spec)
}

func writeDownDriveSerialNumber(tbl *table, d *model.Drive) {
	media := ""
	if d.IsHDD() {
		media = "HDD"
	} else if d.IsSSD() {
		media = "SSD"
	}
	transport := strings.Split(d.Transport, " ")[0]
	spec := fmt.Sprintf("%s %s %s", transport, media, d.SizeString())
	tbl.append("Storage", d.Model, d.SerialNumber, "", spec)
}

func writeDownStorageSerialNumber(tbl *table, r *model.StorageReport) {
	for _, c := range r.NVMeControllers {
		writeDownNVMeSerialNumber(tbl, c)
	}

	for _, c := range r.RAIDControllers {
		tbl.append("RAID Card", c.ProductName, c.SerialNumber, "", "")
		for _, ld := range c.LogDrives {
			for _, pd := range ld.PhyDrives {
				writeDownPhyDriveSerialNumber(tbl, pd)
			}
		}

		for _, pd := range c.UnconfDrives {
			writeDownPhyDriveSerialNumber(tbl, pd)
		}

		for _, pd := range c.PassthroughDrives {
			writeDownPhyDriveSerialNumber(tbl, pd)
		}
	}

	for _, c := range r.AHCIControllers {
		for _, d := range c.Drives {
			writeDownDriveSerialNumber(tbl, d)
		}
	}
}

func writeDownNetworkSerialNumber(tbl *table, r *model.NetworkReport) {
	for _, c := range r.EthControllers {
		tbl.append("Network", c.LongName(), c.SerialNumber, "", "")
	}
}
