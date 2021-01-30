package cpuid

import (
	"sort"
)

// Processor represents a processor spec
type Processor struct {
	VendorID          uint32
	VendorSignature   string
	BrandString       string
	MaxStdLevel       uint32
	MaxExtLevel       uint32
	Family            uint32
	ExtFamily         uint32
	Model             uint32
	ExtModel          uint32
	Stepping          uint32
	Type              uint32
	Prefetch          uint32
	Caches            []*Cache
	TLBs              []*TLB
	BaseFrequency     uint32
	MaximumFrequency  uint32
	SMT               uint32
	LogicalCoreCount  uint32
	PhysicalCoreCount uint32
}

// GetCaches returns caches
func (p *Processor) GetCaches() []*Cache {
	sort.Slice(p.Caches, func(i, j int) bool {
		oI := (p.Caches[i].Level << 8) | cacheTypeOrder[p.Caches[i].Type]
		oJ := (p.Caches[j].Level << 8) | cacheTypeOrder[p.Caches[j].Type]
		return oI < oJ
	})
	return p.Caches
}

// GetTLBs returns tlbs
func (p *Processor) GetTLBs() []*TLB {
	sort.Slice(p.TLBs, func(i, j int) bool {
		oI := (tlbTypeOrder[p.TLBs[i].Type] << 8) | tlbPageOrder[p.TLBs[i].PageSize]
		oJ := (tlbTypeOrder[p.TLBs[j].Type] << 8) | tlbPageOrder[p.TLBs[j].PageSize]
		return oI < oJ
	})
	return p.TLBs
}

// Cache represents a cache spec
type Cache struct {
	Level           uint16
	Type            string
	ThreadsPerCache uint32
	LineSize        uint32
	Partitions      uint32
	Sets            uint32
	Size            uint32
	Ways            uint32
	Flags           []string
}

// TLB represents a tlb spec
type TLB struct {
	Type     string
	PageSize string
	Ways     uint32 // 0xFF = fully associative
	Entries  uint32
}
