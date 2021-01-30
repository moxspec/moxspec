package main

import (
	"encoding/json"
	"fmt"

	"github.com/actapio/moxspec/model"
)

func show(cli *app) error {
	var err error

	var r *model.Report
	r, err = decode(cli)

	if err != nil {
		return err
	}

	if cli.getBool("j") {
		jb, err := json.Marshal(r)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", jb)
		return nil
	}

	p := newPrinter()
	writeDownSystem(r, p)
	writeDownFirmware(r, p)
	writeDownBaseboard(r, p)
	writeDownProcessor(r, p)
	writeDownMemory(r, p)
	writeDownDisk(r, p)
	writeDownNetwork(r, p)
	writeDownAccelerator(r, p)
	writeDownBMC(r, p)
	writeDownPowerSupply(r, p)
	writeDownPlatform(r, p)
	if cli.getString("h") != "" {
		writeDownLastUpdate(r, p)
	}
	p.show()

	return nil
}

func writeDownSystem(r *model.Report, p *printer) {
	sy := r.System
	if sy == nil {
		log.Debug("could not decode system information")
		return
	}

	s := newSection("System")
	s.block.append(sy.Summary())
	p.append(s)
}

func writeDownFirmware(r *model.Report, p *printer) {
	f := r.Firmware
	if f == nil {
		log.Debug("could not decode firmware information")
		return
	}

	s := newSection(string(f.Type))
	s.block.append(f.Summary())
	p.append(s)
}

func writeDownBaseboard(r *model.Report, p *printer) {
	b := r.Baseboard
	if b == nil {
		log.Debug("could not decode baseboard information")
		return
	}

	s := newSection("Baseboard")
	s.block.append(b.Summary())
	p.append(s)
}

func writeDownProcessor(r *model.Report, p *printer) {
	s := newSection("Processor")
	s.block.append(r.Processor.Summary())
	p.append(s)

	sbnd := new(block)
	for _, p := range r.Processor.Packages {
		for _, n := range p.Nodes {
			sbnd.append(fmt.Sprintf("package%d %s", p.ID, n.Summary()))
		}
	}
	s.block.append(newGroupedBlock("Node", sbnd))

	// TODO: to be improved
	proc := r.Processor.Packages[0]
	if len(proc.Caches) == 0 || len(proc.TLBs) == 0 {
		return
	}

	sbcc := new(block)
	for _, c := range proc.Caches {
		sbcc.appendf(c.Summary())
	}
	s.block.append(newGroupedBlock("Cache", sbcc))

	sbtt := new(block)
	for _, t := range proc.TLBs {
		sbtt.appendf(t.Summary())
	}
	s.block.append(newGroupedBlock("TLB", sbtt))

}

func writeDownMemory(r *model.Report, p *printer) {
	s := newSection("Memory")
	s.block.appendf("Total: %s", r.Memory.TotalString())
	s.block.append(newIndentedBlock(r.Memory.ModuleSummaries()))

	if r.Memory.HasDiag() {
		if r.Memory.IsHealthy() {
			s.block.append("Diag: healthy")
		} else {
			s.block.append("Diag: UNHEALTHY")
			s.block.append(newIndentedBlock(r.Memory.DiagSummaries()))
		}
	}

	p.append(s)
}

func writeDownDisk(r *model.Report, p *printer) {
	s := newSection("Disk")

	if r.Storage.AHCIControllers != nil {
		writeDownAHCIControllers(r, s)
	}
	if r.Storage.RAIDControllers != nil {
		writeDownRAIDControllers(r, s)
	}
	if r.Storage.NVMeControllers != nil {
		writeDownNVMeControllers(r, s)
	}
	if r.Storage.VirtControllers != nil {
		writeDownVirtControllers(r, s)
	}
	if r.Storage.NonStdControllers != nil {
		writeDownNonStdControllers(r, s)
	}

	p.append(s)
}

