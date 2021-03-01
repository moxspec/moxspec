package hpacucli

import (
	"fmt"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/raidcli"
	"github.com/moxspec/moxspec/util"
)

var pathList = []string{
	"/usr/sbin/hpssacli",
}

var clipath string

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("hpacucli")
	clipath, _ = util.ScanPathList(pathList)
}

// Available returns whether raid command is available
func Available() bool {
	if clipath == "" {
		return false
	}
	return true
}

// GetControllers retruns hpacucli controllers
func GetControllers() ([]*Controller, error) {
	res, err := raidcli.Run(clipath, "controller", "all", "show")
	if err != nil {
		return nil, err
	}
	return parseCtlList(res)
}

// Decode makes Adapter satisfy the mox.Decoder interface
func (c *Controller) Decode() error {
	res, err := raidcli.Run(clipath, "controller", fmt.Sprintf("slot=%s", c.Slot), "show", "config", "detail")
	if err != nil {
		return err
	}

	ctLines, ldpdLines, err := splitConfigDetailSections(res)
	if err != nil {
		return err
	}

	sn, firm, battery, pciaddr, err := parseCTLines(ctLines)
	if err != nil {
		return err
	}

	c.SerialNumber = sn
	c.Firmware = firm
	c.Battery = battery
	c.PCIAddr = pciaddr

	arrays, unassigned, err := splitArrays(ldpdLines)
	if err != nil {
		return err
	}

	for _, arr := range arrays {
		// An array can have multiple logical drive but actually this situation is rare.
		ldchunks, err := splitLDChunks(arr)
		if err != nil {
			return err
		}

		for _, chunk := range ldchunks {
			ldLines, pdLines, err := splitLDPDSections(chunk)
			if err != nil {
				return err
			}

			ld, err := parseLDLines(ldLines)
			if err != nil {
				return err
			}

			pds, err := parsePDLines(pdLines)
			if err != nil {
				return err
			}

			ld.PhyDrives = append(ld.PhyDrives, pds...)

			c.LogDrives = append(c.LogDrives, ld)
		}
	}

	if len(unassigned) > 0 {
		pds, err := parsePDLines(unassigned)
		if err != nil {
			return err
		}

		c.UnconfDrives = pds
	}

	return nil
}
