package cpuid

const (
	l0Code    = "L0 Code"
	l0Data    = "L0 Data"
	l1Code    = "L1 Code"
	l1Data    = "L1 Data"
	l2Code    = "L2 Code"
	l2Data    = "L2 Data"
	l2Unified = "L2 Unified"
)

var tlbTypeOrder = map[string]uint16{
	l0Code:    0x00,
	l0Data:    0x01,
	l1Code:    0x02,
	l1Data:    0x03,
	l2Code:    0x04,
	l2Data:    0x05,
	l2Unified: 0x06,
}

const (
	page4k     = "4K"
	page2m     = "2M"
	page4m     = "4M"
	page4k2m   = "4K/2M"
	page4k4m   = "4K/4M"
	page2m4m   = "2M/4M"
	page4k2m4m = "4K/2M/4M"
	page1g     = "1G"
)

var tlbPageOrder = map[string]uint16{
	page4k:     0x00,
	page2m:     0x01,
	page4m:     0x02,
	page4k2m:   0x03,
	page4k4m:   0x04,
	page2m4m:   0x05,
	page4k2m4m: 0x06,
	page1g:     0x07,
}

func parseTLB(cpu *Processor) error {
	log.Debug("parsing TLB descriptors")
	if cpu.VendorID == INTEL {
		return parseStd0002h(cpu)
	}

	parseExt0005h(cpu)
	parseExt0006h(cpu)
	parseExt0019h(cpu)

	return nil
}

func appendTLB(cpu *Processor, t *TLB) {
	if t != nil && t.Ways != 0 && t.Entries != 0 {
		cpu.TLBs = append(cpu.TLBs, t)
	}
}

func parseStd0002h(cpu *Processor) error {
	eax, ebx, ecx, edx := cpuid(cpu, 0x2, 0)

	// AX is number of times this level must be queried to obtain all configuration descriptors
	parseStd0002hDescriptors(cpu, eax&0xFFFFFF00)

	parseStd0002hDescriptors(cpu, ebx)
	parseStd0002hDescriptors(cpu, ecx)
	parseStd0002hDescriptors(cpu, edx)
	return nil
}

func newTLB(t, ps string, w, e uint32) *TLB {
	tlb := new(TLB)
	tlb.Type = t
	tlb.PageSize = ps
	tlb.Ways = w
	tlb.Entries = e
	return tlb
}

