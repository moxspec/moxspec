package model

import (
	"fmt"
	"strings"
)

// NetworkReport represents a network report
type NetworkReport struct {
	EthControllers []*EthController `json:"ethControllers,omitempty"`
}

// EthController represents an ethernet controller
type EthController struct {
	PCIBaseSpec
	Interfaces []*NetInterface `json:"interfaces,omitempty"`
}

// Summary returns summarized string
func (e EthController) Summary() string {
	var sum string
	if e.HasDriver() {
		sum = fmt.Sprintf("%s (%s) (node%d)", e.LongName(), e.Driver, e.Numa)
	} else {
		sum = fmt.Sprintf("%s (node%d)", e.LongName(), e.Numa)
	}
	if e.SerialNumber != "" {
		sum = fmt.Sprintf("%s (SN:%s)", sum, e.SerialNumber)
	}
	return sum
}

// NetInterface represents a network interface
type NetInterface struct {
	State            string       `json:"state,omitempty"`
	Name             string       `json:"name,omitempty"`
	HWAddr           string       `json:"hwaddr,omitempty"`
	Speed            uint32       `json:"speed"`
	MTU              uint32       `json:"mtu"`
	SupportedSpeed   []string     `json:"supportedSpeed,omitempty"`
	AdvertisingSpeed []string     `json:"advertisingSpeed,omitempty"`
	IPAddrs          []*IPAddress `json:"ipaddrs,omitempty"`
	Module           *Module      `json:"module,omitempty"`
	FirmwareVersion  string       `json:"firmwareVersion,omitempty"`
	RxErrors         uint64       `json:"rxErrors"`
	TxErrors         uint64       `json:"txErrors"`
	RxDropped        uint64       `json:"rxDropped"`
	TxDropped        uint64       `json:"txDropped"`
}

// Summary returns summarized string
func (n NetInterface) Summary() string {
	return fmt.Sprintf("%s, %s %s", n.Name, n.HWAddr, n.IPAddrsString())
}

// IPAddrsString returns a concatenated ipaddrs as string
func (n NetInterface) IPAddrsString() string {
	var list []string
	for _, i := range n.IPAddrs {
		list = append(list, fmt.Sprintf("%s/%d", i.Addr, i.MaskSize))
	}
	return strings.Join(list, ", ")
}

// StatSummary returns summarized status string
func (n NetInterface) StatSummary() string {
	return fmt.Sprintf("%s, speed %d, mtu %d", n.State, n.Speed, n.MTU)
}

// IPAddress represents a network address
type IPAddress struct {
	Version   byte   `json:"version,omitempty"`
	Addr      string `json:"addr,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	MaskSize  int    `json:"maskSize,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
	Network   string `json:"network,omitempty"`
}

// Module represents a network module
type Module struct {
	FormFactor   string `json:"formFactor,omimtempty"`
	VendorName   string `json:"vendorName,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	CableLength  uint8  `json:"cableLength,omitempty"`
	Connector    string `json:"connector,omitempty"`
}

// Summary returns summarized string
func (m Module) Summary() string {
	sum := fmt.Sprintf("%s %s (SN: %s), %s", m.VendorName, m.ProductName, m.SerialNumber, m.FormFactor)

	if m.CableLength > 0 {
		sum = fmt.Sprintf("%s %s (%dm)", sum, m.Connector, m.CableLength)
	} else {
		sum = fmt.Sprintf("%s %s", sum, m.Connector)
	}

	return sum
}
