package scsi

import (
	"fmt"
	"strings"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

// 7.8.5 Device Constituents VPD page
// Association Field values
const (
	aLogicalUnit = iota // 00b
	aTargetPort         // 01b
	aSCSITarget         // 10b
)

// 7.6.1 Protocol specific parameters introduction
// Protocol Identifier values
const (
	pFCP4   = iota // Fibre Channel
	pSPI5          //  Parallel SCSI
	pSSAS3P        // SSA
	pSBP3          //  IEEE 1394
	pSRP           // SCSI Remote Direct Memory Access Protocol
	pISCSI         // Internet SCSI (iSCSI)
	pSPL           // SAS Protocol Layer
	pADT2          // Automation/Drive Interface Transport Protocol

)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("scsi")
}

// DecodeSerialNumber returns serial number from scsi vpd page80
func DecodeSerialNumber(path string) string {
	bs, err := util.LoadBytes(path)
	if err != nil {
		return ""
	}

	sn, _ := decodePg80(bs)
	return sn
}

func decodePg80(in []byte) (string, error) {
	log.Debugf("decodinfg page80 (len:%d)", len(in))

	if len(in) < 4 {
		return "", fmt.Errorf("pg data is too short")
	}

	if in[1] != 0x80 {
		return "", fmt.Errorf("page code is not 80h")
	}

	if in[3] == 0 { // Page Length
		return "", nil
	}

	return strings.TrimSpace(string(in[4:])), nil
}

// DecodeWWN returns wwn from scsi vpd page83
func DecodeWWN(path string) string {
	log.Debugf("loading vpd from %s", path)
	bs, err := util.LoadBytes(path)
	if err != nil {
		return ""
	}

	descs, err := decodePg83(bs)
	if err != nil {
		log.Debug(err)
		return ""
	}

	for _, d := range descs {
		wwn := parseWWN(d)
		if wwn != "" {
			log.Debugf("found wwn: %s", wwn)
			return wwn
		}
	}

	return ""
}

// VPD format
//     |-------|-------|-------|-------|-------|-------|-------|-------|
//     |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 0 | Peripheral Qualifier  | Periheral Device Type                 |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 1 | Page Code (0x83)                                              |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 2 | (MSB) Page Length                                             |
// | 3 |                                                         (LSB) |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 4 | Designation Descriptor (first)                                |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// |   | ...                                                           |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | n | Designation Descriptor (n)                                    |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
func decodePg83(in []byte) ([]*pg83Descriptor, error) {
	log.Debugf("decodinfg page83 (len:%d)", len(in))

	if len(in) < 4 {
		return nil, fmt.Errorf("pg data is too short")
	}

	if in[1] != 0x83 {
		return nil, fmt.Errorf("page code is not 83h")
	}

	plen := int(in[2])<<8 | int(in[3])
	if plen == 0 {
		return nil, nil
	}

	log.Debugf("page length: %d", plen)

	var descs []*pg83Descriptor

	log.Debug("scanning identification descriptor list")
	i := 4
	for {
		if i > plen {
			break
		}

		if i+3 > plen {
			log.Debugf("invalid descriptor (too short header)")
			log.Debugf("it caused by invalid data or bug")
			break
		}

		ilen := int(in[i+3])

		if ilen == 0 {
			i = i + 4
			continue
		}

		if i+3+ilen > len(in) {
			log.Debugf("offset: %d + identidier len: %d is longer than page length", i, ilen, plen)
			log.Debugf("it caused by invalid data or bug")
			break
		}

		next := i + 4 + ilen
		d := shapeDescriptor(in[i:next])
		if d == nil {
			continue
		}

		descs = append(descs, d)

		i = next
	}

	return descs, nil
}

func parseWWN(d *pg83Descriptor) string {
	log.Debug("parsing wwn")

	log.Debugf("desig type: 0x%x", d.DesigType)
	switch d.DesigType {
	case 0x00:
		log.Debug("Vendor specific, nothing to do")
	case 0x01:
		log.Debug("T10 vendor identification, nothing to do")
	case 0x02:
		log.Debug("EUI-64")
	case 0x03:
		log.Debug("NAA")
		return parseNAAField(d)
	case 0x04:
		log.Debug("Relative target port, nothing to do")
	case 0x05:
		log.Debug("Target port group, nothing to do")
	case 0x06:
		log.Debug("Logical unit group, nothing to do")
	case 0x07:
		log.Debug("MD5 logical unit identifier, nothing to do")
	case 0x08:
		log.Debug("SCSI name string")
		return decodeSCSINameString(string(d.Body))
	case 0x09:
		log.Debug("Protocol specific port identifier, nothing to do")
	case 0x0a:
		log.Debug("SPC5 UUID identifier RFC 4122")
		return decodeSPC5ID(d.Body)
	default:
		log.Debug("reserved, nothing to do")
	}

	return ""
}

