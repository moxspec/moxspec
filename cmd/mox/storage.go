package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/actapio/moxspec/blk/ahci"
	"github.com/actapio/moxspec/blk/nvme"
	"github.com/actapio/moxspec/blk/raid"
	"github.com/actapio/moxspec/blk/virtio"
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/nvmeadm"
	"github.com/actapio/moxspec/pci"
	"github.com/actapio/moxspec/raidcli"
	"github.com/actapio/moxspec/raidcli/hpacucli"
	"github.com/actapio/moxspec/raidcli/megacli"
	"github.com/actapio/moxspec/raidcli/sas3ircu"
	"github.com/actapio/moxspec/spc"
	"github.com/actapio/moxspec/spc/acs"
	"github.com/actapio/moxspec/spc/megaraid"
	"github.com/actapio/moxspec/util"
)

// for switch-selecter
func matcher(f func(string) bool, t string, ts ...string) bool {
	list := append([]string{t}, ts...)
	for _, l := range list {
		if f(l) {
			return true
		}
	}
	return false
}

func shapeDisk(r *model.Report, pcidevs *pci.Devices, cli *app) {
	r.Storage = new(model.StorageReport)

	var nvmeCtls []*model.NVMeController
	var raidCtls []*model.RAIDController
	var ahciCtls []*model.AHCIController
	var virtCtls []*model.VirtController
	var nstdCtls []*model.NonStdController

	for _, ctl := range pcidevs.FilterByClass(pci.MassStorageController) {
		bspec := shapePCIDevice(ctl)

		equal := func(t string, ts ...string) bool {
			return matcher(func(l string) bool {
				if bspec.Driver == l {
					return true
				}
				return false
			}, t, ts...)
		}

		prefix := func(t string, ts ...string) bool {
			return matcher(func(l string) bool {
				if strings.HasPrefix(bspec.Driver, l) {
					return true
				}
				return false
			}, t, ts...)
		}

		var err error
		switch {
		case equal("nvme"):
			c, err := shapeNVMeController(bspec)
			if err == nil && c != nil {
				nvmeCtls = append(nvmeCtls, c)
			}
		case equal("ahci", "ata_piix", "isci"):
			c, err := shapeAHCIController(bspec)
			if err == nil && c != nil {
				ahciCtls = append(ahciCtls, c)
			}
		case equal("virtio-pci"):
			c, err := shapeVirtController(bspec)
			if err == nil && c != nil {
				virtCtls = append(virtCtls, c)
			}
		case prefix("mpt", "megaraid", "hpvsa", "hpsa"):
			c, err := shapeRAIDController(bspec, cli.getBool("noraidcli"))
			if err == nil && c != nil {
				raidCtls = append(raidCtls, c)
			}
		default:
			log.Warnf("unsupported mass storage controller found: %s (driver: %s)", bspec.LongName(), bspec.Driver)
		}

		if err != nil {
			log.Warn(err.Error())
		}

	}

	if len(nvmeCtls) != 0 {
		r.Storage.NVMeControllers = nvmeCtls
	}
	if len(raidCtls) != 0 {
		r.Storage.RAIDControllers = raidCtls
	}
	if len(ahciCtls) != 0 {
		r.Storage.AHCIControllers = ahciCtls
	}
	if len(virtCtls) != 0 {
		r.Storage.VirtControllers = virtCtls
	}
	if len(nstdCtls) != 0 {
		r.Storage.NonStdControllers = nstdCtls
	}
}

func shapeAHCIController(bspec *model.PCIBaseSpec) (*model.AHCIController, error) {
	var err error

	ahcid := ahci.NewDecoder(bspec.Path)
	err = ahcid.Decode()
	if err != nil {
		return nil, err
	}

	ctl := new(model.AHCIController)
	ctl.PCIBaseSpec = *bspec

	for _, disk := range ahcid.Disks {
		drv := new(model.Drive)
		drv.Name = disk.Name
		drv.Model = disk.Model
		drv.Size = disk.Size()
		drv.Blocks = disk.Blocks
		drv.LogBlockSize = disk.LogBlockSize
		drv.PhyBlockSize = disk.PhyBlockSize
		drv.Scheduler = disk.Scheduler
		drv.Driver = disk.Driver

		shapeACSDrive(drv)

		ctl.Drives = append(ctl.Drives, drv)
	}

	sort.Slice(ctl.Drives, func(i, j int) bool {
		return util.BlkLabelAscSorter(ctl.Drives[i].Name, ctl.Drives[j].Name)
	})

	return ctl, nil
}

