package model

import (
	"fmt"
	"strings"
)

// ProcessorReport represents a processor report
type ProcessorReport struct {
	SocketCount    uint32     `json:"socketCount,omitempty"`
	PopulatedCount uint32     `json:"populatedCount,omitempty"`
	CoreCount      uint32     `json:"coreCount,omitempty"`
	ThreadCount    uint32     `json:"threadCount,omitempty"`
	Packages       []*Package `json:"packages,omitempty"`
}

// Summary returns summarized string
func (p ProcessorReport) Summary() string {
	// TODO: improve
	proc := p.Packages[0]

	var sockets = len(p.Packages)
	var brand = proc.ProductName

	fmtr := "%d x %s (%d cores, %d threads) (%d/%d sockets)"
	return fmt.Sprintf(fmtr, sockets, brand, p.CoreCount, p.ThreadCount, p.PopulatedCount, p.SocketCount)
}

// Package represents a processor package
type Package struct {
	ID            uint16   `json:"id,omitempty"`
	Socket        string   `json:"socket,omitempty"`
	Manufacturer  string   `json:"manufacturer,omitempty"`
	ProductName   string   `json:"productName,omitempty"`
	SerialNumber  string   `json:"serialNumber,omitempty"`
	CoreCount     uint32   `json:"coreCount,omitempty"`
	ThreadCount   uint32   `json:"threadCount,omitempty"`
	ThrottleCount uint16   `json:"throttleCount,omitempty"`
	Caches        []*Cache `json:"caches,omitempty"`
	TLBs          []*TLB   `json:"tlbs,omitempty"`
	Nodes         []*Node  `json:"nodes,omitempty"`
}

// Cache represents a data or instruction cache
type Cache struct {
	Level uint16 `json:"level,omitempty"`
	Type  string `json:"type,omitempty"`
	Size  uint32 `json:"size,omitempty"`
	Ways  uint32 `json:"ways,omitempty"`
}

// Summary returns summarized string
func (c Cache) Summary() string {
	return fmt.Sprintf("L%d %s %dKiB/core %d-way", c.Level, c.Type, c.Size, c.Ways)
}

// TLB represents a transaction lookaside buffer
type TLB struct {
	Type     string `json:"type,omitempty"`
	PageSize string `json:"pageSize,omitempty"`
	Entries  uint32 `json:"entries,omitempty"`
	Ways     uint32 `json:"ways,omitempty"`
	Fully    bool   `json:"fully,omitempty"` // indicates fully associative
}

// Summary returns summarized string
func (t TLB) Summary() string {
	var w string
	if t.Fully {
		w = "fully associative"
	} else {
		w = fmt.Sprintf("%d-ways", t.Ways)
	}

	return fmt.Sprintf("%s %d-entries %s", t.Type+" ("+t.PageSize+")", t.Entries, w)
}

// Node represents a numa node
type Node struct {
	ID        uint16  `json:"id,omitempty"`
	AvgTemp   float64 `json:"avgTemp,omitempty"`
	CoreCount uint32  `json:"coreCount,omitempty"`
	Cores     []*Core `json:"cores,omitempty"`
}

// ThrottledCores returns summary of throttled cores
func (n Node) ThrottledCores() string {
	var list []string
	for _, c := range n.Cores {
		if c.ThrottleCount > 0 {
			list = append(list, fmt.Sprintf("%d", c.ID))
		}
	}
	return strings.Join(list, ",")
}

// Summary returns summarized string
func (n Node) Summary() string {
	var sum string

	if n.AvgTemp > 0.0 {
		sum = fmt.Sprintf("node%d (%d cores) (avg %.1f°C)", n.ID, n.CoreCount, n.AvgTemp)
	} else {
		sum = fmt.Sprintf("node%d (%d cores)", n.ID, n.CoreCount)
	}

	throttled := n.ThrottledCores()
	if throttled != "" {
		sum = fmt.Sprintf("%s (thermal: throttled in %s)", sum, throttled)
	} else {
		sum = fmt.Sprintf("%s (thermal: safe)", sum)
	}

	return sum
}

// Core represents a physical core
type Core struct {
	ID            uint16 `json:"id,omitempty"`
	ThrottleCount uint16 `json:"throttleCount"`
	Temp          int16  `json:"temp,omitempty"`
}
