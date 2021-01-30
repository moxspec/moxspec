package raidcli

// ControllerSpec represents common controller spec
type ControllerSpec struct {
	Number       int
	Domain       uint32
	Bus          uint32
	Device       uint32
	Function     uint32
	ProductName  string
	Firmware     string
	BIOS         string
	SerialNumber string
	Battery      bool
}

// LogDriveSpec represents common logical drive spec
type LogDriveSpec struct {
	Label       string // controller specific identidier
	RAIDLv      Level
	Size        uint64
	State       string
	StripSize   uint64
	CachePolicy string
}

// HealthReporter returns whether it is healthy
type HealthReporter interface {
	IsHealthy() bool
}
