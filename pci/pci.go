package pci

import (
	"io/ioutil"
	"path/filepath"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

var pciidsPossible = []string{
	"/etc/mox/pci.ids",
	"/usr/share/hwdata/pci.ids",
	"/usr/share/misc/pci.ids",
}

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("pci")
}

// NewDecoder creates and initializes a new slice of Device as Decoder
func NewDecoder() *Devices {
	ds := new(Devices)
	ds.classes = make(map[byte][]*Device)
	return ds
}

// Decode makes Devices satisfy the mox.Decoder interface
func (devs *Devices) Decode() error {
	var err error
	var pciids string
	for _, possible := range pciidsPossible {
		if util.Exists(possible) {
			pciids = possible
			log.Debugf("found pciids: %s", pciids)
			break
		}
	}
	err = initDB(pciids)
	if err != nil {
		return err
	}

	// TODO: accessing via /sys/bus is DEPRECATED, to be fixed to use /sys/class/pci_bus
	syspath := "/sys/bus/pci/devices"
	dirs, err := ioutil.ReadDir(syspath)
	if err != nil {
		return err
	}

	var oldDB bool
	for _, d := range dirs {
		bpath, err := filepath.EvalSymlinks(filepath.Join(syspath, d.Name()))
		if err != nil {
			log.Warnf("could not read %s (%s)", bpath, err)
		}

		dev := NewDevice(bpath)
		err = dev.Decode()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if !checkDBFreshness(dev) {
			oldDB = true
		}

		devs.append(dev)
	}

	if oldDB {
		log.Debug("failed to decode a pci (vendor|device|class) name completely")
		log.Debug("the pci database possibly be out of date")
	}

	return nil
}

func checkDBFreshness(d *Device) bool {
	fresh := true

	if d.VendorName == "" {
		fresh = false
		log.Debugf("pci: cound not decode vendor id %x", d.VendorID)
	}
	if d.DeviceName == "" {
		fresh = false
		log.Debugf("pci: cound not decode device id %d", d.DeviceID)
	}
	if d.SubSystemName == "" {
		fresh = false
		log.Debugf("pci: cound not decode sub system %x:%x", d.SubSystemVendorID, d.SubSystemDeviceID)
	}
	if d.ClassName == "" {
		fresh = false
		log.Debugf("pci: cound not decode pci class id %x", d.ClassID)
	}

	if d.SubClassID > 0 {
		if d.SubClassName == "" {
			fresh = false
			log.Debugf("pci: cound not decode pci sub class id %x", d.SubClassID)
		}

		if d.InterfaceID > 0 {
			if d.InterfaceName == "" {
				fresh = false
				log.Debugf("pci: cound not decode pci interface name: %d", d.InterfaceID)
			}
		}
	}
	return fresh
}
