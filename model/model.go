package model

import (
	"fmt"
)

// FirmwareType is used to indicate firmware type
type FirmwareType string

// FirmwareTypes
const (
	BIOS FirmwareType = "BIOS"
	UEFI FirmwareType = "UEFI"
)

// Report represents actual data
type Report struct {
	System      *System            `json:"system,omitempty"`
	Chassis     *Chassis           `json:"chassis,omitempty"`
	Firmware    *Firmware          `json:"firmware,omitempty"`
	Baseboard   *Baseboard         `json:"baseboard,omitempty"`
	Processor   *ProcessorReport   `json:"processor,omitempty"`
	Memory      *MemoryReport      `json:"memory,omitempty"`
	Storage     *StorageReport     `json:"storage,omitempty"`
	Network     *NetworkReport     `json:"network,omitempty"`
	Accelerator *AcceleratorReport `json:"accelerator,omitempty"`
	PCIDevice   []*PCIBaseSpec     `json:"pciDevices,omitempty"`
	PowerSupply []*PowerSupply     `json:"powerSupply,omitempty"`
	BMC         *BMC               `json:"bmc,omitempty"`
	SAR         map[string][]SAR   `json:"sar,omitempty"`
	OS          *OS                `json:"os,omitempty"`
	Hostname    string             `json:"hostname,omitempty"`
	Version     string             `json:"version"`
	Timestamp   int64              `json:"timestamp"`
	Datetime    string             `json:"datetime"`
}

// System represents a product
type System struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
}

// Summary returns summarized string
func (s System) Summary() string {
	return fmt.Sprintf("%s %s (SN:%s)", s.Manufacturer, s.ProductName, s.SerialNumber)
}

// Chassis represents a chassis
type Chassis struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
}

// Firmware represents a system firmware
type Firmware struct {
	Type        FirmwareType `json:"type,omitempty"`
	Vendor      string       `json:"vendor,omitempty"`
	Version     string       `json:"version,omitempty"`
	ReleaseDate string       `json:"releaseDate,omitempty"`
}

// Summary returns summarized string
func (f Firmware) Summary() string {
	return fmt.Sprintf("%s, ver %s, release %s", f.Vendor, f.Version, f.ReleaseDate)
}

// Baseboard represents a baseboard
type Baseboard struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
}

// Summary returns summarized string
func (b Baseboard) Summary() string {
	return fmt.Sprintf("%s %s (SN:%s)", b.Manufacturer, b.ProductName, b.SerialNumber)
}

// PowerSupply represents a power supply
type PowerSupply struct {
	Manufacturer    string `json:"manufacturer,omitempty"`
	ProductName     string `json:"productName,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	ModelPartNumber string `json:"modelPartNumber,omitempty"`
	Capacity        uint16 `json:"capacity,omitempty"`
	Present         bool   `json:"present"`
	Plugged         bool   `json:"plugged"`
	HotSwappable    bool   `json:"hotSwappable"`
}

// Summary returns summarized string
func (p PowerSupply) Summary() string {
	sum := fmt.Sprintf("%s %s", p.Manufacturer, p.ProductName)

	if p.Capacity > 0 {
		sum = fmt.Sprintf("%s (%dW)", sum, p.Capacity)
	}
	if p.SerialNumber != "" {
		sum = fmt.Sprintf("%s (SN:%s)", sum, p.SerialNumber)
	}

	return sum
}

// Stat returns status flags as string
func (p PowerSupply) Stat() string {
	return fmt.Sprintf("present: %t, plugged: %t, hotswap: %t", p.Present, p.Plugged, p.HotSwappable)
}

// BMC represents a baseboard management controller
type BMC struct {
	Type     string `json:"type,omitempty"`
	Firmware string `json:"firmware,omitempty"`
	MAC      string `json:"hwaddr,omitempty"`
	IPAddr   string `json:"ipaddr,omitempty"`
	Netmask  string `json:"netmask,omitempty"`
	MaskSize int    `json:"masksize,omitempty"`
	Gateway  string `json:"gateway,omitempty"`
}

// OS represents an operating system
type OS struct {
	Distro string `json:"distro,omitempty"`
	Kernel string `json:"kernel,omitempty"`
}

// SAR represents an sar results
type SAR struct {
	Time   string   `json:"time,omitempty"`
	Dev    string   `json:"dev,omitempty"`
	Values []string `json:"values,omitempty"`
}
