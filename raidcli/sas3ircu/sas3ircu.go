package sas3ircu

import (
	"fmt"
	"strings"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/raidcli"
	"github.com/moxspec/moxspec/util"
)

var pathList = []string{
	"/usr/sbin/sas3ircu",
}

var clipath string

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("sas3ircu")
	clipath, _ = util.ScanPathList(pathList)
}

// Available returns whether raid command is available
func Available() bool {
	if clipath == "" {
		return false
	}
	return true
}

// GetControllers retruns sas3ircu controllers
func GetControllers() ([]*Controller, error) {
	res, err := raidcli.Run(clipath, "LIST")
	if err != nil {
		return nil, err
	}
	return parseCtlList(res)
}

// Decode makes Adapter satisfy the mox.Decoder interface
func (c *Controller) Decode() error {
	var err error

	res, err := raidcli.Run(clipath, fmt.Sprintf("%d", c.Number), "DISPLAY")
	if err != nil {
		return err
	}

	ctLines, ldLines, pdLines, err := splitSections(res)
	if err != nil {
		return err
	}

	firm, bios, err := parseCTLines(ctLines)
	if err != nil {
		return err
	}
	c.Firmware = firm
	c.BIOS = bios

	lds, err := parseLDLines(ldLines)
	if err != nil {
		return err
	}

	pds, err := parsePDLines(pdLines)
	if err != nil {
		return err
	}

	pdmap, err := genPDMap(pds)
	if err != nil {
		return err
	}

	c.ldAddrMap = make(map[string]*LogDrive)

	log.Debugf("linking ld and pd")
	for _, ld := range lds {
		log.Debugf("ld: %d", ld.VolumeID)
		c.ldAddrMap["0x"+ld.WWID] = ld // linux kernel exposes sas address with prefix "0x"

		for _, pdid := range ld.pdIDList {
			log.Debugf("scanning pd(%s)", pdid)

			if pd, ok := pdmap[pdid]; ok {
				log.Debug("found")
				ld.PhyDrives = append(ld.PhyDrives, pd)
				delete(pdmap, pdid)
			}
		}
	}

	// PhyDrives which are not linked with LogDrive are pass-through drives
	c.ptAddrMap = make(map[string]*PhyDrive)
	for _, ptpd := range pdmap {
		c.PTPhyDrives = append(c.PTPhyDrives, ptpd)

		// linux kernel exposes sas address without "-" and prefixed with "0x"
		key := "0x" + strings.Replace(ptpd.SASAddress, "-", "", -1)
		c.ptAddrMap[key] = ptpd
	}

	c.LogDrives = lds
	return nil
}