func shapeACSDrive(drv *model.Drive) {
	acsd := acs.NewDecoder("/dev/" + drv.Name)
	err := acsd.Decode()
	if err != nil {
		log.Warnf("spc-acs: %s", err)
		return
	}

	drv.CurTemp = acsd.CurTemp
	drv.MaxTemp = acsd.MaxTemp
	drv.MinTemp = acsd.MinTemp
	drv.FormFactor = acsd.FormFactor
	drv.Rotation = acsd.Rotation
	drv.Firmware = acsd.FirmwareRevision
	drv.SerialNumber = acsd.SerialNumber
	drv.Model = acsd.ModelNumber
	drv.Transport = acsd.Transport
	drv.ByteRead = uint64(acsd.TotalLBARead) * uint64(drv.LogBlockSize)
	drv.ByteWritten = uint64(acsd.TotalLBAWritten) * uint64(drv.LogBlockSize)
	drv.NegSpeed = acsd.NegSpeed
	drv.SigSpeed = acsd.SigSpeed
	drv.PowerCycleCount = acsd.PowerCycleCount
	drv.PowerOnHours = acsd.PowerOnHours
	drv.UnsafeShutdownCount = acsd.UnsafeShutdownCount
	drv.SelfTest = acsd.SelfTestSupport
	drv.ErrorLogging = acsd.ErrorLoggingSupport

	for _, rec := range acsd.ErrorRecords {
		e := new(model.SMARTRecord)
		e.ID = rec.ID
		e.Current = rec.Current
		e.Worst = rec.Worst
		e.Raw = rec.Raw
		e.Threshold = rec.Threshold
		e.Name = rec.Name
		drv.ErrorRecords = append(drv.ErrorRecords, e)
	}
}

func shapeRAIDController(bspec *model.PCIBaseSpec, noRaidCli bool) (*model.RAIDController, error) {
	var err error

	raidd := raid.NewDecoder(bspec.Path)
	err = raidd.Decode()
	if err != nil {
		return nil, err
	}

	ctl := new(model.RAIDController)
	ctl.PCIBaseSpec = *bspec
	ctl.Battery = model.BatteryUnknown

	for _, disk := range raidd.VirtualDisks {
		drv := new(model.LogDrive)
		drv.Name = disk.Name
		drv.Model = disk.Model
		drv.Size = disk.Size()
		drv.Blocks = disk.Blocks
		drv.LogBlockSize = disk.LogBlockSize
		drv.PhyBlockSize = disk.PhyBlockSize
		drv.Scheduler = disk.Scheduler
		drv.Driver = disk.Driver
		drv.SCSIHost = disk.Host
		drv.SCSIChannel = disk.Channel
		drv.SCSITarget = disk.Target
		drv.SCSILun = disk.Lun
		drv.SASAddress = disk.SASAddress
		drv.WWN = disk.WWN

		ctl.LogDrives = append(ctl.LogDrives, drv)
	}

	sort.Slice(ctl.LogDrives, func(i, j int) bool {
		return util.BlkLabelAscSorter(ctl.LogDrives[i].Name, ctl.LogDrives[j].Name)
	})

	prefix := func(t string, ts ...string) bool {
		return matcher(func(l string) bool {
			if strings.HasPrefix(ctl.Driver, l) {
				return true
			}
			return false
		}, t, ts...)
	}
	if !noRaidCli {
		switch {
		case prefix("megaraid"):
			shapeMegaRAIDController(ctl)
		case prefix("mpt"):
			shapeMPTRAIDController(ctl)
		case prefix("hpvsa", "hpsa"):
			shapeHPSARAIDController(ctl)
		}
	}
	return ctl, nil
}

