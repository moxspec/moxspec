package pci

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/moxspec/moxspec/util"
)

// NewDevice creates and initializes a new Device using contents inside of given path
func NewDevice(path string) *Device {
	d := new(Device)
	d.Path = path
	return d
}

// Decode makes Device satisfy the mox.Decoder interface
func (d *Device) Decode() error {
	var err error

	log.Debugf("parsing %s", d.Path)

	d.Domain, d.Bus, d.Device, d.Function, err = ParseLocater(d.Path)
	if err != nil {
		return err
	}
	log.Debugf("dom: %04x bus: %02x dev: %02x fun:%x", d.Domain, d.Bus, d.Device, d.Function)

	err = parseConfig(d, filepath.Join(d.Path, "config"))
	if err != nil {
		return err
	}

	d.VendorName = vdb.vendorName(d.VendorID)
	d.DeviceName = vdb.deviceName(d.VendorID, d.DeviceID)
	d.SubSystemName = vdb.subSysName(d.VendorID, d.DeviceID, d.SubSystemVendorID, d.SubSystemDeviceID)
	d.ClassName = cdb.className(d.ClassID)
	d.SubClassName = cdb.subClassName(d.ClassID, d.SubClassID)
	d.InterfaceName = cdb.intfaceName(d.ClassID, d.SubClassID, d.InterfaceID)

	log.Debugf("ven: %s dev: %s subsys: %s class: %s", d.VendorName, d.DeviceName, d.SubSystemName, d.ClassName)

	d.Numa, _ = util.LoadUint16(filepath.Join(d.Path, "numa_node"))
	drvDir := filepath.Join(d.Path, "driver")
	if util.Exists(drvDir) {
		linkPath, err := os.Readlink(drvDir)
		if err == nil {
			d.Driver = filepath.Base(linkPath)
		}
	}

	return nil
}

// ParseLocater parses the pci identifier string
func ParseLocater(basePath string) (dom, bus, dev, fun uint32, err error) {
	// e.g: 0000:d7:12.0
	l := strings.Split(filepath.Base(basePath), ":")
	if len(l) != 3 {
		err = fmt.Errorf("the path is invalid, %s", basePath)
		return
	}

	ll := strings.Split(l[2], ".")
	if len(ll) != 2 {
		err = fmt.Errorf("the path has no function number, %s", basePath)
		return
	}

	do, err := strconv.ParseUint(l[0], 16, 32)
	if err != nil {
		return
	}
	dom = uint32(do)

	b, err := strconv.ParseUint(l[1], 16, 32)
	if err != nil {
		return
	}
	bus = uint32(b)

	de, err := strconv.ParseUint(ll[0], 16, 32)
	if err != nil {
		return
	}
	dev = uint32(de)

	f, err := strconv.ParseUint(ll[1], 16, 32)
	if err != nil {
		return
	}
	fun = uint32(f)

	return
}
