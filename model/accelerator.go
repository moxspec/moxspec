package model

import "fmt"

// AcceleratorReport represents the accelerator report
type AcceleratorReport struct {
	GPUs  []*GPU  `json:"gpus,omitempty"`
	FPGAs []*FPGA `json:"fpgas,omitempty"`
}

// GPU represents the GPU device
type GPU struct {
	PCIBaseSpec
	ProductName string        `json:"productName,omitempty"`
	BIOS        string        `json:"bios,omitempty"`
	Power       GPUPower      `json:"power,omitempty"`
	Util        GPUUtil       `json:"util,omitempty"`
	Temp        GPUTemp       `json:"temp,omitempty"`
	CECount     GPUECCCounter `json:"ceCount,omitempty"`
	UECount     GPUECCCounter `json:"ueCount,omitempty"`
}

// IsHealthy returns whether the GPU  is healthy
func (g GPU) IsHealthy() bool {
	if !g.CECount.IsHealthy() {
		return false
	}

	if !g.UECount.IsHealthy() {
		return false
	}

	return g.PCIBaseSpec.IsHealthy()
}

// DiagSummaries returns diag status
func (g GPU) DiagSummaries() []string {
	sum := g.PCIBaseSpec.DiagSummaries()

	if !g.CECount.IsHealthy() {
		sum = append(sum, g.CECount.diagSummariesWithPrefix("ce")...)
	}

	if !g.UECount.IsHealthy() {
		sum = append(sum, g.UECount.diagSummariesWithPrefix("ue")...)
	}

	return sum
}

// LongName returns pretty name
func (g GPU) LongName() string {
	if g.VendorName != "" && g.ProductName != "" {
		return fmt.Sprintf("%s %s", g.VendorName, g.ProductName)
	}

	return g.PCIBaseSpec.LongName()
}

// PowerSummary returns summarized power status
func (g GPU) PowerSummary() string {
	return g.Power.Summary()
}

// TempSummary returns summarized temperature status
func (g GPU) TempSummary() string {
	return g.Temp.Summary()
}

// Summary returns summarized string
func (g GPU) Summary() string {
	var sum string
	if g.Driver == "" {
		sum = fmt.Sprintf("%s (node%d)", g.LongName(), g.Numa)
	} else {
		sum = fmt.Sprintf("%s (%s) (node%d)", g.LongName(), g.Driver, g.Numa)
	}

	if g.SerialNumber != "" {
		sum = fmt.Sprintf("%s (SN:%s)", sum, g.SerialNumber)
	}
	return sum
}

// GPUPower represents GPU power status
type GPUPower struct {
	Current float32 `json:"current,omitempty"`
	Limit   float32 `json:"limit,omitempty"`
}

// Summary returns summarized power status
func (g GPUPower) Summary() string {
	return fmt.Sprintf("cur %.1fW, limit %.1fW", g.Current, g.Limit)
}

// GPUTemp represents GPU temperature status
type GPUTemp struct {
	GPU    float32 `json:"gpu,omitempty"`
	Memory float32 `json:"memory,omitempty"`
}

// Summary returns summarized temperature status
func (g GPUTemp) Summary() string {
	return fmt.Sprintf("GPU %.1f°C, Memory %.1f°C", g.GPU, g.Memory)
}

// GPUUtil represents GPU utility status
type GPUUtil struct {
	GPU    float32 `json:"gpu,omitempty"`
	Memory float32 `json:"memory,omitempty"`
}

// GPUECCCounter represents GPU ECC counters
type GPUECCCounter struct {
	DeviceMemory int `json:"deviceMemory,omitempty"`
	RegisterFile int `json:"registerFile,omitempty"`
	L1Cache      int `json:"l1cache,omitempty"`
	L2Cache      int `json:"l2cache,omitempty"`
	Total        int `json:"total,omitempty"`
}

// IsHealthy returns whether the device is healthy
func (g GPUECCCounter) IsHealthy() bool {
	return (g.Total == 0)
}

func (g GPUECCCounter) diagSummariesWithPrefix(prefix string) []string {
	var sum []string
	if g.DeviceMemory > 0 {
		sum = append(sum, fmt.Sprintf("[%s] DeviceMemory: %d", prefix, g.DeviceMemory))
	}
	if g.RegisterFile > 0 {
		sum = append(sum, fmt.Sprintf("[%s] RegisterFile: %d", prefix, g.RegisterFile))
	}
	if g.L1Cache > 0 {
		sum = append(sum, fmt.Sprintf("[%s] L1Cache: %d", prefix, g.L1Cache))
	}
	if g.L2Cache > 0 {
		sum = append(sum, fmt.Sprintf("[%s] L2Cache: %d", prefix, g.L2Cache))
	}
	sum = append(sum, fmt.Sprintf("[%s] Total: %d", prefix, g.Total))

	return sum
}

// FPGA represents the FPGA device
type FPGA struct {
	PCIBaseSpec
}

// Summary returns summarized string
func (f FPGA) Summary() string {
	var sum string
	if f.Driver == "" {
		sum = fmt.Sprintf("%s (node%d)", f.LongName(), f.Numa)
	} else {
		sum = fmt.Sprintf("%s (%s) (node%d)", f.LongName(), f.Driver, f.Numa)
	}

	if f.SerialNumber != "" {
		sum = fmt.Sprintf("%s (SN:%s)", sum, f.SerialNumber)
	}
	return sum
}
