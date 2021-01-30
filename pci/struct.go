package pci

import (
	"encoding/binary"
	"fmt"
)

// Device represents status of a PCI or PCIe device status
type Device struct {
	Path              string
	DeviceID          uint16
	DeviceName        string
	Revision          byte
	VendorID          uint16
	VendorName        string
	SubSystemVendorID uint16
	SubSystemDeviceID uint16
	SubSystemName     string
	ClassID           byte
	ClassName         string
	SubClassID        byte
	SubClassName      string
	InterfaceID       byte
	InterfaceName     string
	Driver            string
	Domain            uint32
	Bus               uint32
	Device            uint32
	Function          uint32
	Numa              uint16
	MaxGen            byte
	MaxSpeed          float32
	MaxWidth          byte
	LinkGen           byte
	LinkSpeed         float32
	LinkWidth         byte
	SlotPowetLimit    float32
	SerialNumber      string
	UncorrectableErrs []string
	CorrectableErrs   []string
	Express           bool // indicates PCIe
	BasicCaps         []*BasicCap
	ExtCaps           []*ExtCap
}

// PCIID returns the device's identifier as Domain:Bus:Device:FunctionNumber
func (d Device) PCIID() string {
	return IDString(d.Domain, d.Bus, d.Device, d.Function)
}

// IDString returns given Domain, Bus, Device, FunctionNumber as string
func IDString(dom, bus, dev, fun uint32) string {
	return fmt.Sprintf("%04x:%02x:%02x.%x", dom, bus, dev, fun)
}

// IsUnknown returns if the device is unknown
func (d Device) IsUnknown() bool {
	return (d.VendorName == "" || d.DeviceName == "")
}

// Config represents a PCI configuration register
type Config struct {
	path string
	br   []byte
}

// ReadByteFrom reads one-byte as byte value from the given offset
func (c Config) ReadByteFrom(offset uint16) byte {
	if offset < uint16(len(c.br)) {
		return c.br[offset]
	}
	return 0
}

// ReadWordFrom reads two-byte as uint16 value from the given offset
func (c Config) ReadWordFrom(offset uint16) uint16 {
	if offset+1 < uint16(len(c.br)) {
		return binary.LittleEndian.Uint16(c.br[offset : offset+2])
	}
	return 0
}

// ReadDWordFrom reads four-byte as uint32 value from the given offset
func (c Config) ReadDWordFrom(offset uint16) uint32 {
	if offset+3 < uint16(len(c.br)) {
		l := c.ReadWordFrom(offset)
		u := c.ReadWordFrom(offset + 2)
		return (uint32(u)<<16 | uint32(l))
	}
	return 0
}

// BasicCap represents a PCI capability
type BasicCap struct {
	Offset uint16
	ID     byte
	Next   uint16
}

// ExtCap represents a PCIe capability
type ExtCap struct {
	Offset uint16
	ID     uint16
	Ver    byte
	Next   uint16
}

// Devices represents the database of devices
type Devices struct {
	all     []*Device
	classes map[byte][]*Device
}

func (db *Devices) append(dev *Device) {
	db.all = append(db.all, dev)
	db.classes[dev.ClassID] = append(db.classes[dev.ClassID], dev)
}

// FilterByClass returns a slice of device filtered by given class id
func (db Devices) FilterByClass(cids ...byte) []*Device {
	var devs []*Device

	for _, cid := range cids {
		ds, ok := db.classes[cid]
		if !ok {
			continue
		}

		devs = append(devs, ds...)
	}

	return devs
}

// AllDevices returns a slice of all devices
func (db Devices) AllDevices() []*Device {
	return db.all
}
