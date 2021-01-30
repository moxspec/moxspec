package main

import (
	"os"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/model"
)

func lsraid() error {
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

	if r.Storage == nil {
		return nil
	}

	if len(r.Storage.RAIDControllers) == 0 {
		return nil
	}

	tbl := newTable("blk", "conf", "adp", "pos", "stat", "size", "form", "model")
	for _, ctl := range r.Storage.RAIDControllers {
		if len(ctl.LogDrives) == 0 && len(ctl.UnconfDrives) == 0 && len(ctl.PassthroughDrives) == 0 {
			continue
		}

		for _, ld := range ctl.LogDrives {
			if len(ld.PhyDrives) == 0 {
				continue
			}

			for _, pd := range ld.PhyDrives {
				tbl.append(ld.Name, ld.RAIDLv, ctl.AdapterID, pd.Pos(), pd.Status, pd.SizeString(), pd.FormSummary(), pd.Model)
			}
		}

		for _, pd := range ctl.UnconfDrives {
			tbl.append("", "unconf", ctl.AdapterID, pd.Pos(), pd.Status, pd.SizeString(), pd.FormSummary(), pd.Model)
		}

		for _, pd := range ctl.PassthroughDrives {
			tbl.append(pd.Name, "Pass-Through", ctl.AdapterID, pd.Pos(), pd.Status, pd.SizeString(), pd.FormSummary(), pd.Model)
		}
	}
	tbl.print()

	return nil
}

func strSliceToIntfSlice(in []string) []interface{} {
	var out []interface{}
	for _, i := range in {
		out = append(out, interface{}(i))
	}
	return out
}