func shapeMegaRAIDController(ctl *model.RAIDController) {
	var err error

	if !megacli.Available() {
		log.Info("megacli is not installed")
		return
	}

	mctls, err := megacli.GetControllers()
	if err != nil {
		log.Debugf("do not parse megaraid due to %s", err)
		return
	}

	var mctl *megacli.Controller
	for _, m := range mctls {
		loc := ctl.Location
		if m.Bus == loc.Bus && m.Device == loc.Device && m.Function == loc.Function {
			mctl = m
			break
		}
	}
	if mctl == nil {
		log.Debug("the controller not found in megacli")
		return
	}

	err = mctl.Decode()
	if err != nil {
		log.Debugf("failed to decode megacli due to %s", err)
		return
	}

	if mctl.Battery {
		ctl.Battery = model.BatteryPresent
	} else {
		ctl.Battery = model.BatteryNotPresent
	}

	ctl.ProductName = mctl.ProductName
	ctl.BIOS = mctl.BIOS
	ctl.Firmware = mctl.Firmware
	ctl.SerialNumber = mctl.SerialNumber
	ctl.AdapterID = fmt.Sprintf("%d", mctl.Number)

	log.Debug("scanning log drives")

	// ctl.LogDrives may contain pass-through drives that must be separated to ctl.PassthroughDrives.
	ldrvs := ctl.LogDrives
	ctl.LogDrives = []*model.LogDrive{}
	for _, ldrv := range ldrvs {
		log.Debugf("scanning megacli data for %s", ldrv.Summary())
		log.Debugf("scsi address: %s", ldrv.SCSIAddress())

		log.Debugf("wwn: %s", ldrv.WWN)
		if ldrv.WWN != "" {
			log.Debug("this logical drive has wwn. possibly it is pass-through disk (jbod)")

			log.Debug("scanning wwn")
			ptpd := mctl.GetPTPhyDriveByWWN(ldrv.WWN)
			if ptpd != nil {
				p := shapeMegaRAIDPhyDisk(ptpd, mctl)
				p.Name = ldrv.Name
				p.Scheduler = ldrv.Scheduler
				p.Driver = ldrv.Driver
				p.SCSIHost = ldrv.SCSIHost
				p.SCSIChannel = ldrv.SCSIChannel
				p.SCSITarget = ldrv.SCSITarget
				p.SCSILun = ldrv.SCSILun
				p.SASAddress = ldrv.SASAddress
				p.WWN = ldrv.WWN
				p.Blocks = ldrv.Blocks
				p.PhyBlockSize = ldrv.PhyBlockSize
				p.LogBlockSize = ldrv.LogBlockSize
				ctl.PassthroughDrives = append(ctl.PassthroughDrives, p)

				log.Debugf("[enc:slt] = %s:%s", p.Enclosure, p.Slot)
				continue
			}

			log.Debug("cound not find wwn from controller")
			log.Debug("continue scanning")
		}

		ld := mctl.GetLogDriveByTarget(ldrv.SCSITarget)
		if ld == nil {
			continue
		}

		log.Debugf("megacli has the ld data for the target %d", ldrv.SCSITarget)
		ldrv.RAIDLv = string(ld.RAIDLv)
		ldrv.CachePolicy = ld.CachePolicy
		ldrv.Status = ld.State
		ldrv.StripeSize = ld.StripSize
		ldrv.GroupLabel = ld.Label

		if _, ok := interface{}(ld).(raidcli.HealthReporter); ok {
			ldrv.Degraded = !ld.IsHealthy()
		}

		for _, pd := range ld.PhyDrives {
			p := shapeMegaRAIDPhyDisk(pd, mctl)
			log.Debugf("found phy drive: %s", p.Model)
			ldrv.PhyDrives = append(ldrv.PhyDrives, p)
		}
		ctl.LogDrives = append(ctl.LogDrives, ldrv)
	}

	for _, pd := range mctl.UnconfDrives {
		p := shapeMegaRAIDPhyDisk(pd, mctl)
		log.Debugf("found unconfigured phy drive: %s", p.Model)
		ctl.UnconfDrives = append(ctl.UnconfDrives, p)
	}
}

