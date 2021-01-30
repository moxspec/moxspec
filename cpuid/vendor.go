package cpuid

import (
	"fmt"
)

// Internal bit assign to handle vendors
const (
	INTEL = 1 << iota
	AMD
)

var vsigToVID = map[string]uint32{
	"GenuineIntel": INTEL,
	"AMDisbetter!": AMD,
	"AuthenticAMD": AMD,
}

func parseVendorID(cpu *Processor) error {
	log.Debug("parsing vendor info and maximum supported std level")

	eax, ebx, ecx, edx := cpuid(cpu, 0x00, 0)
	log.Debugf("maximum supported standard leaf is %04xh", eax)
	cpu.MaxStdLevel = eax

	vsig := string(readUint32sBackward(ebx, edx, ecx))
	log.Debugf("vendor signature: %s", vsig)
	cpu.VendorSignature = vsig
	if v, ok := vsigToVID[vsig]; ok {
		cpu.VendorID = v
		return nil
	}

	return fmt.Errorf("unsupported vendor: %s", vsig)
}