func decodeSPC5ID(body []byte) string {
	if len(body) != 18 {
		log.Debugf("unexpected length: %d", len(body))
		return ""
	}

	//  FORMAT: 4       - 2  - 2  - 2  - 6
	//  FORMAT: 00000000-0000-0000-0000-000000000000
	return fmt.Sprintf("%x-%x-%x-%x-%x", body[2:6], body[6:8], body[8:10], body[10:12], body[12:18])
}

func decodeSCSINameString(name string) string {
	log.Debugf("scsi name = %s", name)
	for _, prefix := range []string{"NAA.", "naa.", "EUI.", "eui.", "iqn."} {
		if strings.HasPrefix(name, prefix) {
			s := strings.Replace(name, "\x00", "", -1)
			return strings.TrimSpace(strings.TrimPrefix(s, prefix))
		}
	}
	log.Debug("no valid prefix found")
	return ""
}

// NAA Designator Basic Format
//     |-------|-------|-------|-------|-------|-------|-------|-------|
//     |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 0 | NAA                           | NAA Specific Data             |
// |---|-------|-------|-------|-------|                               |
// | 1 |                                                               |
// | n |                                                               |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
func parseNAAField(d *pg83Descriptor) string {
	if d.CodeSet != 0x01 {
		log.Debug("expected binary value")
		return ""
	}

	if len(d.Body) < 1 {
		log.Debug("empty body")
		return ""
	}

	naa := d.Body[0] >> 4 & 0xF

	log.Debugf("naa: 0x%x", naa)
	switch naa {
	case 0x02:
		log.Debug("NAA IEEE extended")
		return formatNAABytes(d.Body, 8)
	case 0x03:
		log.Debug("Locally assigned")
		if d.PID != pSPL || d.Assoc != aTargetPort {
			return formatNAABytes(d.Body, 8)
		}
	case 0x05:
		log.Debug("IEEE Registered")
		if d.PID != pSPL || d.Assoc != aTargetPort {
			return formatNAABytes(d.Body, 8)
		}
	case 0x06:
		log.Debug("NAA IEEE Registered extended")
		return formatNAABytes(d.Body, 16)
	default:
		log.Debug("unknown")
	}

	log.Debug("nothing to do")
	return ""
}

func formatNAABytes(body []byte, l int) string {
	if l < 0 {
		return ""
	}

	if len(body) != l {
		log.Debugf("unexpected length: %d", l)
		return ""
	}

	return fmt.Sprintf("%x", body)
}

// PG83 Descriptor Format
//     |-------|-------|-------|-------|-------|-------|-------|-------|
//     |   7   |   6   |   5   |   4   |   3   |   2   |   1   |   0   |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 0 | Protocol Identifier           | Code Set                      |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 1 | PIV   | Rsrvd | Association   | Designator Type               |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 2 | Reserved                                                      |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 3 | Designator Length (n - 3)                                     |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
// | 4 | Designator Body                                               |
// | n |                                                               |
// |---|-------|-------|-------|-------|-------|-------|-------|-------|
func shapeDescriptor(in []byte) *pg83Descriptor {
	// If the ASSOCIATION field contains a value other than 1h or 2h or the PIV bit is set to zero,
	// then the PROTOCOL IDENTIFIER field should be ignored.
	d := &pg83Descriptor{
		CodeSet:   in[0] & 0xF,        // code set
		PID:       (in[0] >> 4) & 0xF, // protocol identidier
		DesigType: in[1] & 0xF,        // designator type
		Assoc:     (in[1] >> 4) & 0x3, // association
		PIV:       (in[1] >> 7) & 0x1, // piv
		Len:       int(in[3]),
		Body:      in[4:],
	}

	fmtr := "codeset:%d, pid:%d (piv:%d), assoc:%d, type:%d, len:%d, body:%s"
	if d.CodeSet < 2 {
		fmtr = "codeset:%d, pid:%d (piv:%d), assoc:%d, type:%d, len:%d, body:%+v"
	}
	log.Debugf(fmtr, d.CodeSet, d.PID, d.PIV, d.Assoc, d.DesigType, d.Len, d.Body)

	return d
}
