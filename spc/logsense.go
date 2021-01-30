package spc

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 5.2.21 Supported Log Pages log page (00h/00h)
	logSenseSupportLogPage = 0x00
	// Table 299 Error counter log page codes
	logSenseWriteErrorPage = 0x02
	logSenseReadErrorPage  = 0x03
	// Table 66 Page Control (PC) field values
	pcCumulative = 0x40
	// Table 300 Parameter codes for error counter log pages
	logSenseErrorPageTotalByte = 0x05
)

func (d *Device) decodeLogSense(fd *os.File) error {
	supportPages, err := d.readLogSenseSupportPages(fd)
	if err != nil {
		return err
	}

	if _, supported := supportPages[logSenseReadErrorPage]; supported {
		params, err := d.readLogSenseReadError(fd)
		if err != nil {
			return err
		}

		if totalBytes, ok := params[logSenseErrorPageTotalByte]; ok {
			d.TotalLBARead = totalBytes
		}
	}

	if _, supported := supportPages[logSenseWriteErrorPage]; supported {
		params, err := d.readLogSenseWriteError(fd)
		if err != nil {
			return err
		}

		if totalBytes, ok := params[logSenseErrorPageTotalByte]; ok {
			d.TotalLBAWritten = totalBytes
		}
	}

	return nil
}

func (d *Device) readLogSenseSupportPages(fd *os.File) (map[byte]struct{}, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 5.2.21 Supported Log Pages log page (00h/00h)
	cdb := makeReadLogSenseCmd(pcCumulative+logSenseSupportLogPage, 0x00, 0x16)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return nil, err
	}

	return parseReadLogSenseSupportPages(buf)
}

func parseReadLogSenseSupportPages(buf []byte) (map[byte]struct{}, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 5.2.21 Supported Log Pages log page (00h/00h)

	if len(buf) < 4 {
		return nil, fmt.Errorf("read logsense support pages: too short response")
	}

	// first 2 octest: page code and subpage code
	offset := 2
	length := int(binary.BigEndian.Uint16(buf[offset : offset+2]))
	offset += 2

	if len(buf) < length+3 {
		return nil, fmt.Errorf("read logsense support pages: resp is too shorter than page length")
	}

	supportPages := make(map[byte]struct{})
	for i := 0; i < length; i++ {
		supportPages[buf[offset+i]] = struct{}{}
	}

	return supportPages, nil
}

func (d *Device) readLogSenseReadError(fd *os.File) (map[int]uint64, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// Table 299 Error counter log page codes
	cdb := makeReadLogSenseCmd(pcCumulative+logSenseReadErrorPage, 0x00, 0x65)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return nil, err
	}

	return parseLogSenseErrorCounterPages(buf)
}

func (d *Device) readLogSenseWriteError(fd *os.File) (map[int]uint64, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// Table 299 Error counter log page codes
	cdb := makeReadLogSenseCmd(pcCumulative+logSenseWriteErrorPage, 0x00, 0x65)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return nil, err
	}

	return parseLogSenseErrorCounterPages(buf)
}

func parseLogSenseErrorCounterPages(buf []byte) (map[int]uint64, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 5.2.9 Error counter log pages (WRITE, READ, and VERIFY)

	if len(buf) < 4 {
		return nil, fmt.Errorf("read logsense error counter page: too short response")
	}

	// first 2 octest: page code and subpage code
	off := 2
	// Page Length = frame size - 4
	pageLength := int(binary.BigEndian.Uint16(buf[off:off+2])) + 4
	off += 2

	if len(buf) < pageLength {
		return nil, fmt.Errorf("read logsense error counter page: too short response")
	}

	params := map[int]uint64{}

	for off < pageLength {
		code := int(binary.BigEndian.Uint16(buf[off : off+2]))
		// skip param control byte (1octest) + seek param code (2octets)
		off += 3
		length := int(buf[off])
		off++

		if off+length > pageLength {
			break
		}

		value := parseLogSenseErrorCounterParamValue(buf[off : off+length])
		off += length

		params[code] = value
	}

	return params, nil
}

func parseLogSenseErrorCounterParamValue(buf []byte) uint64 {
	// if len > sizeof(uint64), return 0 like as sg3_utils v1.44
	// sg3_utils/src/sg_logs.c / show_error_counter_page
	// sg3_utils/include/sg_unaligned.h / sg_get_unaligned_be
	if len(buf) > 8 {
		return 0
	}

	padding := make([]byte, 8-len(buf))
	return binary.BigEndian.Uint64(append(padding, buf...))
}
