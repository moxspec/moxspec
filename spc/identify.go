package spc

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/moxspec/moxspec/util"
)

type identifyData struct {
	serialNumber        string
	firmwareRevision    string
	modelNumber         string
	sigSpeed            string
	negSpeed            string
	formFactor          string
	rotation            uint16
	transport           string
	selfTestSupport     bool
	errorLoggingSupport bool
}

func (d *Device) decodeIdentify(fd *os.File) error {
	// ACS-4
	// 7.12 IDENTIFY DEVICE ECh, PIO Data-In
	// 7.12.3 Inputs
	// http://www.t13.org/documents/UploadedDocuments/docs2016/di529r14-ATAATAPI_Command_Set_-_4.pdf
	cdb := makeATAPThruCmd(0x00, 0x00, 0x00, 0x00, 0xEC)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return err
	}

	l := parseIdentify(buf)

	d.SerialNumber = l.serialNumber
	d.FirmwareRevision = l.firmwareRevision
	d.ModelNumber = l.modelNumber
	d.SigSpeed = l.sigSpeed
	d.NegSpeed = l.negSpeed
	d.FormFactor = l.formFactor
	d.Rotation = l.rotation
	d.Transport = l.transport
	d.SelfTestSupport = l.selfTestSupport
	d.ErrorLoggingSupport = l.errorLoggingSupport

	return nil
}

func parseIdentify(buf []byte) *identifyData {
	l := new(identifyData)

	if len(buf) < 256 {
		log.Debugf("buf size is too short")
		return nil
	}

	// 9.11.7.2 SERIAL NUMBER field
	ser := readWordsAsATAString(buf, 10, 19)
	log.Debugf("serial number: word 10-19 is %s", ser)
	l.serialNumber = util.SanitizeString(ser)

	// 9.11.7.3 FIRMWARE REVISION field
	fim := readWordsAsATAString(buf, 23, 26)
	log.Debugf("firmware revision: word 23-26 is %s", fim)
	l.firmwareRevision = util.SanitizeString(fim)

	// 9.11.7.4 MODEL NUMBER field
	mod := readWordsAsATAString(buf, 27, 46)
	log.Debugf("model number: word 27-46 is %s", mod)
	l.modelNumber = util.SanitizeString(mod)

	// 9.11.10.2 SATA Capabilities
	l.sigSpeed = decodeSigSpeed(readWord(buf, 76))
	l.negSpeed = decodeNegSpeed(readWord(buf, 77))

	// 9.11.5.5 Form Factor (NOMINAL FORM FACTOR field)
	ffwd := readWordAsBytes(buf, 168)
	ff := decodeFormFactor(ffwd[0])
	log.Debugf("form factor: word 168 is %v (%s)", ffwd, ff)
	l.formFactor = ff

	// 9.11.5.4 NOMINAL MEDIA ROTATION RATE field
	var rotwd uint16
	rotwd = readWord(buf, 217)
	log.Debugf("nominal media rotation rate: word 217 is %d", rotwd)
	if rotwd == 1 || (0x400 < rotwd && rotwd < 0xFFFF) {
		l.rotation = rotwd
	}

	// Transport major version number
	var satawd uint16
	satawd = readWord(buf, 222)
	log.Debugf("transport major version number: word 222 is %d", satawd)
	l.transport = decodeTransport(satawd)

	// self-test and error logging supports
	var feat uint16
	feat = readWord(buf, 87)

	// 9.11.5.2.25 SELF-TEST SUPPORTED field
	l.selfTestSupport = decodeSelfTestSupport(feat)

	// 9.11.5.2.26 ERROR LOGGING SUPPORTED field
	l.errorLoggingSupport = decodeErrLoggingSupport(feat)

	return l
}

// the publication of the SATA Revision 3.4 Specification
// https://sata-io.org/sites/default/files/documents/SATA%20Spec%20Rev%203%204%20PR%20FINAL.pdf

