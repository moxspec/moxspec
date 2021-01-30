package megacli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actapio/moxspec/raidcli"
)

func getAllPD(num int) ([]*PhyDrive, error) {
	res, err := raidcli.Run(clipath, "-PDList", fmt.Sprintf("-a%d", num), "-NoLog")
	if err != nil {
		return nil, err
	}

	return parsePDList(strings.Split(res, "\n"))
}

func parsePDList(lines []string) ([]*PhyDrive, error) {
	var drives []*PhyDrive
	var pd *PhyDrive

	for _, line := range lines {
		key, val, err := raidcli.SplitKeyVal(line, ":")
		if err != nil {
			log.Debug(err)
			continue
		}

		if key == "Enclosure Device ID" {
			if pd != nil {
				drives = append(drives, pd)
			}
			pd = new(PhyDrive)
			pd.EnclosureID = val
		}

		if pd == nil {
			continue
		}

		// NOTE: Ver 8.04.07 has typo
		if key == "Drive's position" || key == "Drive's postion" {
			grp, span, arm, err := parseDrivePos(val)
			if err != nil {
				return nil, err
			}
			pd.Group = uint16(grp)
			pd.Span = uint16(span)
			pd.Arm = uint16(arm)
			log.Debugf("grp: %d, span: %d, arm: %d", pd.Group, pd.Span, pd.Arm)
			continue
		}

		switch key {
		case "WWN":
			pd.WWN = val
		case "Slot Number":
			pd.SlotNumber = val
		case "PD Type":
			pd.Type = val
		case "Raw Size":
			pd.Size = raidcli.ParseSize(val, raidcli.Binary)
		case "Logical Sector Size":
			pd.LogBlockSize = parseSectorSize(val)
		case "Physical Sector Size":
			pd.PhyBlockSize = parseSectorSize(val)
		case "Firmware state":
			pd.State = val
		case "Device Firmware Level":
			pd.FirmwareRevision = val
		case "Connected Port Number":
			pd.ConnectedPort = val
		case "Inquiry Data":
			pd.Model = raidcli.ShapeSpacedString(val)
			pd.InquiryRaw = val
		case "Device Speed":
			pd.DriveSpeed = val
		case "Link Speed":
			pd.LinkSpeed = val
		case "Media Type":
			pd.SolidStateDrive = (strings.Contains(strings.ToLower(val), "solid state"))
		case "Drive Temperature":
			pd.CurTemp = parseTemp(val)
		case "Drive has flagged a S.M.A.R.T alert":
			pd.SMARTAlert = (strings.ToLower(val) == "yes")
		case "Media Error Count":
			pd.MediaErrorCount = parseUint16(val)
		case "Device Id":
			pd.DeviceID = parseUint16(val)
		}
	}

	// finalize
	if pd != nil {
		drives = append(drives, pd)
	}

	log.Debugf("found %d physical drives", len(drives))
	for _, d := range drives {
		log.Debugf("pd[%s:%s] grp:%d, span:%d, arm:%d", d.EnclosureID, d.SlotNumber, d.Group, d.Span, d.Arm)
	}

	return drives, nil
}

func parseSectorSize(val string) uint16 {
	size, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return uint16(size)
}

func parseTemp(val string) int16 {
	flds := strings.SplitN(strings.ToLower(val), "c", 2)
	if len(flds) != 2 {
		return 0
	}
	temp, err := strconv.ParseFloat(flds[0], 64)
	if err != nil {
		return 0
	}
	return int16(temp)
}

func parseDrivePos(val string) (grp, span, arm int, err error) {
	log.Debug("parsing drive position")

	keys := []string{
		"DiskGroup:",
		"Span:",
		"Arm:",
	}

	v := val
	for _, k := range keys {
		if !strings.Contains(val, k) {
			err = fmt.Errorf("not found %s in %s", k, val)
			return
		}
		v = strings.Replace(v, k, "", -1)
	}

	flds := strings.Split(v, ",")
	if len(flds) != 3 {
		err = fmt.Errorf("invalid drive pos")
		return
	}

	grp, err = strconv.Atoi(strings.TrimSpace(flds[0]))
	if err != nil {
		return
	}

	span, err = strconv.Atoi(strings.TrimSpace(flds[1]))
	if err != nil {
		return
	}

	arm, err = strconv.Atoi(strings.TrimSpace(flds[2]))
	return
}

func parseUint16(val string) uint16 {
	count, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return uint16(count)
}
