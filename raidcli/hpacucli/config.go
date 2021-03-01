package hpacucli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/moxspec/moxspec/raidcli"
	"github.com/moxspec/moxspec/util"
)

var lvmap = map[string]raidcli.Level{
	"0":   raidcli.RAID0,
	"1":   raidcli.RAID1,
	"5":   raidcli.RAID5,
	"6":   raidcli.RAID6,
	"1+0": raidcli.RAID10,
	"0+1": raidcli.RAID01,
}

func splitConfigDetailSections(in string) (ctLines, ldpdLines []string, err error) {
	var buf *[]string

	for _, line := range strings.Split(in, "\n") {
		log.Debug(line)

		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		switch {
		case buf == nil && util.HasPrefixIn(l, "Smart Array ", "Smart HBA", "Dynamic Smart Array ") && strings.Contains(l, "in Slot"):
			// hpacucli repeats a mark such as `Smart ... in Slot 0 ...` in case of Smart HBA series
			// so we need to check whether buf is nil
			// eg.
			//   Smart HBA H240ar in Slot 0 (Embedded)
			//         physicaldrive 1I:1:11 (port 1I:box 1:bay 11, SATA, 4000.7 GB, OK)
			//   Smart HBA H240ar in Slot 0 (Embedded)
			//         physicaldrive 1I:1:12 (port 1I:box 1:bay 12, SATA, 4000.7 GB, OK)
			log.Debug("controller lines started")
			buf = &ctLines
			*buf = append(*buf, l)
			continue
		case strings.HasPrefix(l, "Array: "):
			log.Debug("drive lines started")
			buf = &ldpdLines
			*buf = append(*buf, l)
			continue
		}

		// any other no need sections started
		if strings.HasPrefix(l, "SEP (Vendor ID") {
			break
		}

		if buf != nil {
			*buf = append(*buf, l)
		}
	}

	log.Debugf("got %d controller lines", len(ctLines))
	log.Debugf("got %d logical/physical drive lines", len(ldpdLines))

	return
}

func parseCTLines(lines []string) (sn, firm string, battery bool, pciaddr string, err error) {
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		key, val, err := raidcli.SplitKeyVal(l, ": ")
		if err != nil {
			continue
		}

		switch key {
		case "Serial Number":
			if sn == "" {
				sn = val
			}
		case "Firmware Version":
			if firm == "" {
				firm = val
			}
		case "Battery/Capacitor Count":
			cnt, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			if cnt > 0 {
				battery = true
			}
		case "PCI Address (Domain:Bus:Device.Function)":
			pciaddr = val
		}
	}

	return
}

func splitArrays(lines []string) (arrays [][]string, unassigned []string, err error) {
	var buf *[]string

	for _, line := range lines {
		log.Debug(line)

		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		if strings.HasPrefix(l, "Array: ") || l == "unassigned" {
			if buf != nil && len(*buf) > 0 {
				arrays = append(arrays, *buf)
			}
			buf = &[]string{}
		}

		if buf != nil {
			*buf = append(*buf, l)
		}

	}

	if buf != nil && len(*buf) > 0 {
		if (*buf)[0] == "unassigned" {
			unassigned = *buf
		} else {
			arrays = append(arrays, *buf)
		}
	}

	log.Debugf("got %d arrays", len(arrays))
	log.Debugf("got %d unassigned lines", len(unassigned))

	return
}

func splitLDChunks(lines []string) (lds [][]string, err error) {
	var buf *[]string

	for _, line := range lines {
		log.Debug(line)

		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		if strings.HasPrefix(l, "Logical Drive: ") {
			if buf != nil && len(*buf) > 0 {
				lds = append(lds, *buf)
			}
			buf = &[]string{}
		}

		if buf != nil {
			*buf = append(*buf, l)
		}
	}

	if buf != nil && len(*buf) > 0 {
		lds = append(lds, *buf)
	}

	log.Debugf("got %d lds", len(lds))

	return
}

