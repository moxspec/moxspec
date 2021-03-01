package main

import (
	"fmt"
	"os"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/model"
)

const (
	healthy   = "healthy"
	unhealthy = "UNHEALTHY"
)

const (
	exitHealthy   = 0
	exitUnhealthy = 1
)

func lsdiag() (int, error) {
	loglet.SetLevel(loglet.INFO)

	cli := newAppWithoutCmd(os.Args)
	cli.appendFlag("h", "", "hostname")
	err := cli.parse()
	if err != nil {
		return exitUnhealthy, err
	}

	var r *model.Report
	r, err = decode(cli)

	if err != nil {
		return exitUnhealthy, err
	}

	exitCode := exitHealthy

	tbl := newTable("category", "stat", "detail")

	if r.Processor != nil {
		for _, p := range r.Processor.Packages {
			model := fmt.Sprintf("%s %s", p.Socket, p.ProductName)

			stat := healthy
			detail := model

			if p.ThrottleCount > 0 {
				stat = unhealthy
				detail = fmt.Sprintf("%s: throttle cnt = %d", model, p.ThreadCount)

				if exitCode != exitUnhealthy {
					exitCode = exitUnhealthy
				}
			}

			tbl.append("Processor", stat, detail)
		}
	}

	if r.Memory != nil {
		append := func(stat, detail string) {
			tbl.append("Memory", stat, detail)
		}

		for _, ctl := range r.Memory.Controllers {
			ctlName := ctl.Name
			for _, cs := range ctl.CSRows {
				if cs.HasError() {
					append(unhealthy, fmt.Sprintf("%s %s", ctlName, cs.Summary()))

					if exitCode != exitUnhealthy {
						exitCode = exitUnhealthy
					}
				} else {
					append(healthy, fmt.Sprintf("%s %s", ctlName, cs.Name))
				}
			}
		}
	}
	if r.Storage != nil {
		for _, ctl := range r.Storage.RAIDControllers {
			if ctl.IsHealthy() {
				tbl.append("RAID Card", healthy, ctl.LongName())
			} else {
				appendMultiDiags(tbl, "RAID Card", unhealthy, ctl.DiagSummaries())

				if exitCode != exitUnhealthy {
					exitCode = exitUnhealthy
				}
			}

			for _, ld := range ctl.LogDrives {
				detail := fmt.Sprintf("%s: %s, %s", ld.Name, ld.RAIDLv, ld.Status)

				stat := healthy
				if !ld.IsHealthy() {
					stat = unhealthy

					if exitCode != exitUnhealthy {
						exitCode = exitUnhealthy
					}
				}

				tbl.append("RAID Volume", stat, detail)
			}

			for _, pd := range ctl.PassthroughDrives {
				stat := healthy
				if !pd.IsHealthy() {
					stat = unhealthy

					if exitCode != exitUnhealthy {
						exitCode = exitUnhealthy
					}
				}
				detail := fmt.Sprintf("%s: %s, %s", pd.Name, pd.Model, pd.Status)
				tbl.append("Pass-Through Drive", stat, detail)
			}
		}

		for _, ctl := range r.Storage.NVMeControllers {
			if ctl.IsHealthy() {
				tbl.append("NVMe Drive", healthy, ctl.LongName())
			} else {
				appendMultiDiags(tbl, "NVMe Drive", unhealthy, ctl.DiagSummaries())

				if exitCode != exitUnhealthy {
					exitCode = exitUnhealthy
				}
			}
		}

		for _, ctl := range r.Storage.AHCIControllers {
			for _, drv := range ctl.Drives {
				detail := fmt.Sprintf("%s %s", drv.Model, drv.SizeString())

				if drv.IsHealthy() {
					tbl.append("SATA Drive", healthy, detail)
				} else {
					appendMultiDiags(tbl, "SATA Drive", unhealthy, ctl.DiagSummaries())

					if exitCode != exitUnhealthy {
						exitCode = exitUnhealthy
					}
				}
			}
		}
	}

	if r.Network != nil {
		for _, ctl := range r.Network.EthControllers {
			detail := ctl.LongName()

			if ctl.IsHealthy() {
				tbl.append("Network", healthy, detail)
			} else {
				appendMultiDiags(tbl, "Network", unhealthy, ctl.DiagSummaries())

				if exitCode != exitUnhealthy {
					exitCode = exitUnhealthy
				}
			}
		}
	}

	if r.Accelerator != nil {
		for _, g := range r.Accelerator.GPUs {
			if g.IsHealthy() {
				tbl.append("Accelerator", healthy, g.LongName())
			} else {
				appendMultiDiags(tbl, "Accelerator", unhealthy, g.DiagSummaries())

				if exitCode != exitUnhealthy {
					exitCode = exitUnhealthy
				}
			}
		}

		for _, f := range r.Accelerator.FPGAs {
			if f.IsHealthy() {
				tbl.append("Accelerator", healthy, f.LongName())
			} else {
				appendMultiDiags(tbl, "Accelerator", unhealthy, f.DiagSummaries())

				if exitCode != exitUnhealthy {
					exitCode = exitUnhealthy
				}
			}
		}
	}

	tbl.print()
	return exitCode, nil
}

func appendMultiDiags(t *table, cat, stat string, diags []string) {
	for i, d := range diags {
		if i == 0 {
			t.append(cat, stat, d)
		} else {
			t.append("", "", d)
		}
	}
}