func shapeMegaRAIDPhyDisk(pd *megacli.PhyDrive, ctl *megacli.Controller) *model.PhyDrive {
	p := new(model.PhyDrive)
	p.Enclosure = pd.EnclosureID
	p.Slot = pd.SlotNumber
	p.Model = pd.Model
	p.Size = pd.Size
	p.Status = pd.State
	p.Firmware = pd.FirmwareRevision
	p.NegSpeed = pd.DriveSpeed
	p.Transport = pd.Type
	p.SolidStateDrive = pd.SolidStateDrive
	p.ErrorCount = pd.MediaErrorCount

	d := megaraid.NewDecoder(ctl.Number, int(pd.DeviceID), spc.CastDiskType(pd.Type))
	err := d.Decode()
	if err != nil {
		log.Debugf("failed to decode spc/megaraid due to %s", err)
		return p
	}

	blockSize := uint64(pd.LogBlockSize)
	if pd.Type == "SAS" {
		blockSize = 1
	}

	p.ByteWritten = uint64(d.TotalLBAWritten) * blockSize
	p.ByteRead = uint64(d.TotalLBARead) * blockSize
	p.Firmware = d.FirmwareRevision
	p.SerialNumber = d.SerialNumber
	p.Model = d.ModelNumber

	p.PowerCycleCount = d.PowerCycleCount
	p.PowerOnHours = d.PowerOnHours
	p.UnsafeShutdownCount = d.UnsafeShutdownCount

	for _, rec := range d.ErrorRecords {
		e := new(model.SMARTRecord)
		e.ID = rec.ID
		e.Current = rec.Current
		e.Worst = rec.Worst
		e.Raw = rec.Raw
		e.Threshold = rec.Threshold
		e.Name = rec.Name
		p.ErrorRecords = append(p.ErrorRecords, e)
	}

	return p
}

func shapeMPTRAIDController(ctl *model.RAIDController) {
	var err error

	if !sas3ircu.Available() {
		log.Info("sas3ircu is not installed")
		return
	}

	sctls, err := sas3ircu.GetControllers()
	if err != nil {
		log.Debugf("do not parse sas3ircu due to %s", err)
		return
	}

	var sctl *sas3ircu.Controller
	for _, s := range sctls {
		loc := ctl.Location
		if s.Bus == loc.Bus && s.Device == loc.Device && s.Function == loc.Function {
			sctl = s
			break
		}
	}
	if sctl == nil {
		log.Debug("the controller not found in sas3ircu")
		return
	}

	err = sctl.Decode()
	if err != nil {
		log.Debugf("failed to decode sas3ircu due to %s", err)
		return
	}

	ctl.ProductName = sctl.ProductName
	ctl.BIOS = sctl.BIOS
	ctl.Firmware = sctl.Firmware
	ctl.AdapterID = fmt.Sprintf("%d", sctl.Number)

	log.Debug("scanning log drives")

	// ctl.LogDrives may contain pass-through drives that must be separated to ctl.PassthroughDrives.
	ldrvs := ctl.LogDrives
	ctl.LogDrives = []*model.LogDrive{}
	for _, ldrv := range ldrvs {
		log.Debugf("scanning sas3ircu data for %s", ldrv.Summary())

		if ldrv.SASAddress == "" {
			log.Debugf("this device has no sas_address which is required to bind it with sas3ircu outputs")
			continue
		}

		log.Debugf("this device has sas address: %s", ldrv.SASAddress)

		log.Debug("scanning from sas3ircu pass-through drives")
		ptpd := sctl.GetPTDrive(ldrv.SASAddress)
		if ptpd != nil {
			log.Debug("this device is pass-through drive")
			p := shapeMPTRAIDPhyDisk(ptpd)
			p.Name = ldrv.Name
			p.Scheduler = ldrv.Scheduler
			p.Driver = ldrv.Driver
			p.SCSIHost = ldrv.SCSIHost
			p.SCSIChannel = ldrv.SCSIChannel
			p.SCSITarget = ldrv.SCSITarget
			p.SCSILun = ldrv.SCSILun
			p.SASAddress = ldrv.SASAddress
			p.WWN = ldrv.WWN
			p.Blocks = ldrv.Blocks
			p.PhyBlockSize = ldrv.PhyBlockSize
			p.LogBlockSize = ldrv.LogBlockSize
			ctl.PassthroughDrives = append(ctl.PassthroughDrives, p)

			log.Debugf("[enc:slt] = %s:%s", p.Enclosure, p.Slot)
			continue
		}
		log.Debug("not found")

		log.Debug("scanning from sas3ircu log drives")
		ld := sctl.GetLogDrive(ldrv.SASAddress)
		if ld == nil {
			log.Debug("not found")
			continue
		}
		ldrv.RAIDLv = string(ld.RAIDLv)
		ldrv.Status = ld.State
		ldrv.GroupLabel = ld.Label

		if _, ok := interface{}(ld).(raidcli.HealthReporter); ok {
			ldrv.Degraded = !ld.IsHealthy()
		}

		for _, pd := range ld.PhyDrives {
			p := shapeMPTRAIDPhyDisk(pd)
			log.Debugf("found phy drive: %s", p.Model)
			ldrv.PhyDrives = append(ldrv.PhyDrives, p)
		}
		ctl.LogDrives = append(ctl.LogDrives, ldrv)
	}
	return
}

