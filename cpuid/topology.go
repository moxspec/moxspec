package cpuid

func parseTopology(cpu *Processor) error {
	log.Debug("parsing topology")

	var err error
	if cpu.VendorID == INTEL {
		err = parseStd0008Bh(cpu)
	} else {
		err = parseExt001Eh(cpu)
	}

	log.Debugf("Physical Cores: %d / Logical Cores: %d (SMT: %d)", cpu.PhysicalCoreCount, cpu.LogicalCoreCount, cpu.SMT)

	return err
}

func parseStd0008Bh(cpu *Processor) error {
	_, ebx, _, _ := cpuid(cpu, 0x0B, 0x00) // SMT
	cpu.SMT = ebx & 0xFF

	_, ebx, _, _ = cpuid(cpu, 0x0B, 0x01)
	cpu.LogicalCoreCount = ebx & 0xFF

	if cpu.SMT != 0 {
		cpu.PhysicalCoreCount = cpu.LogicalCoreCount / cpu.SMT
	}

	return nil
}

func parseExt001Eh(cpu *Processor) error {
	_, ebx, _, _ := cpuid(cpu, 0x8000001E, 0x00)
	cpu.SMT = ((ebx >> 8) & 0xFF) + 1

	_, _, ecx, _ := cpuid(cpu, 0x80000008, 0x00)
	cpu.LogicalCoreCount = (ecx & 0xFF) + 1

	if cpu.SMT != 0 {
		cpu.PhysicalCoreCount = cpu.LogicalCoreCount / cpu.SMT
	}

	return nil
}
