package pci

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	pcieConfigSpaceSize = 4096
	pcieDeviceIDBitSize = 16
	pcieVendorIDBitSize = 16
)

// NewConfig creates and initializes a new Config by decoding config register from given path
func NewConfig(p string) *Config {
	fd, err := os.Open(p)
	if err != nil {
		log.Warnf("could not read %s. %s", p, err)
		return nil
	}
	defer fd.Close()

	c := new(Config)
	c.path = p
	c.br, err = ioutil.ReadAll(io.LimitReader(fd, pcieConfigSpaceSize))
	if err != nil {
		log.Warnf("could not read %s. %s", p, err)
		return nil
	}

	return c
}

func parseConfig(dev *Device, path string) error {
	conf := NewConfig(path)

	dev.VendorID = conf.ReadWordFrom(0x00)
	if dev.VendorID == 0xFFFF {
		log.Debugf("The vendor id read from config regsiter is 0xFFFF. Trying to read from sysfs vendor object.")
		value, err := readSysfsValue(dev.Path, "vendor", pcieVendorIDBitSize)
		if err != nil {
			log.Warnf("cannot parse %s/vendor. %s", dev.Path, err)
			return err
		}
		dev.VendorID = uint16(value)
	}

	dev.DeviceID = conf.ReadWordFrom(0x02)
	if dev.DeviceID == 0xFFFF {
		log.Debugf("The device id read from config regsiter is 0xFFFF. Trying to read from sysfs device object.")
		value, err := readSysfsValue(dev.Path, "device", pcieDeviceIDBitSize)
		if err != nil {
			log.Warnf("cannot parse %s/device. %s", dev.Path, err)
			return err
		}
		dev.DeviceID = uint16(value)
	}

	dev.Revision = conf.ReadByteFrom(0x08)
	dev.InterfaceID = conf.ReadByteFrom(0x09)
	dev.SubClassID = conf.ReadByteFrom(0x0A)
	dev.ClassID = conf.ReadByteFrom(0x0B)
	dev.SubSystemVendorID = conf.ReadWordFrom(0x2C)
	dev.SubSystemDeviceID = conf.ReadWordFrom(0x2C + 2)

	log.Debugf("ven: %04x dev: %04x rev: %02x class: %02x subclass: %02x intf: %02x subven: %04x subdev: %04x",
		dev.VendorID, dev.DeviceID, dev.Revision,
		dev.ClassID, dev.SubClassID, dev.InterfaceID,
		dev.SubSystemVendorID, dev.SubSystemDeviceID,
	)

	if dev.VendorID == 0 {
		log.Debug("empty device")
		return nil
	}

	var err error
	err = parseBasicCapabilities(dev, conf)
	if err != nil {
		return err
	}

	if dev.Express == false {
		log.Debug("this is not PCIe device")
		return nil
	}

	err = parseExtCapabilities(dev, conf)
	if err != nil {
		return err
	}

	return nil
}

func readSysfsValue(path string, object string, bitSize int) (uint64, error) {
	p := filepath.Join(path, object)

	fd, err := os.Open(p)
	if err != nil {
		log.Warnf("cannot open %s. %s", p, err)
		return 0, err
	}
	defer fd.Close()

	// len("0xFFFFFFFFFFFFFFFF") = 18 (64bit hex value)
	buf, err := ioutil.ReadAll(io.LimitReader(fd, 18))
	if err != nil {
		log.Warnf("cannot read %s. %s", p, err)
		return 0, err
	}

	str := strings.Trim(string(buf), "\n")
	value, err := strconv.ParseUint(str, 0, bitSize)
	if err != nil {
		log.Warnf("cannot convert %s to int: %s", str, err)
		return 0, err
	}

	return value, nil
}

func parseBasicCapabilities(dev *Device, conf *Config) error {
	footPrint := make(map[uint16]bool)

	log.Debug("parsing basic capability")
	capPtr := uint16(conf.ReadByteFrom(0x34))
	if capPtr != 0x00 {
		for {
			bc := new(BasicCap)
			bc.Offset = capPtr
			bc.ID = conf.ReadByteFrom(capPtr)
			if bc.ID == 0x10 { // PCI Express capability register
				dev.Express = true

				// device capabbility register
				devCapReg := conf.ReadDWordFrom(capPtr + 0x04)
				dev.SlotPowetLimit = parseSlotPowerLimit(devCapReg)

				// link capability register
				linkCapReg := conf.ReadDWordFrom(capPtr + 0x0C)
				dev.MaxGen, dev.MaxSpeed, dev.MaxWidth = parseLinkSpec(linkCapReg)

				// link status register
				linkStaReg := conf.ReadDWordFrom(capPtr + 0x12)
				dev.LinkGen, dev.LinkSpeed, dev.LinkWidth = parseLinkSpec(linkStaReg)
			}

			bc.Next = uint16(conf.ReadByteFrom(capPtr + 1))
			log.Debugf("ptr=%02x id=%02x next=%02x", bc.Offset, bc.ID, bc.Next)
			dev.BasicCaps = append(dev.BasicCaps, bc)

			footPrint[bc.Offset] = true

			if bc.Next == 0x00 {
				break
			}
			if _, ok := footPrint[bc.Next]; ok {
				log.Error("circular reference was detected")
				break
			}
			capPtr = bc.Next
		}
	}

	log.Debugf("%d basic capabilities", len(dev.BasicCaps))
	return nil
}