func parseStd0002hDescriptors(cpu *Processor, reg uint32) {
	// this table has ONLY TLB descriptors
	// see also: http://sandpile.org/x86/cpuid.htm#level_0000_0002h
	descriptors := map[byte][]*TLB{
		//0x01: "code TLB, 4K pages, 4 ways, 32 entries",
		0x01: []*TLB{newTLB(l1Code, page4k, 4, 32)},

		//0x02: "code TLB, 4M pages, fully, 2 entries",
		0x02: []*TLB{newTLB(l1Code, page4m, 0xFF, 2)},

		//0x03: "data TLB, 4K pages, 4 ways, 64 entries",
		0x03: []*TLB{newTLB(l1Data, page4k, 4, 64)},

		//0x04: "data TLB, 4M pages, 4 ways, 8 entries",
		0x04: []*TLB{newTLB(l1Data, page4m, 4, 8)},

		//0x05: "data TLB, 4M pages, 4 ways, 32 entries",
		0x05: []*TLB{newTLB(l1Data, page4m, 4, 32)},

		//0x0B: "code TLB, 4M pages, 4 ways, 4 entries",
		0x0B: []*TLB{newTLB(l1Code, page4m, 4, 4)},

		//0x4F: "code TLB, 4K pages, ???, 32 entries",
		0x4F: []*TLB{newTLB(l1Code, page4k, 0x00, 32)},

		//0x50: "code TLB, 4K/4M/2M pages, fully, 64 entries",
		0x50: []*TLB{newTLB(l1Code, page4k2m4m, 0xFF, 64)},

		//0x51: "code TLB, 4K/4M/2M pages, fully, 128 entries",
		0x51: []*TLB{newTLB(l1Code, page4k2m4m, 0xFF, 128)},

		//0x52: "code TLB, 4K/4M/2M pages, fully, 256 entries",
		0x52: []*TLB{newTLB(l1Code, page4k2m4m, 0xFF, 256)},

		//0x55: "code TLB, 2M/4M, fully, 7 entries",
		0x55: []*TLB{newTLB(l1Code, page2m4m, 0xFF, 7)},

		//0x56: "L0 data TLB, 4M pages, 4 ways, 16 entries",
		0x56: []*TLB{newTLB(l0Data, page4m, 4, 16)},

		//0x57: "L0 data TLB, 4K pages, 4 ways, 16 entries",
		0x57: []*TLB{newTLB(l0Data, page4k, 4, 16)},

		//0x59: "L0 data TLB, 4K pages, fully, 16 entries",
		0x59: []*TLB{newTLB(l0Data, page4k, 0xFF, 16)},

		//0x5A: "L0 data TLB, 2M/4M, 4 ways, 32 entries",
		0x5A: []*TLB{newTLB(l0Data, page2m4m, 4, 32)},

		//0x5B: "data TLB, 4K/4M pages, fully, 64 entries",
		0x5B: []*TLB{newTLB(l1Data, page4k4m, 0xFF, 64)},

		//0x5C: "data TLB, 4K/4M pages, fully, 128 entries",
		0x5C: []*TLB{newTLB(l1Data, page4k4m, 0xFF, 128)},

		//0x5D: "data TLB, 4K/4M pages, fully, 256 entries",
		0x5D: []*TLB{newTLB(l1Data, page4k4m, 0xFF, 256)},

		//0x61: "code TLB, 4K pages, fully, 48 entries",
		0x61: []*TLB{newTLB(l1Code, page4k, 0xFF, 48)},

		//0x63: "data TLB, 2M/4M pages, 4-way, 32-entries, and data TLB, 1G pages, 4-way, 4 entries",
		0x63: []*TLB{
			newTLB(l1Data, page2m4m, 4, 32),
			newTLB(l1Data, page1g, 4, 4),
		},

		//0x64: "data TLB, 4K pages, 4-way, 512 entries",
		0x64: []*TLB{newTLB(l1Data, page4k, 4, 512)},

		//0x6A: "L0 data TLB, 4K pages, 8-way, 64 entries",
		0x6A: []*TLB{newTLB(l0Data, page4k, 8, 64)},

		//0x6B: "data TLB, 4K pages, 8-way, 256 entries",
		0x6B: []*TLB{newTLB(l1Data, page4k, 8, 256)},

		//0x6C: "data TLB, 2M/4M pages, 8-way, 126 entries",
		0x6C: []*TLB{newTLB(l1Data, page2m4m, 8, 126)},

		//0x6D: "data TLB, 1G pages, fully, 16 entries",
		0x6D: []*TLB{newTLB(l1Data, page1g, 0xFF, 16)},

		//0x76: "code TLB, 2M/4M pages, fully, 8 entries",
		0x76: []*TLB{newTLB(l1Code, page2m4m, 0xFF, 8)},

		//0xA0: "data TLB, 4K pages, fully, 32 entries",
		0xA0: []*TLB{newTLB(l1Data, page4k, 0xFF, 32)},

		//0xB0: "code TLB, 4K pages, 4 ways, 128 entries",
		0xB0: []*TLB{newTLB(l1Code, page4k, 4, 128)},

		//0xB1: "code TLB, 4M pages, 4 ways, 4 entries and code TLB, 2M pages, 4 ways, 8 entries",
		0xB1: []*TLB{
			newTLB(l1Code, page4m, 4, 4),
			newTLB(l1Code, page2m, 4, 8),
		},

		//0xB2: "code TLB, 4K pages, 4 ways, 64 entries",
		0xB2: []*TLB{newTLB(l1Code, page4k, 4, 64)},

		//0xB3: "data TLB, 4K pages, 4 ways, 128 entries",
		0xB3: []*TLB{newTLB(l1Data, page4k, 4, 128)},

		//0xB4: "data TLB, 4K pages, 4 ways, 256 entries",
		0xB4: []*TLB{newTLB(l1Data, page4k, 4, 256)},

		//0xB5: "code TLB, 4K pages, 8 ways, 64 entries",
		0xB5: []*TLB{newTLB(l1Code, page4k, 8, 64)},

		//0xB6: "code TLB, 4K pages, 8 ways, 128 entries",
		0xB6: []*TLB{newTLB(l1Code, page4k, 8, 128)},

		//0xBA: "data TLB, 4K pages, 4 ways, 64 entries",
		0xBA: []*TLB{newTLB(l1Data, page4k, 4, 64)},

		//0xC0: "data TLB, 4K/4M pages, 4 ways, 8 entries",
		0xC0: []*TLB{newTLB(l1Data, page4k4m, 4, 8)},

		//0xC1: "L2 code and data TLB, 4K/2M pages, 8 ways, 1024 entries",
		0xC1: []*TLB{newTLB(l2Unified, page4k2m, 8, 1024)},

		//0xC2: "data TLB, 2M/4M pages, 4 ways, 16 entries",
		0xC2: []*TLB{newTLB(l1Data, page2m4m, 4, 16)},

		//0xC3: "L2 code and data TLB, 4K/2M pages, 6 ways, 1536 entries and L2 code and data TLB, 1G pages, 4 ways, 16 entries",
		0xC3: []*TLB{
			newTLB(l2Unified, page4k2m, 6, 1536),
			newTLB(l2Unified, page1g, 4, 16),
		},

		//0xC4: "data TLB, 2M/4M pages, 4-way, 32 entries",
		0xC4: []*TLB{newTLB(l1Data, page2m4m, 4, 32)},

		//0xCA: "L2 code and data TLB, 4K pages, 4 ways, 512 entries",
		0xCA: []*TLB{newTLB(l2Unified, page4k, 4, 512)},
	}

	prefetching := map[byte]uint32{
		0xF0: 64,  // 64 byte prefetching,
		0xF1: 128, // 128 byte prefetching,
	}

	for _, b := range uint32toBytes(reg) {
		if b == 0x00 {
			log.Debug("00 = null descriptor (=unused descriptor)")
			continue
		}

		if b == 0xFE { // query standard level 0000_0018h instead
			log.Debugf("%02x = query standard level 0000_0018h instead", b)
			log.Debug("CPUID(0000_0018h) decoder is not implemented.")
			continue
		}
		if b == 0xFF { // query standard level 0000_0004h instead
			log.Debugf("%02x = query standard level 0000_0004h instead", b)
			continue
		}

		if descs, ok := descriptors[b]; ok {
			log.Debugf("%02x = %+v", b, descs)
			if b == 0xC3 && isSkylakeSP(cpu) {
				log.Debug("Skylake-SP reports TLB info inaccurately. L2 TLB shoud be 12-way set associative.")
				l2u4k2m := newTLB(l2Unified, page4k2m, 12, 1536)
				log.Debugf("WRONG: %+v", descs[0])
				log.Debugf("RIGHT: %+v", l2u4k2m)
				appendTLB(cpu, l2u4k2m)

				l2u1g := newTLB(l2Unified, page1g, 0x04, 16)
				appendTLB(cpu, l2u1g)
				continue
			}

			for _, desc := range descs {
				appendTLB(cpu, desc)
			}
			continue
		}

		if pref, ok := prefetching[b]; ok {
			cpu.Prefetch = pref
			log.Debugf("%02x = %d byte prefetching", b, cpu.Prefetch)
			continue
		}

		log.Debugf("CPUID(0000_0002h) Code %02x is undefined.", b)
	}
}