func splitLDPDSections(lines []string) (ldLines, pdLines []string, err error) {
	var buf *[]string

	for _, line := range lines {
		log.Debug(line)

		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		switch {
		case strings.HasPrefix(l, "Logical Drive: "):
			log.Debug("ld lines started")
			buf = &ldLines
			*buf = append(*buf, l)
			continue
		case strings.HasPrefix(l, "physicaldrive") && len(strings.Fields(l)) == 2:
			log.Debug("pd lines started")
			buf = &pdLines
			*buf = append(*buf, l)
			continue
		}

		if buf != nil {
			*buf = append(*buf, l)
		}
	}

	log.Debugf("got %d logical drive lines", len(ldLines))
	log.Debugf("got %d physical drive lines", len(pdLines))

	return
}

func parseLDLines(lines []string) (*LogDrive, error) {
	if len(lines) == 0 {
		log.Debug("no ld info found")
		return nil, nil
	}

	var ld *LogDrive

	for _, line := range lines {
		l := strings.TrimSpace(line)

		if strings.HasPrefix(l, "Logical Drive: ") {
			if ld != nil {
				return nil, fmt.Errorf("multiple log drive found")
			}

			_, val, err := raidcli.SplitKeyVal(l, ":")
			if err != nil {
				continue
			}

			ld = new(LogDrive)
			ld.VolumeID = val
			ld.Label = fmt.Sprintf("vol:%s", val)
			continue
		}

		if ld == nil {
			continue
		}

		if strings.HasPrefix(l, "physicaldrive") {
			flds := strings.Fields(l)
			if len(flds) < 2 {
				continue
			}
			if strings.Count(flds[1], ":") != 2 {
				continue
			}
			ld.pdIDList = append(ld.pdIDList, flds[1])
		}

		key, val, err := raidcli.SplitKeyVal(l, ":")
		if err != nil {
			continue
		}

		switch key {
		case "Size":
			ld.Size = raidcli.ParseSize(val, raidcli.Decimal)
		case "Fault Tolerance":
			ld.RAIDLv = formatRAIDLv(val)
		case "Strip Size":
			ld.StripSize = raidcli.ParseSize(val, raidcli.Binary)
		case "Status":
			ld.State = val
		case "Unique Identifier":
			ld.UUID = val
		case "Disk Name":
			ld.DiskName = val
		}
	}

	return ld, nil
}

func formatRAIDLv(val string) raidcli.Level {
	if lv, ok := lvmap[val]; ok {
		return lv
	}
	return raidcli.Unknown
}

func parsePDLines(lines []string) ([]*PhyDrive, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("no pd info found, this controller has no drive")
	}

	var pd *PhyDrive
	var pds []*PhyDrive

	for _, line := range lines {
		l := strings.TrimSpace(line)

		if util.HasPrefixIn(l, "physicaldrive", "Enclosure ", "Expander ") {
			if pd != nil {
				pds = append(pds, pd)
			}
			pd = new(PhyDrive)
			continue
		}

		if pd == nil {
			continue
		}

		key, val, err := raidcli.SplitKeyVal(l, ": ")
		if err != nil {
			continue
		}

		switch key {
		case "Port":
			pd.Port = val
		case "Box":
			pd.Box = val
		case "Bay":
			pd.Bay = val
		case "Status":
			pd.Status = val
		case "Interface Type":
			pd.Protocol = val
		case "Size":
			pd.Size = raidcli.ParseSize(val, raidcli.Decimal)
		case "Drive exposed to OS": // pass-through flag
		case "Rotational Speed":
			r, err := strconv.Atoi(val)
			if err == nil {
				pd.Rotation = uint(r)
			}
		case "Firmware Revision":
			pd.Firmware = val
		case "Serial Number":
			pd.SerialNumber = val
		case "Model":
			pd.Model = raidcli.ShapeSpacedString(val)
		case "Current Temperature (C)":
			t, err := strconv.Atoi(val)
			if err == nil {
				pd.CurTemp = t
			}
		case "Maximum Temperature (C)":
			t, err := strconv.Atoi(val)
			if err == nil {
				pd.MaxTemp = t
			}
		case "PHY Transfer Rate":
			flds := strings.Split(val, ", ")
			if len(flds) == 2 {
				pd.NegSpeed = flds[0]
			}
		}
	}

	if pd != nil {
		pds = append(pds, pd)
	}

	return pds, nil
}
