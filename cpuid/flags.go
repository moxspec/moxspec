package cpuid

func parseFlags(cpu *Processor) error {
	log.Debug("parsing processor type/family/model/stepping")
	parseStd0001h(cpu)
	parseExt0000h(cpu)
	return nil
}

func parseStd0001h(cpu *Processor) error {
	eax, _, _, _ := cpuid(cpu, 0x01, 0)

	cpu.ExtFamily = (eax >> 20) & 0xFF
	cpu.ExtModel = (eax >> 16) & 0xF
	cpu.Type = (eax >> 12) & 0x3
	cpu.Family = (eax >> 8) & 0xF
	cpu.Model = (eax >> 4) & 0xF
	cpu.Stepping = eax & 0xF

	log.Debugf("family: %x, ext-family: %x, model: %x, ext-model:%x, stepping: %d, type: %d",
		cpu.Family, cpu.ExtFamily,
		cpu.Model, cpu.ExtModel,
		cpu.Stepping, cpu.Type,
	)

	return nil
}

func parseExt0000h(cpu *Processor) error {
	cpu.MaxExtLevel, _, _, _ = cpuid(cpu, 0x80000000, 0)
	return nil
}