func parseExt0005h(cpu *Processor) error {
	// eax = 4/2 MB L1 TLB configuration descriptor
	// ebx = 4 KB L1 TLB configuration descriptor
	eax, ebx, _, _ := cpuid(cpu, 0x80000005, 0)

	l1d2m4m := new(TLB)
	l1d2m4m.Type = l1Data
	l1d2m4m.PageSize = page2m4m
	l1d2m4m.Ways = (eax >> 24) & 0xFF
	l1d2m4m.Entries = (eax >> 16) & 0xFF
	log.Debugf("%+v", l1d2m4m)
	appendTLB(cpu, l1d2m4m)

	l1d4k := new(TLB)
	l1d4k.Type = l1Data
	l1d4k.PageSize = page4k
	l1d4k.Ways = (ebx >> 24) & 0xFF
	l1d4k.Entries = (ebx >> 16) & 0xFF
	log.Debugf("%+v", l1d4k)
	appendTLB(cpu, l1d4k)

	l1i2m4m := new(TLB)
	l1i2m4m.Type = l1Code
	l1i2m4m.PageSize = page2m4m
	l1i2m4m.Ways = (eax >> 8) & 0xFF
	l1i2m4m.Entries = eax & 0xFF
	log.Debugf("%+v", l1i2m4m)
	appendTLB(cpu, l1i2m4m)

	l1i4k := new(TLB)
	l1i4k.Type = l1Code
	l1i4k.PageSize = page4k
	l1i4k.Ways = (ebx >> 8) & 0xFF
	l1i4k.Entries = ebx & 0xFF
	log.Debugf("%+v", l1i4k)
	appendTLB(cpu, l1i4k)

	return nil
}

