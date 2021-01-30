package main

import (
	"strings"

	"github.com/actapio/moxspec/cpu"
	"github.com/actapio/moxspec/cpuid"
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/msr"
	"github.com/actapio/moxspec/smbios"
	"github.com/actapio/moxspec/util"
)

func shapeProcessor(r *model.Report, sm []*smbios.Processor) {
	if sm == nil {
		return
	}

	var err error
	r.Processor = new(model.ProcessorReport)

	var sockets, populated uint32
	var pkgs []*model.Package
	for _, smproc := range sm {
		sockets++
		// TODO: to be improved
		if strings.Contains(strings.ToLower(strings.Join(smproc.Status, ",")), "unpopulated") {
			continue
		}
		populated++

		p := new(model.Package)
		p.Manufacturer = smproc.Manufacturer
		p.SerialNumber = smproc.SerialNumber
		p.Socket = smproc.SocketDesignation
		p.ProductName = smproc.Version
		p.CoreCount = uint32(smproc.CoreCount)
		p.ThreadCount = uint32(smproc.ThreadCount)

		pkgs = append(pkgs, p)
	}
	r.Processor.SocketCount = sockets
	r.Processor.PopulatedCount = populated

	cpuidd := cpuid.NewDecoder()
	err = cpuidd.Decode()
	if err != nil {
		log.Debug(err)
		log.Info("using smbibos as primary data source")
	} else {
		for i, p := range pkgs {
			if cpuidd.BrandString != "" {
				p.ProductName = cpuidd.BrandString
				log.Debugf("package %d: productname is set from cpuid: %s", i, p.ProductName)
			}

			if cpuidd.PhysicalCoreCount > 0 {
				p.CoreCount = cpuidd.PhysicalCoreCount
				log.Debugf("package %d: core count is set from cpuid: %d", i, p.CoreCount)
			}

			if cpuidd.LogicalCoreCount > 0 {
				p.ThreadCount = cpuidd.LogicalCoreCount
				log.Debugf("package %d: thread count is set from cpuid: %d", i, p.ThreadCount)
			}

			// an old virtualized platform does not support cpuid(0x0Bh)
			if p.CoreCount > 0 && p.ThreadCount > 0 {
				shapeProcessorCache(p, cpuidd)
				shapeProcessorTLB(p, cpuidd)
			}
		}
	}

	cput := cpu.NewDecoder()
	err = cput.Decode()
	if err != nil {
		log.Error(err)
		log.Error("could not parse processor information from the kernel")
		return
	}

	shapeProcessorNode(pkgs, cput)

	var cores, threads uint32
	for _, p := range pkgs {
		cores = cores + p.CoreCount
		threads = threads + p.ThreadCount
	}
	r.Processor.CoreCount = cores
	r.Processor.ThreadCount = threads

	r.Processor.Packages = pkgs
}

func shapeProcessorCache(p *model.Package, cpuidd *cpuid.Processor) {
	for _, c := range cpuidd.GetCaches() {
		size := c.Size
		tpc := c.ThreadsPerCache
		if tpc > cpuidd.LogicalCoreCount {
			tpc = cpuidd.LogicalCoreCount
		}

		if cpuidd.SMT == 0 {
			log.Warn("cpuid SMT is 0")
			return
		}

		perCore := uint32(tpc / cpuidd.SMT)
		if perCore > 1 {
			size = uint32(size / perCore)
		}

		cache := new(model.Cache)
		cache.Level = c.Level
		cache.Type = c.Type
		cache.Size = size
		cache.Ways = c.Ways

		p.Caches = append(p.Caches, cache)
	}
}

func shapeProcessorTLB(p *model.Package, cpuidd *cpuid.Processor) {
	for _, t := range cpuidd.GetTLBs() {
		tlb := new(model.TLB)
		tlb.Type = t.Type
		tlb.PageSize = t.PageSize
		tlb.Entries = t.Entries
		tlb.Ways = t.Ways
		if t.Ways == 0xFF {
			tlb.Fully = true
		}

		p.TLBs = append(p.TLBs, tlb)
	}
}

func shapeProcessorNode(pkgs []*model.Package, cput *cpu.Topology) {
	cputPkgs := cput.Packages()
	if len(cputPkgs) != len(pkgs) {
		log.Warn("package count reported by the kernel is not match with smbios report")
		log.Warn("please check your kernel implementation")
		return
	}

	for i, p := range cputPkgs {
		if i >= len(pkgs) {
			break
		}

		pkgs[i].ThrottleCount = p.ThrottleCount
		if pkgs[i].CoreCount == 0 {
			pkgs[i].CoreCount = uint32(p.CoreCount())
			log.Debugf("package %d: core count is set from the kernel: %d", i, pkgs[i].CoreCount)
		}
		if pkgs[i].ThreadCount == 0 {
			pkgs[i].ThreadCount = uint32(p.ThreadCount())
			log.Debugf("package %d: thread count is set from the kernel: %d", i, pkgs[i].ThreadCount)
		}

		for _, nd := range p.Nodes() {
			n := new(model.Node)
			n.ID = nd.ID

			for _, cr := range nd.Cores() {
				c := new(model.Core)
				c.ID = cr.ID
				c.ThrottleCount = cr.ThrottleCount

				r := msr.NewDecoder(cr.Threads[0], msr.INTEL)
				err := r.Decode()
				if err == nil {
					c.Temp = r.Temp
				}

				n.Cores = append(n.Cores, c)
			}

			var cnt uint32
			var tsum float64
			for _, c := range n.Cores {
				cnt++
				if c.Temp > 0 {
					tsum = tsum + float64(c.Temp)
				}
			}
			if cnt > 0 {
				n.AvgTemp = util.Round(tsum/float64(cnt), 1)
			}
			n.CoreCount = cnt

			pkgs[i].Nodes = append(pkgs[i].Nodes, n)
		}

		pkgs[i].ID = p.ID
	}
}
