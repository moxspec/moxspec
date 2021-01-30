package sas3ircu

import (
	"strings"

	"github.com/actapio/moxspec/raidcli"
)

// Controller represents a mpt controller
type Controller struct {
	raidcli.ControllerSpec
	LogDrives   []*LogDrive
	ldAddrMap   map[string]*LogDrive // key = sas_address
	PTPhyDrives []*PhyDrive          // Pass-Through
	ptAddrMap   map[string]*PhyDrive // key = sas_address
}

// GetLogDrive returns a logical drive which has given sas address
func (c Controller) GetLogDrive(sasAddr string) *LogDrive {
	if c.ldAddrMap == nil {
		return nil
	}

	if ld, ok := c.ldAddrMap[sasAddr]; ok {
		return ld
	}

	return nil
}

// GetPTDrive returns a physical drive which has given sas address
func (c Controller) GetPTDrive(sasAddr string) *PhyDrive {
	if c.ptAddrMap == nil {
		return nil
	}

	if pt, ok := c.ptAddrMap[sasAddr]; ok {
		return pt
	}

	return nil
}

// LogDrive represents a logical drive
type LogDrive struct {
	raidcli.LogDriveSpec
	VolumeID  int
	WWID      string // kernel exposes wwid as sas_address
	PhyDrives []*PhyDrive
	pdIDList  []string
}

// IsHealthy returns whether a logical drive is healthy
func (l *LogDrive) IsHealthy() bool {
	return isHealthy(l.LogDriveSpec.State)
}

func isHealthy(stat string) bool {
	if strings.HasSuffix(stat, "(OKY)") {
		return true
	}
	return false
}

// PhyDrive represents a physical drive
type PhyDrive struct {
	EnclosureID     string
	SlotNumber      string
	SASAddress      string
	Protocol        string // SAS, SATA
	Model           string
	SerialNumber    string
	Size            uint64
	Firmware        string
	State           string
	DriveType       string
	SolidStateDrive bool
}
