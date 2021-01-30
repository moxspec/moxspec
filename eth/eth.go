package eth

import (
	"strings"

	"github.com/actapio/moxspec/loglet"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("eth")
}

// Device represents the ethernet device
type Device struct {
	Name             string
	Speed            uint32
	Link             bool
	SupportedSpeed   []string
	AdvertisingSpeed []string
	Module           *Module
	FirmwareVersion  string
}

// Decode makes Device satisfy the mox.Decoder interface
func (d *Device) Decode() error {
	var err error

	ndev, err := newNetdev(d.Name)
	if err != nil {
		return err
	}
	defer ndev.close()

	linkup, err := getLinkStat(ndev)
	if err != nil {
		return err
	}
	d.Link = linkup
	log.Debugf("linkup: %v", d.Link)

	lset, err := getLinkSettings(ndev)
	if err != nil {
		return err
	}

	if d.Link {
		d.Speed = lset.speed
		log.Debugf("speed: %d", d.Speed)
	}

	d.SupportedSpeed = lset.supportedSpeed()
	log.Debugf("supported speed: %s", strings.Join(d.SupportedSpeed, " "))

	d.AdvertisingSpeed = lset.advertisingSpeed()
	log.Debugf("advertising speed: %s", strings.Join(d.AdvertisingSpeed, " "))

	if d.Link && lset.port != portTP {
		mod, err := getModule(ndev)
		if err != nil { // not fatal
			log.Debug(err)
		} else {
			d.Module = mod
		}
	}

	driver, err := getDrvInfo(ndev)
	if err != nil {
		return err
	}

	d.FirmwareVersion = driver.firmwareVersion()

	return nil
}

// NewDecoder creates and initializes a Device as Decoder
func NewDecoder(name string) *Device {
	d := new(Device)
	d.Name = name
	return d
}
