package spc

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	smartAttrSize   = 12
	maxSmartAttrNum = 30
)

// SmartRecord represents a s.m.a.r.t record
type SmartRecord struct {
	ID        byte
	Current   byte
	Worst     byte
	Raw       int64
	Threshold byte
	Name      string
}

func (d *Device) appendErrorRecord(a *smartAttr, t *smartThreshold, name string) {
	e := newSmartRecord(a.id, a.current, a.worst, a.raw, t.threshold, name)
	if e != nil {
		d.ErrorRecords = append(d.ErrorRecords, e)
	}
}

func (d *Device) decodeSmart(fd *os.File) error {
	// ACS-4
	// http://www.t13.org/documents/UploadedDocuments/docs2016/di529r14-ATAATAPI_Command_Set_-_4.pdf

	// 7.44 SMART
	// 7.44.2.3 Inputs
	cdb := makeATAPThruCmd(0xD0, 0x00, 0x4F, 0xC2, 0xB0)
	buf, err := d.post(fd, cdb)
	if err != nil {
		return err
	}

	attrs, err := parseSmartAttrs(buf)
	if err != nil {
		return err
	}

	// 7.44 SMART
	// 7.44.2.3 Inputs
	cdb = makeATAPThruCmd(0xD1, 0x00, 0x4F, 0xC2, 0xB0)
	buf, err = d.post(fd, cdb)
	if err != nil {
		return err
	}

	tDict, err := parseSmartThresholdsDict(buf)
	if err != nil {
		return err
	}

	isUnhealthy := func(a *smartAttr) bool {
		// can not judge condition
		if _, ok := tDict[a.id]; !ok {
			return false
		}

		t := tDict[a.id]

		log.Debugf("%03d(0x%02X): current=%3d, worst=%3d, threshold=%3d, raw=%d", a.id, a.id, a.current, a.worst, t.threshold, a.raw)

		if t.threshold < a.worst {
			return false
		}

		return true
	}

	log.Debug("scanning smart values")

	for i := 0; i < len(attrs); i++ {
		attr := attrs[i]
		switch attr.id {
		case 0x01: // 1. Raw Read Error Rate
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Raw Read Error Rate")
			}
		case 0x05: // 5. Reallocated Sectors Count
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Reallocated Sectors Count")
			}
		case 0x09: // 9. Power-On Hours
			if uint64(attr.raw) > d.PowerOnHours {
				d.PowerOnHours = uint64(attr.raw)
			}
		case 0x0A: // 10. Spin Retry Count
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Spin Retry Count")
			}
		case 0x0C: // 12. Power Cycle Count
			if uint64(attr.raw) > d.PowerCycleCount {
				d.PowerCycleCount = uint64(attr.raw)
			}
		case 0x0D: // 13. Soft Read Error Rate
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Soft Read Error Rate")
			}
		case 0xB8: // 184.End-to-End error / IOEDC
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "End-to-End error / IOEDC")
			}
		case 0xBB: // 187. Reported Uncorrectable Errors
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Reported Uncorrectable Errors")
			}
		case 0xBC: // 188. Command Timeout
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Command Timeout")
			}
		case 0xBE: // 190. Temperature Difference or Airflow Temperature (INTEL)
			cur := int16(attr.rawBytes[0])
			if cur > 0 {
				d.CurTemp = cur
			}
			min := int16(attr.rawBytes[2])
			if min > 0 {
				d.MinTemp = min
			}
			max := int16(attr.rawBytes[3])
			if max > 0 {
				d.MaxTemp = max
			}
		case 0xC0: // 192. Power-off Retract Count or Unsafe Shutdown Count
			if uint64(attr.raw) > d.UnsafeShutdownCount {
				d.UnsafeShutdownCount = uint64(attr.raw)
			}
		case 0xC2: // 194. Temperature or Temperature Celsius
			cur, min, max := readSmartTemp(attr.rawBytes)
			if cur > 0 {
				d.CurTemp = cur
			}
			if min > 0 {
				d.MinTemp = min
			}
			if max > 0 {
				d.MaxTemp = max
			}
		case 0xC3: // 195. Hardware ECC recovered
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Hardware ECC recovered")
			}
		case 0xC4: // 196. Reallocation Event Count
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Reallocation Event Count")
			}
		case 0xC5: // 197. Current Pending Sector Count
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Current Pending Sector Count")
			}
		case 0xC6: // 198. Off-Line Scan Uncorrectable Sector Count
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Off-Line Scan Uncorrectable Sector Count")
			}
		case 0xDC: // 220. Disk Shift
			if isUnhealthy(attr) {
				d.appendErrorRecord(attr, tDict[attr.id], "Disk Shift")
			}
		case 0xF1: // 241. Total LBAs Written
			if uint64(attr.raw) > d.TotalLBAWritten {
				d.TotalLBAWritten = uint64(attr.raw)
			}
		case 0xF2: // 242. Total LBAs Read
			if uint64(attr.raw) > d.TotalLBARead {
				d.TotalLBARead = uint64(attr.raw)
			}
		}
	}

	return nil
}

