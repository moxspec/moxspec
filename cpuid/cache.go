package cpuid

import (
	"strings"
)

const (
	dataCache    = "Data"
	codeCache    = "Code"
	unifiedCache = "Unified"
)

var cacheTypeOrder = map[string]uint16{
	dataCache:    0x00,
	codeCache:    0x01,
	unifiedCache: 0x02,
}

func parseCache(cpu *Processor) error {
	log.Debug("parsing cache descriptors")
	if cpu.VendorID == INTEL {
		return parseCacheBy(cpu, 0x00000004)
	}
	return parseCacheBy(cpu, 0x8000001D)
}

func getCacheType(typeCode uint32) string {
	switch typeCode {
	case 0x01:
		return dataCache
	case 0x02:
		return codeCache
	case 0x03:
		return unifiedCache
	}
	return ""
}

func parseCacheBy(cpu *Processor, eaxIn uint32) error {
	var i uint32
	for i = 0; ; i++ {
		eax, ebx, ecx, edx := cpuid(cpu, eaxIn, i)

		typeCode := eax & 0x1F
		if typeCode == 0 {
			break
		}

		c := new(Cache)
		c.Type = getCacheType(typeCode)
		if c.Type == "" {
			log.Warnf("%x is unknown cache type", typeCode)
			continue
		}

		c.Level = uint16((eax >> 5) & 0x7)

		if eax&(1<<8) != 0 {
			c.Flags = append(c.Flags, "self-initializing")
		}
		if eax&(1<<9) != 0 {
			c.Flags = append(c.Flags, "fully associative")
		}

		c.ThreadsPerCache = ((eax >> 14) & 0xFFF) + 1

		c.LineSize = (ebx & 0xFFF) + 1
		c.Partitions = ((ebx >> 12) & 0x3FF) + 1
		c.Ways = ((ebx >> 22) & 0x3FF) + 1

		c.Sets = ecx + 1

		c.Size = (c.Ways * c.Partitions * c.LineSize * c.Sets) >> 10

		if edx&1 == 0 {
			c.Flags = append(c.Flags, "write-back")
		}

		if edx&(1<<1) != 0 {
			c.Flags = append(c.Flags, "inclusive")
		}

		if eaxIn == 0x00000004 && edx&(1<<2) != 0 {
			c.Flags = append(c.Flags, "complex indexing")
		}

		cpu.Caches = append(cpu.Caches, c)

		log.Debugf("L%d %s %dKB (per %d threads): %dB line-size, %d sets, %d-ways (%s)", c.Level, c.Type, c.Size, c.ThreadsPerCache, c.LineSize, c.Sets, c.Ways, strings.Join(c.Flags, ", "))
	}
	return nil
}
