package cpuid

// NOTE: CPUID(0x16)  is supported Skylake or later
func parseFrequency(cpu *Processor) error {
	log.Debug("parsing frequency")

	if cpu.VendorID != INTEL {
		return nil
	}

	eax, ebx, _, _ := cpuid(cpu, 0x16, 0)
	if eax == 0 && ebx == 0 {
		return nil
	}

	cpu.BaseFrequency = eax & 0xFFFF
	cpu.MaximumFrequency = ebx & 0xFFFF

	log.Debugf("Base Frequency: %d MHz", cpu.BaseFrequency)
	log.Debugf("Maximum Frequency: %d MHz", cpu.MaximumFrequency)

	return nil
}
