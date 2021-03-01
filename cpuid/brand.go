package cpuid

import (
	"github.com/moxspec/moxspec/util"
)

func parseBrandString(cpu *Processor) error {
	log.Debug("parsing brand string")
	cpu.BrandString += string(readUint32sBackward(cpuid(cpu, 0x80000002, 0)))
	cpu.BrandString += string(readUint32sBackward(cpuid(cpu, 0x80000003, 0)))
	cpu.BrandString += string(readUint32sBackward(cpuid(cpu, 0x80000004, 0)))
	log.Debugf("Brand String: %s", cpu.BrandString)
	cpu.BrandString = util.ShortenProcName(cpu.BrandString)
	return nil
}