func newSmartRecord(id byte, c, w byte, r int64, t byte, name string) *SmartRecord {
	s := new(SmartRecord)
	s.ID = id
	s.Current = c
	s.Worst = w
	s.Raw = r
	s.Threshold = t
	s.Name = name
	return s
}

type smartAttr struct {
	id       byte
	flags    uint16
	current  byte
	worst    byte
	rawBytes []byte
	raw      int64
	reserv   byte
}

func (s smartAttr) summary() string {
	return fmt.Sprintf("%03d(0x%02X): current=%3d, worst=%3d, raw=%d", s.id, s.id, s.current, s.worst, s.raw)
}

func newSmartAttr(buf []byte) (*smartAttr, error) {
	if len(buf) != smartAttrSize {
		return nil, fmt.Errorf("smartAttr: invalid bytes given (%d)", len(buf))
	}

	a := new(smartAttr)
	a.id = buf[0]
	a.flags = binary.LittleEndian.Uint16(buf[1:3])
	a.current = buf[3]
	a.worst = buf[4]
	a.rawBytes = buf[5:11]
	a.raw = readSmartRawVal(a.rawBytes)
	a.reserv = buf[11]

	return a, nil
}

func readSmartRawVal(buf []byte) int64 {
	if len(buf) != 6 {
		return 0
	}

	var raw int64
	for i := 0; i < 6; i++ {
		raw = raw | (int64(buf[i]) << (8 * uint(i)))
	}

	return raw
}

type smartThreshold struct {
	id        byte
	threshold byte
	reserve   []byte // 10 bytes
}

func (s smartThreshold) summary() string {
	return fmt.Sprintf("%03d(0x%02X): threshold=%3d", s.id, s.id, s.threshold)
}

func newSmartThreshold(buf []byte) (*smartThreshold, error) {
	if len(buf) != smartAttrSize {
		return nil, fmt.Errorf("smartThreshold: invalid bytes given (%d)", len(buf))
	}

	s := new(smartThreshold)
	s.id = buf[0]
	s.threshold = buf[1]
	s.reserve = buf[2:]

	return s, nil
}

func parseSmartAttrs(buf []byte) ([]*smartAttr, error) {
	var attrs []*smartAttr

	// first 2 bytes are revision number.
	for i := 2; i < maxSmartAttrNum*smartAttrSize; i += smartAttrSize {
		if i+smartAttrSize >= len(buf) {
			break
		}

		a, err := newSmartAttr(buf[i : i+smartAttrSize])
		if err != nil {
			return nil, err
		}

		if a.id == 0 {
			continue
		}

		log.Debugf("%s", a.summary())
		attrs = append(attrs, a)
	}

	return attrs, nil
}

func parseSmartThresholdsDict(buf []byte) (map[byte]*smartThreshold, error) {
	dict := make(map[byte]*smartThreshold)

	// first 2 bytes are revision number.
	for i := 2; i < maxSmartAttrNum*smartAttrSize; i = i + smartAttrSize {
		if i+smartAttrSize >= len(buf) {
			break
		}

		t, err := newSmartThreshold(buf[i : i+smartAttrSize])
		if err != nil {
			return nil, err
		}

		if t.id == 0 {
			continue
		}

		log.Debugf("%s", t.summary())
		dict[t.id] = t
	}

	return dict, nil
}

// guess possible value
func readSmartTemp(v []byte) (cur, min, max int16) {
	if len(v) != 6 {
		return
	}

	//  5  4  3  2  1  0 : Index
	// xx LL xx HH xx TT : Pattern 1
	// xx HH xx LL xx TT : Pattern 2
	// 00 00 HH LL xx TT : Pattern 3
	// 00 00 00 HH LL TT : Pattern 4
	//
	// TT: Current,  LL: Min, HH: Max
	cur = int16(v[0])

	// Pattern 4
	if v[3] == 0 && v[4] == 0 && v[5] == 0 {
		min = int16(v[1])
		max = int16(v[2])
		return
	}

	// Pattern 3
	if v[4] == 0 && v[5] == 0 {
		min = int16(v[2])
		max = int16(v[3])
		return
	}

	// Pattern 2
	min = int16(v[2])
	max = int16(v[4])

	// Pattern 1
	if min > max {
		tmp := min
		min = max
		max = tmp
	}

	return
}
