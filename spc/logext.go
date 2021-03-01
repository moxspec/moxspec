package spc

import (
	"fmt"
	"os"

	"github.com/moxspec/moxspec/util"
)

type logExtData struct {
	powerCycleCount uint64
	powerOnHours    uint64
	totalLBAWritten uint64
	totalLBARead    uint64
}

func (d *Device) decodeLogExt(fd *os.File) error {
	// ACS-4
	// 7.22 READ LOG EXT – 2Fh, PIO Data-In
	// 7.22.3 Inputs
	// http://www.t13.org/documents/UploadedDocuments/docs2016/di529r14-ATAATAPI_Command_Set_-_4.pdf
	//
	// 9.5 Device Statistics log (Log Address 04h)
	// 9.5.4 General Statistics (log page 01h)
	cdb := makeATAPThruCmd(0x00, 0x04, 0x01, 0x00, 0x2F)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return err
	}

	l, err := parseLogExt(buf)
	if err != nil {
		return err
	}

	if l.powerCycleCount > 0 && l.powerCycleCount > d.PowerCycleCount {
		log.Debug("set power cycle count from gplog 0x04")
		d.PowerCycleCount = l.powerCycleCount
	}

	if l.powerOnHours > 0 && l.powerOnHours > d.PowerOnHours {
		log.Debug("set power on hours from gplog 0x04")
		d.PowerOnHours = l.powerOnHours
	}

	if l.totalLBAWritten > 0 && l.totalLBAWritten > d.TotalLBAWritten {
		log.Debug("set total lba written from gplog 0x04")
		d.TotalLBAWritten = l.totalLBAWritten
	}

	if l.totalLBARead > 0 && l.totalLBARead > d.TotalLBARead {
		log.Debug("set total lba read from gplog 0x04")
		d.TotalLBARead = l.totalLBARead
	}

	return nil
}

func parseLogExt(buf []byte) (*logExtData, error) {
	if 48 > len(buf) {
		return nil, fmt.Errorf("failed to read log ext due to invalid length data given")
	}

	l := new(logExtData)
	// Table 222 General Statistics (part 1 of 3)
	l.powerCycleCount = util.BytesToUint64(buf[8:15])  // NOTE: last 1-byte is DEVICE STATISTICS FLAGS field
	l.powerOnHours = util.BytesToUint64(buf[16:23])    // NOTE: last 1-byte is DEVICE STATISTICS FLAGS field
	l.totalLBAWritten = util.BytesToUint64(buf[24:31]) // NOTE: last 1-byte is DEVICE STATISTICS FLAGS field
	l.totalLBARead = util.BytesToUint64(buf[40:47])    // NOTE: last 1-byte is DEVICE STATISTICS FLAGS field

	return l, nil
}
