package model

import (
	"fmt"

	"github.com/moxspec/moxspec/util"
)

// Helth Check Threshold
const (
	CECountThreshold = 1000
)

// MemorySizeSpec represents a memory size spec
type MemorySizeSpec struct {
	Size uint64 `json:"size"`
}

// SizeString returns the string representation of the size of the device
func (m MemorySizeSpec) SizeString() string {
	u, _ := util.ConvUnitDecFit(m.Size, util.GIGA)
	return u
}

// MemoryReport represents a memory report
type MemoryReport struct {
	Total       uint64              `json:"total,omitempty"`
	Empty       byte                `json:"empty,omitempty"`
	Modules     []*MemoryModule     `json:"modules,omitempty"`
	Controllers []*MemoryController `json:"controllers,omitempty"`
}

// TotalString returns total size string in GB
func (m MemoryReport) TotalString() string {
	u, _ := util.ConvUnitDecFit(m.Total, util.GIGA)
	return u
}

// ModuleSummaries returns memory module summaries
func (m MemoryReport) ModuleSummaries() []string {
	var devs = make(map[string]int)
	for _, md := range m.Modules {
		if md.Size == 0 {
			continue
		}
		devs[md.Summary()]++
	}

	var res []string
	for k, v := range devs {
		res = append(res, fmt.Sprintf("%d x %s", v, k))
	}
	res = append(res, fmt.Sprintf("%d x empty", m.Empty))
	return res
}

// HasDiag returns whether memory has diag status
func (m MemoryReport) HasDiag() bool {
	return (len(m.Controllers) != 0)
}

// IsHealthy returns whether memory is healthy
func (m MemoryReport) IsHealthy() bool {
	for _, ctl := range m.Controllers {
		for _, cs := range ctl.CSRows {
			if cs.HasError() {
				return false
			}
		}
	}
	return true
}

// DiagSummaries returns diag summaries
func (m MemoryReport) DiagSummaries() []string {
	var res []string
	for _, ctl := range m.Controllers {
		if ctl.HasError() {
			res = append(res, ctl.Summaries()...)
		}
	}
	return res
}

// MemoryModule represents a memory module
type MemoryModule struct {
	Locator         string  `json:"locator,omitempty"`
	Manufacturer    string  `json:"manufacturer,omitempty"`
	PartNumber      string  `json:"partNumber,omitempty"`
	SerialNumber    string  `json:"serialNumber,omitempty"`
	FormFactor      string  `json:"formFactor,omitempty"`
	Type            string  `json:"type,omitempty"`
	TypeDetail      string  `json:"typeDetail,omitempty"`
	Speed           uint16  `json:"speed,omitempty"`
	ConfiguredSpeed uint16  `json:"configuredSpeed,omitempty"`
	Voltage         float32 `json:"voltage,omitempty"`
	IsPersistent    bool    `json:"isPersistent,omitempty"`
	MemorySizeSpec
}

// Spec returns module spec string
func (m MemoryModule) Spec() string {
	summary := fmt.Sprintf("%s-%d %s %s", m.Type, m.Speed, m.FormFactor, m.SizeString())
	if m.IsPersistent {
		summary = fmt.Sprintf("%s (Persistent)", summary)
	}

	return summary
}

// Summary returns summarized string
func (m MemoryModule) Summary() string {
	return fmt.Sprintf("%s %s", m.Manufacturer, m.Spec())
}

// MemoryController represents a memory controller
type MemoryController struct {
	Name          string           `json:"name,omitempty"`
	CECount       uint64           `json:"ceCount"`
	CENoInfoCount uint64           `json:"ceNoInfoCount"`
	UECount       uint64           `json:"ueCount"`
	UENoInfoCount uint64           `json:"ueNoInfoCount"`
	CSRows        []*ChipSelectRow `json:"csRows,omitempty"`
	MemorySizeSpec
}

// HasError returns whether controller has error
func (m MemoryController) HasError() bool {
	return (m.CECount > CECountThreshold || m.CENoInfoCount > 0 || m.UECount > 0 || m.UENoInfoCount > 0)
}

// Summary returns summarized strings
func (m MemoryController) Summary() string {
	fmtr := "mc: %s, %s, ce=%d(noinfo=%d), ue=%d(noinfo=%d)"
	return fmt.Sprintf(fmtr, m.Name, m.SizeString(), m.CECount, m.CENoInfoCount, m.UECount, m.UENoInfoCount)
}

// Summaries returns summarized strings
func (m MemoryController) Summaries() []string {
	var res []string

	res = append(res, m.Summary())
	for _, row := range m.CSRows {
		if row.HasError() {
			res = append(res, row.Summaries()...)
		}
	}

	return res
}

// ChipSelectRow represents a Chip-Select Row (csrowX)
type ChipSelectRow struct {
	Name     string           `json:"name,omitempty"`
	CECount  uint64           `json:"ceCount"`
	UECount  uint64           `json:"ueCount"`
	Channels []*MemoryChannel `json:"channels,omitempty"`
	MemorySizeSpec
}

// HasError returns whether csrow has error
func (c ChipSelectRow) HasError() bool {
	return (c.CECount > CECountThreshold || c.UECount > 0)
}

// Summary returns summarized strings
func (c ChipSelectRow) Summary() string {
	return fmt.Sprintf("cs: %s, %s, ce=%d, ue=%d", c.Name, c.SizeString(), c.CECount, c.UECount)
}

// Summaries returns summarized strings
func (c ChipSelectRow) Summaries() []string {
	var res []string
	res = append(res, c.Summary())
	for _, ch := range c.Channels {
		if ch.HasError() {
			res = append(res, ch.Summary())
		}
	}

	return res
}

// MemoryChannel represents a Channel table (chX)
type MemoryChannel struct {
	Name    string `json:"name,omitempty"`
	Label   string `json:"label,omitempty"`
	CECount uint64 `json:"ceCount"`
}

// HasError returns whether ch has error
func (m MemoryChannel) HasError() bool {
	return (m.CECount > 0)
}

// Summary returns summarized string
func (m MemoryChannel) Summary() string {
	return fmt.Sprintf("ch: %s, %s, ce=%d", m.Name, m.Label, m.CECount)
}
