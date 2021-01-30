package megacli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actapio/moxspec/raidcli"
)

var lvmap = map[string]raidcli.Level{
	"Primary-0, Secondary-0, RAID Level Qualifier-0": raidcli.RAID0,
	"Primary-1, Secondary-0, RAID Level Qualifier-0": raidcli.RAID1,
	"Primary-5, Secondary-0, RAID Level Qualifier-3": raidcli.RAID5,
	"Primary-6, Secondary-0, RAID Level Qualifier-3": raidcli.RAID6,
	"Primary-1, Secondary-3, RAID Level Qualifier-0": raidcli.RAID10,
}

func getLDList(num int) ([]*LogDrive, error) {
	res, err := raidcli.Run(clipath, "-LDPDInfo", fmt.Sprintf("-a%d", num), "-NoLog")
	if err != nil {
		return nil, err
	}

	return parseLDList(res)
}

func parseLDList(in string) ([]*LogDrive, error) {
	var drives []*LogDrive
	var ld *LogDrive

	var pdLinesStarted bool
	var pdLines []string

	for _, line := range strings.Split(in, "\n") {
		key, val, err := raidcli.SplitKeyVal(line, ":")
		if err != nil {
			log.Debug(err)
			continue
		}

		if key == "Virtual Drive" {
			if ld != nil {
				pds, err := parsePDList(pdLines)
				if err != nil {
					return nil, err
				}
				ld.PhyDrives = pds

				drives = append(drives, ld)
			}

			pdLinesStarted = false

			ld = new(LogDrive)
			grp, tgt, err := parseLogID(val)
			if err != nil {
				log.Debug(err)
				continue
			}
			ld.GroupID = uint(grp)
			ld.TargetID = uint(tgt)
			ld.Label = fmt.Sprintf("grp:%d", grp)
		}

		if ld == nil {
			continue
		}

		if pdLinesStarted {
			pdLines = append(pdLines, line)
			continue
		}

		switch key {
		case "RAID Level":
			ld.RAIDLv = parseRAIDLv(val)
		case "Size":
			ld.Size = raidcli.ParseSize(val, raidcli.Binary)
		case "State":
			ld.State = val
		case "Strip Size":
			ld.StripSize = raidcli.ParseSize(val, raidcli.Binary)
		case "Span Depth":
			// NOTE:
			// we should parse it to determine what raid is it
			// for example:
			//   RAID Level : Primary-1, Secondary-0, RAID Level Qualifier-0 with Span Depth: 1 means RAID 1
			//   RAID Level : Primary-1, Secondary-0, RAID Level Qualifier-0 with Span Depth: 2 means RAID 0+1
			if ld.RAIDLv != raidcli.RAID1 {
				continue
			}

			dep, err := strconv.Atoi(val)
			if err != nil {
				log.Debug(err) // it's not fatal
				continue
			}

			if dep > 1 {
				ld.RAIDLv = raidcli.RAID01
			}
		case "Current Cache Policy":
			ld.CachePolicy = parseCachePolicy(val)
		case "Number of Spans":
			pdLinesStarted = true
			pdLines = []string{}
		}
	}

	// finalize
	if ld != nil {
		pds, err := parsePDList(pdLines)
		if err != nil {
			return nil, err
		}
		ld.PhyDrives = pds
		drives = append(drives, ld)
	}

	log.Debugf("found %d logical drives", len(drives))
	for _, d := range drives {
		log.Debugf("ld grp:%d, tgt:%d, lv:%s", d.GroupID, d.TargetID, d.RAIDLv)
	}

	return drives, nil
}

func parseLogID(val string) (grp int, tgt int, err error) {
	if !strings.Contains(val, "(Target Id:") {
		err = fmt.Errorf("invalid format")
		return
	}

	// 1 (Target Id: 0)
	v := strings.Replace(val, "(Target Id:", "", -1)
	v = strings.Replace(v, ")", "", -1)
	flds := strings.Fields(v)
	if len(flds) != 2 {
		err = fmt.Errorf("invalid format")
		return
	}

	grp, err = strconv.Atoi(flds[0])
	if err != nil {
		return
	}
	if grp < 0 {
		err = fmt.Errorf("invalid format")
		return
	}

	tgt, err = strconv.Atoi(flds[1])
	if err != nil {
		return
	}
	if tgt < 0 {
		err = fmt.Errorf("invalid format")
	}

	return
}

func parseRAIDLv(val string) raidcli.Level {
	if lv, ok := lvmap[val]; ok {
		return lv
	}
	return raidcli.Unknown
}

func parseCachePolicy(val string) string {
	flds := strings.Split(val, ",")
	if len(flds) == 0 {
		return "unknown"
	}

	pol := strings.TrimSpace(flds[0])

	if !strings.Contains(strings.ToLower(pol), "write") {
		return "unknown"
	}

	return pol
}
