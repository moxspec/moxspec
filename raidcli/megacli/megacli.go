package megacli

import (
	"fmt"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/raidcli"
	"github.com/actapio/moxspec/util"
)

var pathList = []string{
	"/opt/MegaRAID/MegaCli/MegaCli64",
}

var clipath string

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("megacli")
	clipath, _ = util.ScanPathList(pathList)
}

// Available returns whether raid command is available
func Available() bool {
	if clipath == "" {
		return false
	}
	return true
}

// GetControllers retruns megaraid controllers
func GetControllers() ([]*Controller, error) {
	res, err := raidcli.Run(clipath, "-AdpGetPciInfo", "-aAll", "-NoLog")
	if err != nil {
		return nil, err
	}
	return parsePCIInfo(res)
}

// Decode makes Adapter satisfy the mox.Decoder interface
func (c *Controller) Decode() error {
	err := setAdpInfo(c)
	if err != nil {
		return err
	}

	lds, err := getLDList(c.Number)
	if err != nil {
		return err
	}

	ldmap := make(map[uint]*LogDrive)
	for _, ld := range lds {
		if _, ok := ldmap[ld.GroupID]; ok {
			return fmt.Errorf("duplicate group id found: %d", ld.GroupID)
		}

		ldmap[ld.GroupID] = ld
	}

	c.LogDrives = lds
	c.logDeriveMap = ldmap

	// scan unconfigured drives
	apd, err := getAllPD(c.Number)
	if err != nil {
		return err
	}
	if len(apd) == 0 {
		return nil
	}

	wwnmap := make(map[string]bool)
	inqmap := make(map[string]bool)
	for _, ld := range c.LogDrives {
		for _, pd := range ld.PhyDrives {
			if pd.WWN != "" {
				wwnmap[pd.WWN] = true
			}

			if pd.InquiryRaw != "" {
				inqmap[pd.InquiryRaw] = true
			}
		}
	}

	for _, pd := range apd {
		if wwnmap[pd.WWN] { // means 'this drive is under a logical drive'
			continue
		}

		// means 'this drive is under a logical drive'
		// same SSD does not return WWN
		if inqmap[pd.InquiryRaw] {
			continue
		}

		if pd.State == "JBOD" {
			c.PTPhyDrives = append(c.PTPhyDrives, pd)
		} else {
			c.UnconfDrives = append(c.UnconfDrives, pd)
		}
	}

	return nil
}