func parseExt0006h(cpu *Processor) error {
	// eax = 4/2 MB L2 TLB configuration descriptor
	// ebx = 4 KB L2 TLB configuration descriptor
	eax, ebx, _, _ := cpuid(cpu, 0x80000006, 0)

	l2d2m4m := new(TLB)
	l2d2m4m.Type = l2Data
	l2d2m4m.PageSize = page2m4m
	switch (eax >> 28) & 0xF {
	case 0x02:
		l2d2m4m.Ways = 2
	case 0x03:
		l2d2m4m.Ways = 3
	default:
		log.Debugf("CPUID(8000_0006) EAX[31..28] %x is undefined.", (eax>>28)&0xF)
	}
	l2d2m4m.Entries = (eax >> 16) & 0xFFF
	log.Debugf("%+v", l2d2m4m)
	appendTLB(cpu, l2d2m4m)

	l2d4k := new(TLB)
	l2d4k.Type = l2Data
	l2d4k.PageSize = page4k
	switch (ebx >> 28) & 0xF {
	case 0x05:
		l2d4k.Ways = 6
	case 0x06:
		l2d4k.Ways = 8
	default:
		log.Debugf("CPUID(8000_0006) EBX[31..28] %x is undefined.", (ebx>>28)&0xF)
	}
	l2d4k.Entries = (ebx >> 16) & 0xFFF
	log.Debugf("%+v", l2d4k)
	appendTLB(cpu, l2d4k)

	l2i2m4m := new(TLB)
	l2i2m4m.Type = l2Code
	l2i2m4m.PageSize = page2m4m
	switch (eax >> 12) & 0xF {
	case 0x06:
		l2i2m4m.Ways = 8
	default:
		log.Warnf("CPUID(8000_0006) EAX[15..12] %x is undefined.", (eax>>12)&0xF)
	}
	l2i2m4m.Entries = eax & 0xFFF
	log.Debugf("%+v", l2i2m4m)
	appendTLB(cpu, l2i2m4m)

	l2i4k := new(TLB)
	l2i4k.Type = l2Code
	l2i4k.PageSize = page4k
	switch (ebx >> 12) & 0xF {
	case 0x06:
		l2i4k.Ways = 8
	default:
		log.Warnf("CPUID(8000_0006) EBX[15..12] %x is undefined.", (ebx>>12)&0xF)
	}
	l2i4k.Entries = ebx & 0xFFF
	log.Debugf("%+v", l2i4k)
	appendTLB(cpu, l2i4k)

	return nil
}

func parseExt0019h(cpu *Processor) error {
	// eax = 1 GB L1 TLB configuration descriptor
	// ebx = 1 GB L2 TLB configuration descriptor
	eax, ebx, _, _ := cpuid(cpu, 0x80000019, 0)

	l1d1g := new(TLB)
	l1d1g.Type = l1Data
	l1d1g.PageSize = page1g
	l1d1g.Ways = (eax >> 28) & 0xF
	l1d1g.Entries = (eax >> 16) & 0xFFF
	log.Debugf("%+v", l1d1g)
	appendTLB(cpu, l1d1g)

	l1i1g := new(TLB)
	l1i1g.Type = l1Code
	l1i1g.PageSize = page1g
	l1i1g.Ways = (eax >> 12) & 0xF
	l1i1g.Entries = eax & 0xFFF
	log.Debugf("%+v", l1i1g)
	appendTLB(cpu, l1i1g)

	l2d1g := new(TLB)
	l2d1g.Type = l2Data
	l2d1g.PageSize = page1g
	l2d1g.Ways = (ebx >> 28) & 0xF
	l2d1g.Entries = (ebx >> 16) & 0xFFF
	log.Debugf("%+v", l2d1g)
	appendTLB(cpu, l2d1g)

	l2i1g := new(TLB)
	l2i1g.Type = l2Code
	l2i1g.PageSize = page1g
	l2d1g.Ways = (ebx >> 12) & 0xF
	l2d1g.Entries = ebx & 0xFFF
	log.Debugf("%+v", l2i1g)
	appendTLB(cpu, l2i1g)

	return nil
}
