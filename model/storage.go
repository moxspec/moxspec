package model

import (
	"fmt"
	"strings"

	"github.com/moxspec/moxspec/util"
)

// batteryStatus represents battery unit presence
type batteryStatus string

// These are the battery presence status
const (
	BatteryPresent    batteryStatus = "present"
	BatteryNotPresent batteryStatus = "not present"
	BatteryUnknown    batteryStatus = "unknown"
)

// Health Check Threshold
const (
	ErrorCountThreshold = 5
)

// StorageReport represents a storage report
type StorageReport struct {
	NVMeControllers   []*NVMeController   `json:"nvmeControllers,omitempty"`
	RAIDControllers   []*RAIDController   `json:"raidControllers,omitempty"`
	AHCIControllers   []*AHCIController   `json:"ahciControllers,omitempty"`
	VirtControllers   []*VirtController   `json:"virtControllers,omitempty"`
	NonStdControllers []*NonStdController `json:"nonStdControllers,omitempty"`
}

// StorageSizeSpec represents a storage size spec
type StorageSizeSpec struct {
	Size uint64 `json:"size,omitempty"` // should be bytes
}

// SizeString returns the string representation of the size of the device
func (s StorageSizeSpec) SizeString() string {
	u, _ := util.ConvUnitDecFit(s.Size, util.GIGA)
	return u
}

// StorageIOSpec represents an io statistics spec
type StorageIOSpec struct {
	ByteRead    uint64 `json:"byteRead,omitempty"`
	ByteWritten uint64 `json:"byteWritten,omitempty"`
}

// ByteReadString returns read bytes in string
func (s StorageIOSpec) ByteReadString() string {
	str, err := util.ConvUnitBinFit(s.ByteRead, util.KILO)
	if err != nil {
		return ""
	}
	return str
}

// ByteWrittenString returns written bytes in string
func (s StorageIOSpec) ByteWrittenString() string {
	str, err := util.ConvUnitBinFit(s.ByteWritten, util.KILO)
	if err != nil {
		return ""
	}
	return str
}

// IORatio returns read/write ratio
func (s StorageIOSpec) IORatio() (w float64, r float64) {
	total := float64(s.ByteRead) + float64(s.ByteWritten)
	if total == 0 {
		return 0.0, 0.0
	}

	w = float64(s.ByteWritten) / total * 100.0
	r = float64(s.ByteRead) / total * 100.0
	return
}

// IOSummary returns summarized io accounting
func (s StorageIOSpec) IOSummary() string {
	wr, rr := s.IORatio()
	return fmt.Sprintf("written %s, read %s (w:%.1f%%/r:%.1f%%)", s.ByteWrittenString(), s.ByteReadString(), wr, rr)
}

// StorageTempSpec represents a temperature spec
type StorageTempSpec struct {
	CurTemp  int16 `json:"curTemp,omitempty"`
	MaxTemp  int16 `json:"maxTemp,omitempty"`
	MinTemp  int16 `json:"minTemp,omitempty"`
	WarnTemp int16 `json:"warnTemp,omitempty"`
	CritTemp int16 `json:"critTemp,omitempty"`
}

// TempWarnCritSummary returns warn/crit style temp spec summary
func (s StorageTempSpec) TempWarnCritSummary() string {
	return fmt.Sprintf("cur %d°C, warn %d°C, crit %d°C", s.CurTemp, s.WarnTemp, s.CritTemp)
}

// TempMaxMinSummary returns max/min style temp spec summary
func (s StorageTempSpec) TempMaxMinSummary() string {
	return fmt.Sprintf("cur %d°C, max %d°C, min %d°C", s.CurTemp, s.MaxTemp, s.MinTemp)
}

// TempSummary returns simple temp summary
func (s StorageTempSpec) TempSummary() string {
	return fmt.Sprintf("%d°C", s.CurTemp)
}