func parseExtCapabilities(dev *Device, conf *Config) error {
	footPrint := make(map[uint16]bool)

	log.Debug("parsing extended capability")
	var exCapPtr uint16 = 0x100
	for {
		ec := new(ExtCap)

		dw := conf.ReadDWordFrom(exCapPtr)
		ec.Offset = exCapPtr
		ec.ID = uint16(dw & 0xFFFF)
		ec.Ver = byte((dw >> 16) & 0xF)
		ec.Next = uint16((dw >> 20) & 0xFFF)

		log.Debugf("ptr=%04x id=%04x ver=%02x next=%04x", ec.Offset, ec.ID, ec.Ver, ec.Next)

		switch ec.ID {
		case 0x0001: // Advanced Error Reporting Extended Capability
			// UncorrectableErrorStatusRegister(Offset04h)
			ucStat := conf.ReadDWordFrom(exCapPtr + 0x04)
			parseUncorrectableErrs(ucStat)

			// CorrectableErrorStatusRegister(Offset10h)
			crStat := conf.ReadDWordFrom(exCapPtr + 0x10)
			parseCorrectableErrs(crStat)

		case 0x0003: // Device Serial Number Capability
			// SerialNumberRegister(Offset04h)
			dw1 := conf.ReadDWordFrom(exCapPtr + 0x04) // 1st DW
			dw2 := conf.ReadDWordFrom(exCapPtr + 0x08) // 2nd DW
			dev.SerialNumber = parseSerialNumber(dw1, dw2)
			log.Debugf("serial number: %s", dev.SerialNumber)
		}

		dev.ExtCaps = append(dev.ExtCaps, ec)

		footPrint[ec.Offset] = true

		exCapPtr = ec.Next
		if exCapPtr == 0x00 {
			break
		}
		if _, ok := footPrint[ec.Next]; ok {
			log.Error("circular reference was detected")
			break
		}
	}
	log.Debugf("%d extended capabilities", len(dev.ExtCaps))

	return nil
}

func parseSlotPowerLimit(reg uint32) float32 {
	val := float32((reg >> 18) & 0xFF)
	scale := (reg >> 26) & 0x3
	switch scale {
	case 0x01:
		val = val * 0.1
	case 0x02:
		val = val * 0.01
	case 0x03:
		val = val * 0.001
	}
	return val
}

func parseLinkSpec(reg uint32) (gen byte, speed float32, width byte) {
	gen = byte(reg & 0x1F)
	speeds := []float32{
		0.0,  // 0: unknown
		2.5,  // 1: Gen 1
		5.0,  // 2: Gen 2
		8.0,  // 3: Gen 3
		16.0, // 4: Gen 4
		32.0, // 4: Gen 5
	}
	if gen < byte(len(speeds)) {
		speed = speeds[gen]
	}

	w := (reg >> 4) & 0x3F
	if w <= 32 {
		width = byte(w)
	}

	log.Debugf("gen: %d / speed: %.1f / width: %d", gen, speed, width)
	return
}

func parseSerialNumber(dw1, dw2 uint32) string {
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x-%02x-%02x",
		(dw2>>24)&0xFF, (dw2>>16)&0xFF, (dw2>>8)&0xFF, dw2&0xFF,
		(dw1>>24)&0xFF, (dw1>>16)&0xFF, (dw1>>8)&0xFF, dw1&0xFF,
	)
}

func parseUncorrectableErrs(reg uint32) []string {
	// Table 7-31: Uncorrectable Error Status Register
	ucDefs := map[byte]string{
		4:  "Data Link Protocol Error",
		5:  "Surprise Down Error",
		12: "Poisoned TLP",
		13: "Flow Control Protocol Error",
		14: "Completion Timeout",
		15: "Completer Abort",
		16: "Unexpected Completion",
		17: "Receiver Overflow",
		18: "Malformed TLP",
		19: "ECRC Error",
		20: "Unsupported Request Error",
		21: "ACS Violation",
		22: "Uncorrectable Internal Error",
	}
	return parseBitDefs(reg, ucDefs)
}

func parseCorrectableErrs(reg uint32) []string {
	// Table 7-34: Correctable Error Status Register
	crDefs := map[byte]string{
		0:  "Receiver Error",
		6:  "Bad TLP",
		7:  "Bad DLLP",
		8:  "REPLAY_NUM Rollover",
		12: "Replay Timer Timeout",
		13: "Advisory Non-Fatal Error",
		14: "Corrected Internal Error",
		15: "Header Log Overflow",
	}
	return parseBitDefs(reg, crDefs)
}

func parseBitDefs(reg uint32, defs map[byte]string) []string {
	var stats []string
	var i byte
	for i = 0; i < 32; i++ {
		if reg&(1<<i) == 0 {
			continue
		}
		if stat, ok := defs[i]; ok {
			stats = append(stats, stat)
		}
	}

	return stats
}