func shapeMPTRAIDPhyDisk(pd *sas3ircu.PhyDrive) *model.PhyDrive {
	p := new(model.PhyDrive)
	p.Enclosure = pd.EnclosureID
	p.Slot = pd.SlotNumber
	p.Model = pd.Model
	p.Size = pd.Size
	p.Status = pd.State
	p.Firmware = pd.Firmware
	p.SerialNumber = pd.SerialNumber
	p.Transport = pd.Protocol
	p.SolidStateDrive = pd.SolidStateDrive
	return p
}

func shapeHPSARAIDController(ctl *model.RAIDController) {
	var err error

	if !hpacucli.Available() {
		log.Info("hpssacli is not installed")
		return
	}

	hctls, err := hpacucli.GetControllers()
	if err != nil {
		log.Debugf("do not parse hpacucli due to %s", err)
		return
	}

	pciaddr := fmt.Sprintf("%04x:%02x:%02x.%x", ctl.Location.Domain, ctl.Location.Bus, ctl.Location.Device, ctl.Location.Function)

	var hctl *hpacucli.Controller
	for _, h := range hctls {
		err := h.Decode()
		if err != nil {
			log.Debugf("failed to decode hpacucli due to %s", err)
			return
		}

		if h.PCIAddr == pciaddr {
			hctl = h
			break
		}
	}

	if hctl == nil {
		log.Debug("the controller not found in hpacucli")
		return
	}

	ctl.ProductName = hctl.ProductName
	ctl.Firmware = hctl.Firmware
	ctl.SerialNumber = hctl.SerialNumber
	ctl.AdapterID = hctl.Slot

	if hctl.Battery {
		ctl.Battery = model.BatteryPresent
	} else {
		ctl.Battery = model.BatteryNotPresent
	}

	log.Debug("scanning log drives")
	for _, ldrv := range ctl.LogDrives {
		log.Debugf("scanning hpacucli for %s", ldrv.Summary())

		ld := hctl.GetLogDrive("/dev/" + ldrv.Name)
		if ld == nil {
			log.Debug("not found")
			continue
		}

		ldrv.RAIDLv = string(ld.RAIDLv)
		ldrv.Status = ld.State
		ldrv.GroupLabel = ld.Label

		if _, ok := interface{}(ld).(raidcli.HealthReporter); ok {
			ldrv.Degraded = !ld.IsHealthy()
		}

		for _, pd := range ld.PhyDrives {
			p := shapeHPSARAIDPhyDisk(pd)
			log.Debugf("found phy drive: %s", p.Model)
			ldrv.PhyDrives = append(ldrv.PhyDrives, p)
		}
	}

	for _, pd := range hctl.UnconfDrives {
		p := shapeHPSARAIDPhyDisk(pd)
		log.Debugf("found unconfigured phy drive: %s", p.Model)
		ctl.UnconfDrives = append(ctl.UnconfDrives, p)
	}

	return
}

