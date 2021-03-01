package nvmeadm

import (
	"os"
	"strings"

	"github.com/moxspec/moxspec/loglet"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("nvmeadm")
}

// NewDecoder creates and initializes a Device as Decoder
func NewDecoder(p string) *Device {
	if !strings.Contains(p, "nvme") {
		return nil
	}

	// TODO: test path
	d := new(Device)
	d.path = p
	return d
}

// Decode makes Device satisfy the mox.Decoder interface
func (d *Device) Decode() error {
	var err error
	fd, err := os.OpenFile(d.path, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return err
	}
	defer fd.Close()

	return decode(fd, d)
}

func decode(fd *os.File, d *Device) error {
	var err error

	err = readSmartLog(fd, d)
	if err != nil {
		return err
	}

	err = identifyController(fd, d)
	if err != nil {
		return err
	}

	return nil
}
