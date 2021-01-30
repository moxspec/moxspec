package spc

import "os"

// DiskType represents Divice Protocol (SATA/SAS)
type DiskType string

// PostFunc represents type of post function
type PostFunc func(*os.File, []byte) ([]byte, error)

const (
	// SATADisk represents SATA protocol
	SATADisk DiskType = "SATA"
	// SASDisk represents SAS protocol
	SASDisk DiskType = "SAS"
	// UnknownTypeDisk represents unknown protocol
	UnknownTypeDisk DiskType = "Unknown"
)

// CastDiskType cast from strint to DiskType
func CastDiskType(s string) DiskType {
	switch s {
	case "SATA":
		return SATADisk
	case "SAS":
		return SASDisk
	default:
		return UnknownTypeDisk
	}
}

// Device represents a SCSI/ATA device
type Device struct {
	CurTemp             int16
	MinTemp             int16
	MaxTemp             int16
	SpareSpace          byte
	LifeLeft            byte
	Used                byte
	FormFactor          string
	Rotation            uint16
	ModelNumber         string
	SerialNumber        string
	FirmwareRevision    string
	Transport           string
	LogSectorSize       uint32
	PowerCycleCount     uint64
	PowerOnHours        uint64
	UnsafeShutdownCount uint64
	TotalLBAWritten     uint64
	TotalLBARead        uint64
	SigSpeed            string
	NegSpeed            string
	ErrorRecords        []*SmartRecord
	SelfTestSupport     bool
	ErrorLoggingSupport bool

	ioctlDeviceFilePath string
	post                PostFunc
	diskType            DiskType
}

// NewDevice make new SCSI/ATA device
func NewDevice(post PostFunc, ioctlDeviceFilePath string, diskType DiskType) *Device {
	d := new(Device)
	d.post = post
	d.ioctlDeviceFilePath = ioctlDeviceFilePath

	if diskType != SASDisk && diskType != SATADisk {
		log.Debugf("unsupported disk type (%s)", diskType)
		return nil
	}
	d.diskType = diskType

	return d
}
