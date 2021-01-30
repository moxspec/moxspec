package spc

import (
	"fmt"
	"os"
	"strings"
)

type inquiryData struct {
	vendor   string
	model    string
	revision string
}

func (d *Device) decodeInquiry(fd *os.File) error {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 3.6.1 INQUIRY command introduction
	cdb := makeInquiryCmd(0x00, false, 0x60)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return err
	}

	inq, err := parseInquiryData(buf)
	d.ModelNumber = inq.model
	d.FirmwareRevision = inq.revision

	// 5.4.18 Supported Vital Product Data pages (00h)
	cdb = makeInquiryCmd(0x00, true, 0xFC)
	buf, err = d.post(fd, cdb)
	if err != nil {
		return err
	}

	supportedVPDs, err := parseInquiryVPDSupportedPages(buf)
	if err != nil {
		return err
	}

	if _, supported := supportedVPDs[0x80]; supported {
		// 5.4.19 Unit Serial Number page (80h)
		cdb = makeInquiryCmd(0x80, true, 0x18)
		buf, err = d.post(fd, cdb)
		if err != nil {
			return err
		}

		serial, err := parseInquiryVPDSerialNumber(buf)
		if err != nil {
			return err
		}
		d.SerialNumber = serial
	}

	return nil
}

func parseInquiryVPDSupportedPages(buf []byte) (map[byte]struct{}, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 5.4.18 Supported Vital Product Data pages (00h)

	if len(buf) < 4 {
		return nil, fmt.Errorf("data length is too short")
	}

	// skit fields not to use
	off := 3

	// length
	length := int(buf[off])
	off++

	if length > len(buf)+3 || length <= 0 {
		return nil, fmt.Errorf("datalen is invalid (bufsize %d, len %d)", len(buf), length)
	}

	supportedVPDs := make(map[byte]struct{})
	for i := 0; i < length; i++ {
		supportedVPDs[buf[i+off]] = struct{}{}
	}

	return supportedVPDs, nil
}

func parseInquiryData(buf []byte) (*inquiryData, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 3.6.2 Standard INQUIRY data

	if len(buf) < 5 {
		return nil, fmt.Errorf("data length is too short")
	}

	// skip feilds not to use
	off := 4

	// length
	length := int(buf[off])
	off++

	// fields not to use (3 octets)
	// + vendor identification (8 octets) + product identification (16 octets)
	// + product revsion level (4 octets)  = 31
	if length < 31 || len(buf) < 31+4 || length+4 > len(buf) {
		return nil, fmt.Errorf("data length is invalid")
	}

	inq := new(inquiryData)

	// skip feilds not to use
	off += 3

	// vendor indentification (8 octets)
	inq.vendor = strings.TrimSpace(string(buf[off : off+8]))
	off += 8

	// product identification (16 octets)
	inq.model = strings.TrimSpace(string(buf[off : off+16]))
	off += 16

	// product revision level (4 octets)
	inq.revision = strings.TrimSpace(string(buf[off : off+4]))

	return inq, nil
}

func parseInquiryVPDSerialNumber(buf []byte) (string, error) {
	// SCSI Commands Reference Manual
	// https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf
	// 5.4.19 Unit Serial Number page (80h)

	if len(buf) < 4 {
		return "", fmt.Errorf("data length is too short")
	}

	// skit fields not to use
	off := 3

	// length
	length := int(buf[off])
	off++

	if length <= 0 || len(buf) < length+4 {
		return "", fmt.Errorf("data length is invalid (page len: %d, buf len: %d)", length, len(buf))
	}

	serial := strings.Trim(string(buf[off:off+length]), "\u0000\n\t\r ")
	return serial, nil
}
