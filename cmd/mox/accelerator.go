package main

import (
	"github.com/moxspec/moxspec/gpu/nvidia"
	"github.com/moxspec/moxspec/model"
	"github.com/moxspec/moxspec/pci"
)

func shapeAccelerater(r *model.Report, pcidevs *pci.Devices) {
	ar := new(model.AcceleratorReport)

	// GPU
	var nvGPUs []*model.GPU
	for _, dctl := range pcidevs.FilterByClass(pci.DisplayController) {
		if dctl.SubClassID != 0x02 { // 3D controller
			continue
		}

		g := new(model.GPU)
		g.PCIBaseSpec = *shapePCIDevice(dctl)

		switch dctl.Driver {
		case "nvidia":
			nvGPUs = append(nvGPUs, g)
		}

		ar.GPUs = append(ar.GPUs, g)
	}

	if len(nvGPUs) > 0 {
		nvd := nvidia.NewDecoder()
		err := nvd.Decode()
		if err == nil {
			for _, g := range nvGPUs {
				ngpu := nvd.GetGPU(g.PCIID())
				if ngpu == nil {
					continue
				}
				shapeNvidiaGPU(g, ngpu)
			}
		}
	}

	// FPGA
	for _, fctl := range pcidevs.FilterByClass(pci.CommunicationController, pci.ProcessingAccelerator) {
		if fctl.ClassID == pci.CommunicationController && fctl.SubClassID != 0x00 { // Serial controller
			continue
		}

		f := new(model.FPGA)
		f.PCIBaseSpec = *shapePCIDevice(fctl)
		ar.FPGAs = append(ar.FPGAs, f)
	}

	r.Accelerator = ar
}

func shapeNvidiaGPU(g *model.GPU, ngpu *nvidia.GPU) {
	g.ProductName = ngpu.ProductName
	g.SerialNumber = ngpu.Serial
	g.BIOS = ngpu.VBIOSVersion

	g.Power.Current = ngpu.Power.Draw
	g.Power.Limit = ngpu.Power.Limit

	g.Util.GPU = ngpu.Util.GPU
	g.Util.Memory = ngpu.Util.Memory

	g.Temp.GPU = ngpu.Temp.GPU
	g.Temp.Memory = ngpu.Temp.Memory

	g.CECount = shapeNvidiaECCCounter(ngpu.ECCErrors.Aggregate.SingleBit)
	g.UECount = shapeNvidiaECCCounter(ngpu.ECCErrors.Aggregate.DoubleBit)
}

func shapeNvidiaECCCounter(ne nvidia.ECCCounter) model.GPUECCCounter {
	var e model.GPUECCCounter
	e.DeviceMemory = ne.DeviceMemory
	e.RegisterFile = ne.RegisterFile
	e.L1Cache = ne.L1Cache
	e.L2Cache = ne.L2Cache
	e.Total = ne.Total
	return e
}
