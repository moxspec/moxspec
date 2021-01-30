package sas3ircu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actapio/moxspec/raidcli"
)

var lvmap = map[string]raidcli.Level{
	"RAID0":  raidcli.RAID0,
	"RAID1":  raidcli.RAID1,
	"RAID1E": raidcli.RAID1,
	"RAID5":  raidcli.RAID5,
	"RAID6":  raidcli.RAID6,
	"RAID10": raidcli.RAID10,
}

func genPDMap(pds []*PhyDrive) (map[string]*PhyDrive, error) {
	log.Debugf("generating pdmap[enc:slot]")
	pdmap := make(map[string]*PhyDrive)

	for _, pd := range pds {
		if pd == nil {
			continue
		}

		key := fmt.Sprintf("%s:%s", pd.EnclosureID, pd.SlotNumber)
		if _, ok := pdmap[key]; ok {
			return nil, fmt.Errorf("duplicate pd id found")
		}

		pdmap[key] = pd
	}

	return pdmap, nil
}

func splitSections(in string) (ctLines, ldLines, pdLines []string, err error) {
	var buf *[]string

	for _, line := range strings.Split(in, "\n") {
		log.Debug(line)

		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		// separater line
		if len(l) == strings.Count(l, "-") {
			continue
		}

		switch l {
		case "Controller information":
			log.Debug("controller lines started")
			buf = &ctLines
			continue
		case "IR Volume information":
			log.Debug("logical drive lines started")
			buf = &ldLines
			continue
		case "Physical device information":
			log.Debug("physical drive lines started")
			buf = &pdLines
			continue
		}

		// any other no need sections started
		if strings.HasSuffix(l, "information") && !strings.Contains(l, ":") {
			break
		}

		if buf != nil && !strings.HasPrefix(l, "-----") {
			*buf = append(*buf, l)
		}
	}

	log.Debugf("got %d controller lines", len(ctLines))
	log.Debugf("got %d logical drive lines", len(ldLines))
	log.Debugf("got %d physical drive lines", len(pdLines))

	return
}

func parseCTLines(lines []string) (firm, bios string, err error) {
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		key, val, err := raidcli.SplitKeyVal(l, ":")
		if err != nil {
			continue
		}

		switch key {
		case "BIOS version":
			bios = val
		case "Firmware version":
			firm = val
		}
	}

	return
}

func parseLDLines(lines []string) ([]*LogDrive, error) {
	if len(lines) == 0 {
		log.Debug("no ld info found, all disks are exposed as pass-through drive")
		return nil, nil
	}

	var ld *LogDrive
	var lds []*LogDrive

	for _, line := range lines {
		l := strings.TrimSpace(line)

		if strings.HasPrefix(l, "IR volume") {
			if ld != nil {
				lds = append(lds, ld)
			}
			ld = new(LogDrive)
			continue
		}

		if ld == nil {
			continue
		}

		key, val, err := raidcli.SplitKeyVal(l, ":")
		if err != nil {
			continue
		}

		if strings.HasPrefix(key, "Size (in") {
			ld.Size = parseSizeIn(key, val)
			continue
		}

		switch key {
		case "Volume ID":
			id, err := strconv.Atoi(val)
			if err == nil {
				ld.VolumeID = id
				ld.Label = fmt.Sprintf("vol:%d", ld.VolumeID)
			}
		case "Volume wwid":
			ld.WWID = val
		case "Status of volume":
			ld.State = val
		case "RAID level":
			ld.RAIDLv = formatRAIDLv(val)
		default:
			if strings.HasPrefix(key, "PHY[") && strings.HasSuffix(key, "] Enclosure#/Slot#") {
				ld.pdIDList = append(ld.pdIDList, val)
			}
		}
	}

	if ld != nil {
		lds = append(lds, ld)
	}

	return lds, nil
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

		if strings.HasPrefix(l, "Device is a") {
			if isValidPD(pd) {
				pds = append(pds, pd)
			}
			pd = new(PhyDrive)
			continue
		}

		if pd == nil {
			continue
		}

		key, val, err := raidcli.SplitKeyVal(l, ":")
		if err != nil {
			continue
		}

		switch key {
		case "Enclosure #":
			pd.EnclosureID = val
		case "Slot #":
			pd.SlotNumber = val
		case "Model Number":
			pd.Model = val
		case "Firmware Revision":
			pd.Firmware = val
		case "Serial No":
			pd.SerialNumber = val
		case "SAS Address":
			pd.SASAddress = val
		case "State":
			pd.State = val
		case "Size (in MB)/(in sectors)":
			pd.Size = parseMBSize(val)
		case "Protocol":
			pd.Protocol = val
		case "Drive Type":
			pd.DriveType = val
			pd.SolidStateDrive = strings.Contains(strings.ToUpper(val), "SSD")
		case "Device Type":
			// NOTE:
			// An enclosure services device has "Device Type" key.
			// It is very similar to "Drive Type" key but it is NOT typo.
			// See also unit test.
			// DAMN IT.
			pd.DriveType = val
		}
	}

	if isValidPD(pd) {
		pds = append(pds, pd)
	}

	return pds, nil
}

func isValidPD(pd *PhyDrive) bool {
	if pd == nil {
		return false
	}

	if pd.Size == 0 {
		return false
	}

	if pd.State == "Missing (MIS)" {
		return false
	}

	invalids := []string{
		"undetermined",
		"enclosure",
	}

	dtype := strings.ToLower(pd.DriveType)
	for _, i := range invalids {
		if strings.Contains(dtype, i) {
			return false
		}
	}

	return true
}

func parseSizeIn(key, val string) uint64 {
	unit := strings.TrimPrefix(key, "Size (in ")
	unit = strings.TrimSuffix(unit, ")")
	unit = strings.TrimSpace(unit)

	sz, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}

	return uint64(sz * raidcli.GetMultiplier(unit, raidcli.Binary))
}

func parseMBSize(in string) uint64 {
	flds := strings.Split(in, "/")
	if len(flds) != 2 {
		return 0
	}

	sz, err := strconv.Atoi(flds[0])
	if err != nil {
		return 0
	}

	return uint64(sz) * 1024 * 1024
}