// StoragePowerStatSpec represents a power statistics spec
type StoragePowerStatSpec struct {
	PowerOnHours        uint64 `json:"powerOnHours,omitempty"`
	PowerCycleCount     uint64 `json:"powerCycleCount,omitempty"`
	UnsafeShutdownCount uint64 `json:"unsafeShutdownCount,omitempty"`
}

// SMARTDiagSpec represents a diagnostic spec
type SMARTDiagSpec struct {
	ErrorRecords []*SMARTRecord `json:"errorRecords,omitempty"`
}

// IsHealthy returns whether a disk is healthy
func (s SMARTDiagSpec) IsHealthy() bool {
	return (len(s.ErrorRecords) == 0)
}

// DiagSummaries returns summarized strings
func (s SMARTDiagSpec) DiagSummaries() []string {
	var m []string
	for _, e := range s.ErrorRecords {
		m = append(m, e.String())
	}
	return m
}

// SMARTRecord represents a s.m.a.r.t error counter
type SMARTRecord struct {
	ID        byte   `json:"id"`
	Current   byte   `json:"current"`
	Worst     byte   `json:"worst"`
	Raw       int64  `json:"raw"`
	Threshold byte   `json:"threshold"`
	Name      string `json:"name"`
}

// String makes SMARTRecord satisfy the Stringer interface
func (s SMARTRecord) String() string {
	return fmt.Sprintf("%03d.%s = current:%d worst:%d threshold:%d raw:%d", s.ID, s.Name, s.Current, s.Worst, s.Threshold, s.Raw)
}

// SCSIAddressSpec represents SCSI address
type SCSIAddressSpec struct {
	SCSIHost    uint16 `json:"scsiHost"`
	SCSIChannel uint16 `json:"scsiChannel"`
	SCSITarget  uint16 `json:"scsiTarget"`
	SCSILun     uint16 `json:"scsiLun"`
}

// SCSIAddress returns a scsi address
func (s SCSIAddressSpec) SCSIAddress() string {
	return fmt.Sprintf("%d:%d:%d:%d", s.SCSIHost, s.SCSIChannel, s.SCSITarget, s.SCSILun)
}

// AHCIController represents an ahci storage controller
type AHCIController struct {
	PCIBaseSpec
	Drives []*Drive `json:"drives,omitempty"`
}

// Summary returns summarized string
func (a AHCIController) Summary() string {
	if a.Driver == "" {
		return fmt.Sprintf("%s (node%d)", a.LongName(), a.Numa)
	}
	return fmt.Sprintf("%s (%s) (node%d)", a.LongName(), a.Driver, a.Numa)
}

// VirtController represents a virtio storage controller
type VirtController struct {
	PCIBaseSpec
	Drives []*Drive `json:"drives,omitempty"`
}

// Summary returns summarized string
func (v VirtController) Summary() string {
	if v.Driver == "" {
		return fmt.Sprintf("%s (node%d)", v.LongName(), v.Numa)
	}
	return fmt.Sprintf("%s (%s) (node%d)", v.LongName(), v.Driver, v.Numa)
}

