package hpacucli

import (
	"strings"

	"github.com/moxspec/moxspec/raidcli"
)

// Controller represents a mpt controller
type Controller struct {
	raidcli.ControllerSpec
	Slot         string
	SerialNumber string
	PCIAddr      string
	Battery      bool
	LogDrives    []*LogDrive
	PTPhyDrives  []*PhyDrive // Pass-Through
	UnconfDrives []*PhyDrive // Unassigned, Unconfigured
}

// GetLogDrive returns a LogDrive which is associated with given drive name
func (c Controller) GetLogDrive(name string) *LogDrive {
	if !strings.HasPrefix(name, "/dev") {
		return nil
	}

	for _, ld := range c.LogDrives {
		if ld.DiskName == name {
			return ld
		}
	}

	return nil
}

// LogDrive represents a logical drive
type LogDrive struct {
	raidcli.LogDriveSpec
	VolumeID  string
	DiskName  string // /dev/sda
	UUID      string
	StripSize uint64
	pdIDList  []string
	PhyDrives []*PhyDrive
}

// IsHealthy returns whether a logical drive is healthy
func (l *LogDrive) IsHealthy() bool {
	return isHealthy(l.LogDriveSpec.State)
}

func isHealthy(stat string) bool {
	if stat == "OK" {
		return true
	}
	return false
}

// PhyDrive represents a physical drive
type PhyDrive struct {
	Port         string
	Box          string
	Bay          string
	Status       string
	Protocol     string // SAS, SATA
	Size         uint64
	Firmware     string
	SerialNumber string
	Model        string
	CurTemp      int
	MaxTemp      int
	Rotation     uint
	NegSpeed     string
}
