package cpuid

import (
	"fmt"
	"os"

	"github.com/moxspec/moxspec/util"
)

// > INPUT
//   (uint32) EBX = 'u' 'n' 'e' 'G'
//   (uint32) EDX = 'I' 'e' 'n' 'i'
//   (uint32) ECX = 'l' 'e' 't' 'n'
//
// > OUTPUT
//   ([]byte) GenuineIntel
func readUint32sBackward(us ...uint32) []byte {
	var ret []byte

	for _, u := range us {
		ret = append(ret, uint32toBytes(u)...)
	}

	return ret
}

func uint32toBytes(u uint32) []byte {
	var ret []byte

	ret = append(ret,
		byte(u&0xFF),       // 1st byte
		byte((u>>8)&0xFF),  // 2nd byte
		byte((u>>16)&0xFF), // 3rd byte
		byte((u>>24)&0xFF), // 4th byte
	)

	return ret
}

func cpuid(p *Processor, eaxIn, ecxIn uint32) (eax, ebx, ecx, edx uint32) {
	ok, err := isValidEaxIn(eaxIn, p.MaxStdLevel, p.MaxExtLevel)
	if !ok {
		log.Debug(err.Error())
		return
	}

	// TODO: to be configurable
	fd, err := os.OpenFile("/dev/cpu/0/cpuid", os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return
	}
	defer fd.Close()

	d := make([]byte, 16, 16)
	var ptr int64
	ptr = (int64(ecxIn) << 32) | int64(eaxIn)
	fd.ReadAt(d, ptr)

	eax = util.BytesToUint32(d[0:4])
	ebx = util.BytesToUint32(d[4:8])
	ecx = util.BytesToUint32(d[8:12])
	edx = util.BytesToUint32(d[12:])

	log.Debugf("CPUID(eaxIn:%08x ecxIn:%08x) => eax:%08x ebx:%08x ecx:%08x edx:%08x", eaxIn, ecxIn, eax, ebx, ecx, edx)
	return
}

func isValidEaxIn(eaxIn, maxStdLv, maxExtLv uint32) (bool, error) {
	var cat = uint16(eaxIn >> 16)
	if cat == 0x0000 && maxStdLv != 0 && eaxIn > maxStdLv {
		return false, fmt.Errorf("CPUID(%08x) is NOT supported (max: %08x)", eaxIn, maxStdLv)
	}
	if cat == 0x8000 && maxExtLv != 0 && eaxIn > maxExtLv {
		return false, fmt.Errorf("CPUID(%08x) is NOT supported (max: %08x)", eaxIn, maxExtLv)
	}
	return true, nil
}

// Skylake-SP has CPUID errata
// https://www.intel.com/content/dam/www/public/us/en/documents/specification-updates/xeon-scalable-spec-update.pdf
//
// 59. CPUID TLB Associativity Information is Inaccurate
// Problem: CPUID leaf 2 (EAX=02H) TLB information inaccurately reports
//          that the shared second- Level TLB is 6-way set associative (value C3H),
//          although it is 12-way set associative.
//          Other information reported by CPUID leaf 2 is accurate.
// Status:  No Fix
func isSkylakeSP(cpu *Processor) bool {
	return (cpu.Model == 5 && cpu.ExtModel == 5 && cpu.Family == 6 && cpu.ExtFamily == 0 && cpu.Type == 0)
}
