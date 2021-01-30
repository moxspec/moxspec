package model

import (
	"fmt"

	"github.com/actapio/moxspec/pci"
)

// PCIBaseSpec represents a basic pci spec
type PCIBaseSpec struct {
	Location struct {
		Domain   uint32 `json:"domain"`
		Bus      uint32 `json:"bus"`
		Device   uint32 `json:"device"`
		Function uint32 `json:"function"`
	} `json:"location"`
	Path              string    `json:"path,omitempty"`
	VendorID          uint16    `json:"vendorID"`
	VendorName        string    `json:"vendorName,omitempty"`
	DeviceID          uint16    `json:"deviceID"`
	DeviceName        string    `json:"deviceName,omitempty"`
	SubSystemVendorID uint16    `json:"SubSystemVendorID"`
	SubSystemDeviceID uint16    `json:"SubSystemDeviceID"`
	SubSystemName     string    `json:"SubSystemName,omitempty"`
	ClassID           byte      `json:"classID"`
	ClassName         string    `json:"className,omitempty"`
	SubClassID        byte      `json:"subClassID"`
	SubClassName      string    `json:"subClassName,omitempty"`
	InterfaceID       byte      `json:"interfaceID"`
	InterfaceName     string    `json:"interfaceName,omitempty"`
	SerialNumber      string    `json:"serialNumber,omitempty"`
	Driver            string    `json:"driver,omitempty"`
	Numa              byte      `json:"numa"`
	CurLink           *PCIeLink `json:"currentLink,omitempty"`
	MaxLink           *PCIeLink `json:"maxLink,omitempty"`
	PowerLimit        float32   `json:"powerLimit,omitempty"`
	UEList            []string  `json:"ueList,omitempty"`
	CEList            []string  `json:"ceList,omitempty"`
}

// LongName returns pretty name
func (p PCIBaseSpec) LongName() string {
	if p.VendorName != "" && p.DeviceName != "" {
		return fmt.Sprintf("%s %s", p.VendorName, p.DeviceName)
	}

	ln := fmt.Sprintf("%04x:%04x", p.VendorID, p.DeviceID)
	if p.VendorName == "" {
		return ln
	}

	return fmt.Sprintf("%s %s", p.VendorName, ln)
}

// PCIID returns pci location in the system
func (p PCIBaseSpec) PCIID() string {
	return pci.IDString(p.Location.Domain, p.Location.Bus, p.Location.Device, p.Location.Function)
}

// HasLinkStatus returns whether a device has link status
func (p PCIBaseSpec) HasLinkStatus() bool {
	return (p.CurLink != nil && p.MaxLink != nil)
}

// HasPowerStatus returns whether a device has power status
func (p PCIBaseSpec) HasPowerStatus() bool {
	return (p.PowerLimit != 0)
}

// HasDriver returns whether a device has a driver
func (p PCIBaseSpec) HasDriver() bool {
	return (p.Driver != "")
}

// LinkSummary returns summarized link status
func (p PCIBaseSpec) LinkSummary() string {
	return fmt.Sprintf("Gen%d %.1fGT/s x%d (max: Gen%d %.1fGT/s x%d)",
		p.CurLink.Gen, p.CurLink.Speed, p.CurLink.Width,
		p.MaxLink.Gen, p.MaxLink.Speed, p.MaxLink.Width,
	)
}

// IsHealthy returns whether a device is healthy
func (p PCIBaseSpec) IsHealthy() bool {
	return (len(p.UEList) == 0 && len(p.CEList) == 0)
}

// DiagSummaries returns diag status
func (p PCIBaseSpec) DiagSummaries() []string {
	var sum []string
	for _, uc := range p.UEList {
		sum = append(sum, fmt.Sprintf("[ue] %s", uc))
	}
	for _, ce := range p.CEList {
		sum = append(sum, fmt.Sprintf("[ce] %s", ce))
	}
	return sum
}

// PCIeLink represents a PCIeLink spec
type PCIeLink struct {
	Gen   byte    `json:"gen,omitempty"`
	Speed float32 `json:"speed,omitempty"`
	Width byte    `json:"width,omitempty"`
}