// 7.13 IDENTIFY DEVICE ECh, PIO Data-In
// Table 50 IDENTIFY DEVICE data (part 5 of 19)
//
// bit  def
// --------
// 7:4  Reserved
//   3  Supports SATA Gen3 Signaling Speed (6.0Gb/s)
//   2  Supports SATA Gen2 Signaling Speed (3.0Gb/s)
//   1  Supports SATA Gen1 Signaling Speed (1.5Gb/s)
//   0  Shall be cleared to zero
var speeds = []string{
	"unknown", // bit0
	"1.5Gb/s",
	"3.0Gb/s",
	"6.0Gb/s", // bit3
	"6.0Gb/s",
	"6.0Gb/s",
	"6.0Gb/s",
	"6.0Gb/s", // bit7
}

func decodeSigSpeed(w uint16) string {
	log.Debugf("signaling speed: word 76 is %d", w)
	i, err := util.FindMSB(byte(w))
	if err != nil {
		return "unknown"
	}
	if i > 3 {
		log.Debugf("reserved bit is set in %d", i)
		log.Debug("parse as max speed")
	}
	return speeds[i]

}

func decodeNegSpeed(w uint16) string {
	log.Debugf("negotiated speed: word 77 is %d", w)
	b := int(w) & 0x07
	return speeds[b]
}

func decodeSelfTestSupport(w uint16) bool {
	log.Debugf("commands and feature sets supported: word 87 is %d", w)
	if b := int(w) & 0x02; b != 0 {
		log.Debugf("self-test supported: word 87, bit 1 is %d", b>>1)
		return true
	}
	return false
}

func decodeErrLoggingSupport(w uint16) bool {
	log.Debugf("commands and feature sets supported: word 87 is %d", w)
	if b := int(w) & 0x01; b != 0 {
		log.Debugf("error logging supported: word 87, bit 0 is %d", b)
		return true
	}
	return false
}

// Table 50 IDENTIFY DEVICE data (part 19 of 19)
var transports = []string{
	"ATA8-AST", // bit0
	"SATA 1.0a",
	"SATA II: Extensions",
	"SATA 2.5",
	"SATA 2.6",
	"SATA 3.0", // bit5
	"SATA 3.1",
	"SATA 3.2",
	"SATA 3.3",
	"SATA 3.4",
}

func decodeTransport(w uint16) string {
	log.Debugf("transport id is %d", w)
	if w == 0 || w == 0xFFFF {
		return ""
	}

	tlen := len(transports)
	if tlen == 0 {
		return ""
	}

	for i := tlen - 1; i >= 0; i-- {
		if w&(1<<uint16(i)) != 0 {
			return transports[i]
		}
	}
	return ""
}

func decodeFormFactor(b byte) string {
	bb := b & 0x0F
	var ff string
	switch bb {
	case 0x01:
		ff = "5.25\""
	case 0x02:
		ff = "3.5\""
	case 0x03:
		ff = "2.5\""
	case 0x04:
		ff = "1.8\""
	case 0x05:
		ff = "< 1.8\""
	case 0x06:
		ff = "mSATA"
	case 0x07:
		ff = "M.2"
	case 0x08:
		ff = "MicroSSD"
	case 0x09:
		ff = "CFast"
	}
	return ff
}

func readDWord(b []byte, at int) uint32 {
	l := readWord(b, at)
	h := readWord(b, at+1)
	return (uint32(h)<<16 | uint32(l))
}

func readWord(b []byte, at int) uint16 {
	wb := readWordAsBytes(b, at)
	if wb == nil {
		return 0
	}

	return (uint16(wb[1])<<8 | uint16(wb[0]))
}

func readWordAsBytes(b []byte, at int) []byte {
	if at < 0 {
		return nil
	}

	if len(b) < at*2+1 {
		return nil
	}

	return []byte{
		b[at*2],
		b[at*2+1],
	}
}

// 3.4.9 ATA string convention
func readWordsAsATAString(b []byte, from, to int) string {
	if from > to {
		return ""
	}

	if from < 0 { // guard for out of range
		return ""
	}

	if len(b)%2 != 0 { // guard for out of range
		return ""
	}

	if len(b) < to*2+1 {
		return ""
	}

	var res []byte
	for i := from * 2; i <= to*2; i = i + 2 {
		res = append(res, b[i+1], b[i])
	}
	res = bytes.Trim(res, "\x00")

	return strings.TrimSpace(fmt.Sprintf("%s", res))
}