// Drive represents a typical storage drive
type Drive struct {
	ID           string `json:"id,omitempty"`   // maj:min
	Name         string `json:"name,omitempty"` // kernel name, eg. sda
	Model        string `json:"model,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	Firmware     string `json:"firmware,omitempty"`
	Transport    string `json:"transport,omitempty"`
	Blocks       uint64 `json:"blocks,omitempty"`
	PhyBlockSize uint16 `json:"phyBlockSize,omitempty"`
	LogBlockSize uint16 `json:"logBlockSize,omitempty"`
	Scheduler    string `json:"scheduler,omitempty"`
	FormFactor   string `json:"formFactor,omitempty"`
	Driver       string `json:"driver,omitempty"`
	SigSpeed     string `json:"sigSpeed,omitempty"`
	NegSpeed     string `json:"negSpeed,omitempty"`
	Rotation     uint16 `json:"rotation,omitempty"`
	SelfTest     bool   `json:"selfTest,omitempty"`
	ErrorLogging bool   `json:"errorLogging,omitempty"`
	StorageIOSpec
	StorageTempSpec
	StorageSizeSpec
	StoragePowerStatSpec
	SMARTDiagSpec
	SCSIAddressSpec
}

// Summary returns summarized string
func (d Drive) Summary() string {
	fmtr := "%s %s %s (%s) (log:%dB/phy:%dB) (sched:%s)"
	s := fmt.Sprintf(fmtr, d.Name, d.Model, d.SizeString(), d.Driver, d.LogBlockSize, d.PhyBlockSize, d.Scheduler)
	if d.SerialNumber != "" {
		s = fmt.Sprintf("%s (SN:%s)", s, d.SerialNumber)
	}
	return s
}

// IsSSD returns whether a device is ssd
func (d Drive) IsSSD() bool {
	return (d.Rotation == 1)
}

// IsHDD returns whether a device is hdd
func (d Drive) IsHDD() bool {
	return (d.Rotation > 1)
}

// FormSummary returns form factor summary
func (d Drive) FormSummary() string {
	var sum string

	if d.IsSSD() {
		sum = fmt.Sprintf("SSD (%s)", d.Transport)
	}
	if d.IsHDD() {
		sum = fmt.Sprintf("HDD (%drpm) (%s)", d.Rotation, d.Transport)
	}

	if d.FormFactor != "" {
		sum = fmt.Sprintf("%s %s", d.FormFactor, sum)
	}

	if d.IsSSD() || d.IsHDD() {
		if d.SigSpeed == "" || d.NegSpeed == "" {
			return sum
		}
		return fmt.Sprintf("%s (cur:%s, max:%s)", sum, d.NegSpeed, d.SigSpeed)
	}

	return d.FormFactor
}

// NVMeController represents a NVMe controller
type NVMeController struct {
	PCIBaseSpec
	Name       string       `json:"name,omitempty"`
	Model      string       `json:"model,omitempty"`
	Firmware   string       `json:"firmware,omitempty"`
	Namespaces []*Namespace `json:"namespaces,omitempty"`
	StorageIOSpec
	StorageTempSpec
	StorageSizeSpec
	StoragePowerStatSpec
}

// LongName returns pretty name
func (n NVMeController) LongName() string {
	if n.VendorName != "" && n.Model != "" {
		return fmt.Sprintf("%s %s", n.VendorName, n.Model)
	}

	return n.PCIBaseSpec.LongName()
}

// Summary returns summarized string
func (n NVMeController) Summary() string {
	var sum string
	if n.Driver == "" {
		sum = fmt.Sprintf("%s (node%d)", n.LongName(), n.Numa)
	} else {
		sum = fmt.Sprintf("%s (%s) (node%d)", n.LongName(), n.Driver, n.Numa)
	}
	if n.SerialNumber != "" {
		sum = fmt.Sprintf("%s (SN:%s)", sum, n.SerialNumber)
	}
	return sum
}

// Namespace represents a NVMe namespace
type Namespace struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Scheduler    string `json:"scheduler,omitempty"`
	PhyBlockSize uint16 `json:"phyBlockSize,omitempty"`
	LogBlockSize uint16 `json:"logBlockSize,omitempty"`
	StorageSizeSpec
}

// Summary returns summarized string
func (n Namespace) Summary() string {
	return fmt.Sprintf("%s %s (log:%dB/phy:%dB) (sched:%s)", n.Name, n.SizeString(), n.LogBlockSize, n.PhyBlockSize, n.Scheduler)
}

// RAIDController represents a RAID controller
// NOTE:
//   pci.ids reports a chipset name not a raid product name.
//   mox tries to get a product name from raid cli outputs and set it if it is available.
type RAIDController struct {
	PCIBaseSpec
	ProductName       string        `json:"productName,omitempty"`
	SerialNumber      string        `json:"serialNumber,omitempty"`
	AdapterID         string        `json:"adapterId,omitempty"`
	Firmware          string        `json:"firmware,omitempty"`
	BIOS              string        `json:"bios,omitempty"`
	Battery           batteryStatus `json:"battery"`
	LogDrives         []*LogDrive   `json:"logDrives,omitempty"`
	UnconfDrives      []*PhyDrive   `json:"unconfDrives,omitempty"`
	PassthroughDrives []*PhyDrive   `json:"passthroughDrives,omitempty"`
}

// Summary returns summarized string
func (r RAIDController) Summary() string {
	var sum string

	sum = r.LongName()

	if r.ProductName != "" {
		sum = fmt.Sprintf("%s, %s", sum, r.ProductName)
	}

	if r.Driver != "" {
		sum = fmt.Sprintf("%s (%s)", sum, r.Driver)
	}

	sum = fmt.Sprintf("%s (node%d)", sum, r.Numa)

	if r.SerialNumber == "" {
		return sum
	}

	return fmt.Sprintf("%s (SN:%s)", sum, r.SerialNumber)
}

// SpecSummary returns summarized string
func (r RAIDController) SpecSummary() string {
	var sum []string

	if r.Firmware != "" {
		sum = append(sum, fmt.Sprintf("firm: %s", r.Firmware))
	}

	if r.BIOS != "" {
		sum = append(sum, fmt.Sprintf("bios: %s", r.Firmware))
	}

	sum = append(sum, fmt.Sprintf("battery: %s", r.Battery))

	return strings.Join(sum, ", ")
}

// LogDrive represents a logical drive under a raid card
// if it is marked as pass-through drive it will be treated as normal drive
type LogDrive struct {
	Drive

	// RAID specific attributes
	GroupLabel  string      `json:"groupLabel,omitempty"`
	RAIDLv      string      `json:"raidLv,omitempty"`
	StripeSize  uint64      `json:"stripeSize,omitempty"`
	Status      string      `json:"status,omitempty"`
	CachePolicy string      `json:"cachePolicy,omitempty"`
	WWN         string      `json:"wwn,omitempty"`
	SASAddress  string      `json:"sasAddress"`
	Degraded    bool        `json:"degraded"`
	PhyDrives   []*PhyDrive `json:"phyDrives,omitempty"`
}

// LDSummary returns summarized string
func (l LogDrive) LDSummary() string {
	var cfgs []string

	if l.RAIDLv != "" {
		cfgs = append(cfgs, l.RAIDLv)

		if l.CachePolicy != "" {
			cfgs = append(cfgs, l.CachePolicy)
		}

		if l.StripeSize > 0 {
			ssize, err := util.ConvUnitBinFit(l.StripeSize, util.KILO)
			if err == nil {
				cfgs = append(cfgs, fmt.Sprintf("stripe: %s", ssize))
			}
		}
	}

	fmtr := "%s %s (%s) (log:%dB/phy:%dB) (sched:%s)"
	sum := fmt.Sprintf(fmtr, l.Name, l.SizeString(), l.Driver, l.LogBlockSize, l.PhyBlockSize, l.Scheduler)

	if l.GroupLabel != "" {
		sum = fmt.Sprintf("[%s] %s", l.GroupLabel, sum)
	}

	if len(cfgs) == 0 {
		return sum
	}

	return fmt.Sprintf("%s (%s)", sum, strings.Join(cfgs, ", "))
}

// IsHealthy returns whether a logical disk is healthy
func (l LogDrive) IsHealthy() bool {
	for _, pd := range l.PhyDrives {
		if !pd.IsHealthy() {
			return false
		}
	}

	return !l.Degraded
}

// PhyDrive represents a physical drive under a raid card
type PhyDrive struct {
	Drive

	// RAID PD specific attributes
	Enclosure       string `json:"enclosure,omitempty"`
	Slot            string `json:"slot,omitempty"`
	Status          string `json:"status,omitempty"`
	SolidStateDrive bool   `json:"ssd,omitempty"`
	ErrorCount      uint16 `json:"errorCount,omitempty"`

	// Pass-Through specific attributes
	SASAddress string `json:"sasAddress"`
	WWN        string `json:"wwn,omitempty"`
}

// Pos returns drive position string
func (p PhyDrive) Pos() string {
	return fmt.Sprintf("%s:%s", p.Enclosure, p.Slot)
}

// PTSummary returns summarized string
func (p PhyDrive) PTSummary() string {
	if p.ErrorCount > ErrorCountThreshold {
		fmtr := "[%s:%s] %s %s %s (%s) (log:%dB/phy:%dB) (sched:%s), %s, err=%d"
		return fmt.Sprintf(fmtr, p.Enclosure, p.Slot, p.Name, p.Model, p.SizeString(), p.Driver, p.LogBlockSize, p.PhyBlockSize, p.Scheduler, p.Status, p.ErrorCount)
	}

	fmtr := "[%s:%s] %s %s %s (%s) (log:%dB/phy:%dB) (sched:%s), %s"
	return fmt.Sprintf(fmtr, p.Enclosure, p.Slot, p.Name, p.Model, p.SizeString(), p.Driver, p.LogBlockSize, p.PhyBlockSize, p.Scheduler, p.Status)
}

// Summary returns summarized string
func (p PhyDrive) Summary() string {
	if p.ErrorCount > ErrorCountThreshold {
		return fmt.Sprintf("[%s:%s] %s %s, %s, %s, err=%d", p.Enclosure, p.Slot, p.SizeString(), p.FormSummary(), p.Model, p.Status, p.ErrorCount)
	}

	return fmt.Sprintf("[%s:%s] %s %s, %s, %s", p.Enclosure, p.Slot, p.SizeString(), p.FormSummary(), p.Model, p.Status)
}

// FormSummary returns form factor summary
func (p PhyDrive) FormSummary() string {
	var sum string
	if p.SolidStateDrive {
		sum = fmt.Sprintf("%s SSD", p.Transport)
	} else {
		sum = fmt.Sprintf("%s HDD", p.Transport)
	}

	if p.NegSpeed == "" {
		return sum
	}

	return fmt.Sprintf("%s (%s)", sum, p.NegSpeed)
}

// IsHealthy returns whether a physical disk is healthy
func (p PhyDrive) IsHealthy() bool {
	return p.ErrorCount == 0
}

// NonStdController represents an non-standard storage controller
type NonStdController struct {
	PCIBaseSpec
	Name         string         `json:"name,omitempty"` // kernel name, eg. fct0
	Model        string         `json:"model,omitempty"`
	SerialNumber string         `json:"serialNumber,omitempty"`
	Firmware     string         `json:"firmware,omitempty"`
	ErrorRecords []string       `json:"errors,omitempty"`
	Drives       []*NonStdDrive `json:"drives,omitempty"`
	StorageIOSpec
	StorageTempSpec
	StoragePowerStatSpec
}

// Summary returns summarized string
func (n NonStdController) Summary() string {
	// this should never be called because non standard means that a device requires a driver
	if n.Driver == "" {
		return fmt.Sprintf("%s (node%d)", n.LongName(), n.Numa)
	}
	return fmt.Sprintf("%s (%s) (node%d)", n.LongName(), n.Driver, n.Numa)
}

// HasErrorRecords returns whether a device is healthy
func (n NonStdController) HasErrorRecords() bool {
	return (len(n.ErrorRecords) != 0)
}

// NonStdDrive represents a non-standard storage drive
type NonStdDrive struct {
	ID           string `json:"id,omitempty"`   // maj:min
	Name         string `json:"name,omitempty"` // kernel name
	Blocks       uint64 `json:"blocks,omitempty"`
	PhyBlockSize uint16 `json:"phyBlockSize,omitempty"`
	LogBlockSize uint16 `json:"logBlockSize,omitempty"`
	Scheduler    string `json:"scheduler,omitempty"`
	StorageSizeSpec
}

// Summary returns summarized string
func (n NonStdDrive) Summary() string {
	fmtr := "%s %s (log:%dB/phy:%dB) (sched:%s)"
	return fmt.Sprintf(fmtr, n.Name, n.SizeString(), n.LogBlockSize, n.PhyBlockSize, n.Scheduler)
}