func writeDownAHCIControllers(r *model.Report, s *section) {
	for _, ctl := range r.Storage.AHCIControllers {
		s.block.appendf(ctl.Summary())

		sb := new(block)
		if ctl.HasLinkStatus() {
			sb.appendf("Link: %s", ctl.LinkSummary())
		}
		if ctl.IsHealthy() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(ctl.DiagSummaries()))
		}
		s.block.append(sb)

		if len(ctl.Drives) == 0 {
			continue
		}

		var bs []*block
		for _, d := range ctl.Drives {
			sbbd := new(block)

			if d.IsSSD() || d.IsHDD() {
				sbbd.appendf("Temp: %s", d.TempMaxMinSummary())
				sbbd.appendf("Wear: %s", d.IOSummary())
				sbbd.appendf("Form: %s", d.FormSummary())
				sbbd.appendf("Firm: %s", d.Firmware)

				if d.IsHealthy() {
					sbbd.append("Diag: healthy")
				} else {
					sbbd.append("Diag: UNHEALTHY")
					sbbd.append(newIndentedBlock(d.DiagSummaries()))
				}
			}

			if sbbd.hasContents() {
				bs = append(bs, newGroupedBlock(d.Summary(), sbbd))
			} else {
				// optial drives or others
				misc := new(block)
				misc.append(d.Summary())
				bs = append(bs, misc)
			}

		}
		s.block.append(newGroupedBlock(fmt.Sprintf("Drive: %d", len(ctl.Drives)), bs...))
	}
}

func writeDownRAIDControllers(r *model.Report, s *section) {
	for _, ctl := range r.Storage.RAIDControllers {
		s.block.appendf(ctl.Summary())

		sb := new(block)
		if ctl.HasLinkStatus() {
			sb.appendf("Link: %s", ctl.LinkSummary())
		}

		sb.append(fmt.Sprintf("Spec: %s", ctl.SpecSummary()))

		if ctl.IsHealthy() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(ctl.DiagSummaries()))
		}
		s.block.append(sb)

		vdb := new(block)
		for _, ld := range ctl.LogDrives {
			vdb.append(ld.LDSummary())

			// this means that a topology parser did not work
			// so the mox can not display more information about it
			if len(ld.PhyDrives) == 0 {
				log.Debug("this ld has no pd info")
				continue
			}

			vdsb := new(block)
			vdsb.append(fmt.Sprintf("Stat: %s", ld.Status))

			if ld.IsHealthy() {
				vdsb.append("Diag: healthy")
			} else {
				vdsb.append("Diag: UNHEALTHY")
			}

			vdb.append(vdsb)

			pdb := new(block)
			for _, pd := range ld.PhyDrives {
				pdb.append(pd.Summary())
			}

			vdb.append(newGroupedBlock(fmt.Sprintf("Physical Drive: %d", len(ld.PhyDrives)), pdb))

		}

		if len(ctl.LogDrives) > 0 {
			s.block.append(newGroupedBlock(fmt.Sprintf("Logical Drive: %d", len(ctl.LogDrives)), vdb))
		}

		ptdb := new(block)
		for _, pt := range ctl.PassthroughDrives {
			ptdb.append(pt.PTSummary())
		}
		if len(ctl.PassthroughDrives) > 0 {
			s.block.append(newGroupedBlock(fmt.Sprintf("Pass-Through Drive: %d", len(ctl.PassthroughDrives)), ptdb))
		}

		ucdb := new(block)
		for _, pd := range ctl.UnconfDrives {
			ucdb.append(pd.Summary())
		}
		if len(ctl.UnconfDrives) > 0 {
			s.block.append(newGroupedBlock(fmt.Sprintf("Unconfigured Drive: %d", len(ctl.UnconfDrives)), ucdb))
		}
	}
}

func writeDownNVMeControllers(r *model.Report, s *section) {
	for _, ctl := range r.Storage.NVMeControllers {
		s.block.appendf(ctl.Summary())

		sb := new(block)

		if ctl.HasLinkStatus() {
			sb.appendf("Link: %s", ctl.LinkSummary())
		}

		sb.appendf("Temp: %s", ctl.TempWarnCritSummary())
		sb.appendf("Wear: %s", ctl.IOSummary())
		sb.appendf("Firm: %s", ctl.Firmware)
		if ctl.IsHealthy() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(ctl.DiagSummaries()))
		}

		s.block.append(sb)

		if len(ctl.Namespaces) == 0 {
			continue
		}

		sb = new(block)
		for _, ns := range ctl.Namespaces {
			sb.append(ns.Summary())
		}
		s.block.append(newGroupedBlock(fmt.Sprintf("Namespace: %d", len(ctl.Namespaces)), sb))
	}
}

func writeDownVirtControllers(r *model.Report, s *section) {
	for _, ctl := range r.Storage.VirtControllers {
		s.block.appendf(ctl.Summary())

		if len(ctl.Drives) == 0 {
			continue
		}

		db := new(block)
		for _, d := range ctl.Drives {
			db.append(d.Summary())
		}
		s.block.append(newGroupedBlock(fmt.Sprintf("Drive: %d", len(ctl.Drives)), db))
	}
}

