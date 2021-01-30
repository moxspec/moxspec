package spc

import (
	"fmt"
	"os"

	"github.com/actapio/moxspec/loglet"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("spc")
}

// Decode shapes SCSI/ATA device via mox.Decoder interface
func (d *Device) Decode() error {

	if d.ioctlDeviceFilePath == "" {
		return fmt.Errorf("ioctl device is empty")
	}

	fd, err := os.OpenFile(d.ioctlDeviceFilePath, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = d.decodeInquiry(fd)
	if err != nil {
		return err
	}

	if d.diskType == SASDisk {
		err = d.decodeLogSense(fd)
		if err != nil {
			return err
		}
	}

	if d.diskType == SATADisk {
		err = d.decodeSmart(fd)
		if err != nil {
			return err
		}

		err = d.decodeLogExt(fd)
		if err != nil {
			return err
		}

		err = d.decodeIdentify(fd)
		if err != nil {
			return err
		}
	}

	return nil
}