func shapeHPSARAIDPhyDisk(pd *hpacucli.PhyDrive) *model.PhyDrive {
	p := new(model.PhyDrive)
	p.Enclosure = pd.Box
	p.Slot = pd.Bay
	p.Model = pd.Model
	p.Size = pd.Size
	p.Status = pd.Status
	p.Firmware = pd.Firmware
	p.SerialNumber = pd.SerialNumber
	p.Transport = pd.Protocol
	if pd.Rotation == 0 {
		p.SolidStateDrive = true
	}
	return p
}

func shapeNVMeController(bspec *model.PCIBaseSpec) (*model.NVMeController, error) {
	var err error

	nvmed := nvme.NewDecoder(bspec.Path)
	err = nvmed.Decode()
	if err != nil {
		return nil, err
	}

	ctl := new(model.NVMeController)
	ctl.Name = nvmed.Name
	ctl.PCIBaseSpec = *bspec

	admd := nvmeadm.NewDecoder("/dev/" + ctl.Name)
	err = admd.Decode()
	if err == nil {
		ctl.CurTemp = admd.CurTemp
		ctl.WarnTemp = admd.WarnTemp
		ctl.CritTemp = admd.CritTemp
		ctl.SerialNumber = admd.SerialNumber
		ctl.Model = admd.ModelNumber
		ctl.Firmware = admd.FirmwareRevision
		ctl.ByteRead = admd.ByteRead
		ctl.ByteWritten = admd.ByteWritten
		ctl.Size = admd.Size
		ctl.PowerCycleCount = admd.PowerCycleCount
		ctl.PowerOnHours = admd.PowerOnHours
		ctl.UnsafeShutdownCount = admd.UnsafeShutdownCount
	}

	for _, n := range nvmed.Namespaces {
		ns := new(model.Namespace)
		ns.Name = n.Name

		sz := admd.GetNamespaceSize(n.ID())
		if sz > 0 {
			ns.Size = sz
		} else {
			ns.Size = n.Size()
		}

		ns.Scheduler = n.Scheduler
		ns.PhyBlockSize = n.PhyBlockSize
		ns.LogBlockSize = n.LogBlockSize

		ctl.Namespaces = append(ctl.Namespaces, ns)
	}

	sort.Slice(ctl.Namespaces, func(i, j int) bool {
		return util.BlkLabelAscSorter(ctl.Namespaces[i].Name, ctl.Namespaces[j].Name)
	})

	return ctl, nil
}

func shapeVirtController(bspec *model.PCIBaseSpec) (*model.VirtController, error) {
	var err error

	virtd := virtio.NewDecoder(bspec.Path)
	err = virtd.Decode()
	if err != nil {
		return nil, err
	}

	ctl := new(model.VirtController)
	ctl.PCIBaseSpec = *bspec

	for _, disk := range virtd.Disks {
		drv := new(model.Drive)
		drv.Name = disk.Name
		drv.Model = disk.Model
		drv.Size = disk.Size()
		drv.Blocks = disk.Blocks
		drv.LogBlockSize = disk.LogBlockSize
		drv.PhyBlockSize = disk.PhyBlockSize
		drv.Scheduler = disk.Scheduler
		drv.Driver = disk.Driver
		ctl.Drives = append(ctl.Drives, drv)
	}

	sort.Slice(ctl.Drives, func(i, j int) bool {
		return util.BlkLabelAscSorter(ctl.Drives[i].Name, ctl.Drives[j].Name)
	})

	return ctl, nil
}
