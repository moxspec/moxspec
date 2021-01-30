package megacli

import "github.com/actapio/moxspec/raidcli"

// Controller represents a megacli controller
type Controller struct {
	raidcli.ControllerSpec
	LogDrives    []*LogDrive
	logDeriveMap map[uint]*LogDrive
	PTPhyDrives  []*PhyDrive // Pass-Through, JBOD
	UnconfDrives []*PhyDrive
}

// GetLogDriveByTarget returns LogDrive which has specified target id
func (c Controller) GetLogDriveByTarget(tgt uint16) *LogDrive {
	if ld, ok := c.logDeriveMap[uint(tgt)]; ok {
		return ld
	}
	return nil
}

// GetPTPhyDriveByWWN returns PhyDrive which has specified wwn
func (c Controller) GetPTPhyDriveByWWN(wwn string) *PhyDrive {
	for _, pd := range c.PTPhyDrives {
		if pd.WWN == wwn {
			return pd
		}
	}
	return nil
}

// LogDrive represents a logical drive
type LogDrive struct {
	raidcli.LogDriveSpec
	GroupID   uint
	TargetID  uint
	PhyDrives []*PhyDrive
}

// IsHealthy returns whether a logical drive is healthy
func (l *LogDrive) IsHealthy() bool {
	return isHealthy(l.LogDriveSpec.State)
}

func isHealthy(stat string) bool {
	if stat == "Optimal" {
		return true
	}
	return false
}

// PhyDrive represents a physical drive
type PhyDrive struct {
	WWN              string
	EnclosureID      string
	SlotNumber       string
	Group            uint16
	Span             uint16
	Arm              uint16
	Type             string // SAS, SATA
	Model            string
	InquiryRaw       string
	Size             uint64
	FirmwareRevision string
	State            string
	PhyBlockSize     uint16
	LogBlockSize     uint16
	ConnectedPort    string
	DriveSpeed       string // drive speed (assume it negotiated speed)
	LinkSpeed        string // port speed
	CurTemp          int16
	SolidStateDrive  bool
	SMARTAlert       bool
	MediaErrorCount  uint16
	DeviceID         uint16
}