func writeDownNonStdControllers(r *model.Report, s *section) {
	for _, ctl := range r.Storage.NonStdControllers {
		s.block.appendf(ctl.Summary())

		sb := new(block)

		if ctl.HasLinkStatus() {
			sb.appendf("Link: %s", ctl.LinkSummary())
		}

		sb.appendf("Temp: %s", ctl.TempSummary())
		sb.appendf("Wear: %s", ctl.IOSummary())
		sb.appendf("Firm: %s", ctl.Firmware)
		if ctl.IsHealthy() && !ctl.HasErrorRecords() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(append(ctl.DiagSummaries(), ctl.ErrorRecords...)))
		}

		s.block.append(sb)

		if len(ctl.Drives) == 0 {
			continue
		}

		db := new(block)
		for _, d := range ctl.Drives {
			db.append(d.Summary())
		}
		s.block.append(newGroupedBlock(fmt.Sprintf("Drive: %d", len(ctl.Drives)), db))
	}
}

func writeDownNetwork(r *model.Report, p *printer) {
	if r.Network.EthControllers == nil {
		return
	}

	s := newSection("Network")
	for _, ctl := range r.Network.EthControllers {
		s.block.appendf(ctl.Summary())

		sb := new(block)

		if ctl.HasLinkStatus() {
			sb.appendf("Link: %s", ctl.LinkSummary())
		}

		if ctl.HasDriver() {
			for _, intf := range ctl.Interfaces {
				sb.appendf("Intf: %s", intf.Summary())
				sb.appendf("Stat: %s", intf.StatSummary())
				if intf.Module != nil {
					sb.appendf("Modl: %s", intf.Module.Summary())
				}
			}
		}

		if ctl.IsHealthy() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(ctl.DiagSummaries()))
		}

		s.block.append(sb)
	}
	p.append(s)
}

func writeDownAccelerator(r *model.Report, p *printer) {
	s := newSection("Accelerator")

	if r.Accelerator == nil {
		return
	}

	for _, g := range r.Accelerator.GPUs {
		s.block.appendf(g.Summary())

		sb := new(block)
		s.block.append(sb)

		if g.HasLinkStatus() {
			sb.appendf("Link: %s", g.LinkSummary())
		}

		sb.appendf("Powr: %s", g.PowerSummary())
		sb.appendf("Temp: %s", g.TempSummary())
		sb.appendf("Firm: %s", g.BIOS)

		if g.IsHealthy() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(g.DiagSummaries()))
		}

		s.block.append(sb)
	}

	for _, f := range r.Accelerator.FPGAs {
		s.block.appendf(f.Summary())

		sb := new(block)
		s.block.append(sb)

		if f.HasLinkStatus() {
			sb.appendf("Link: %s", f.LinkSummary())
		}

		if f.IsHealthy() {
			sb.append("Diag: healthy")
		} else {
			sb.append("Diag: UNHEALTHY")
			sb.append(newIndentedBlock(f.DiagSummaries()))
		}

		s.block.append(sb)
	}

	p.append(s)
}

func writeDownPowerSupply(r *model.Report, p *printer) {
	s := newSection("Power Supply")
	for _, ps := range r.PowerSupply {
		s.block.append(ps.Summary())
		sb := new(block)
		sb.append(fmt.Sprintf("Stat: %s", ps.Stat()))
		s.block.append(sb)
	}
	p.append(s)
}

func writeDownBMC(r *model.Report, p *printer) {
	if r.BMC == nil {
		return
	}

	s := newSection("BMC")
	s.block.appendf("Intf: %s, %s/%d", r.BMC.MAC, r.BMC.IPAddr, r.BMC.MaskSize)
	s.block.appendf("Firm: %s", r.BMC.Firmware)
	p.append(s)
}

func writeDownPlatform(r *model.Report, p *printer) {
	sos := newSection("OS")
	sos.block.appendf("%s, %s", r.OS.Distro, r.OS.Kernel)
	p.append(sos)

	cl := newSection("Client")
	cl.block.appendf("v%s", r.Version)
	p.append(cl)

	sh := newSection("Hostname")
	sh.block.append(r.Hostname)
	p.append(sh)
}

func writeDownLastUpdate(r *model.Report, p *printer) {
	lu := newSection("Last Update")
	lu.block.appendf("%s", r.Datetime)
	p.append(lu)
}
