package pci

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/actapio/moxspec/util"
)

const (
	fmtVendor                = "v%08x"
	fmtVendorDevice          = "v%08xd%08x"
	fmtVendorDeviceSubsystem = "v%08xd%08xsv%08xsd%08x"
	fmtClass                 = "c%04x"
	fmtClassSubclass         = "c%04xsc%04x"
	fmtClassSubclassIntf     = "c%04xsc%04xi%04x"
)

type (
	vendorDB map[string]string
	classDB  map[string]string
)

func (v vendorDB) vendorName(vid uint16) string {
	return util.ShortenVendorName(v[fmt.Sprintf(fmtVendor, vid)])
}

func (v vendorDB) deviceName(vid, devid uint16) string {
	return v[fmt.Sprintf(fmtVendorDevice, vid, devid)]
}

func (v vendorDB) subSysName(vid, devid, ssvid, ssdevid uint16) string {
	return v[fmt.Sprintf(fmtVendorDeviceSubsystem, vid, devid, ssvid, ssdevid)]
}

func (c classDB) className(cid byte) string {
	return c[fmt.Sprintf(fmtClass, cid)]
}

func (c classDB) subClassName(cid, scid byte) string {
	return c[fmt.Sprintf(fmtClassSubclass, cid, scid)]
}

func (c classDB) intfaceName(cid, scid, intfid byte) string {
	return c[fmt.Sprintf(fmtClassSubclassIntf, cid, scid, intfid)]
}

var (
	vdb vendorDB
	cdb classDB
)

func initDB(path string) error {
	pciids, _ := ioutil.ReadFile(path)

	// first half = vendor, device, subsystem
	// second half = class, sub_class, programming if
	blocks := strings.Split(string(pciids), "List of known device classes, subclasses and programming interfaces")
	if len(blocks) != 2 {
		fmt.Println("bad pci.ids format")
		os.Exit(1)
	}

	var err error
	vdb, err = genVendorDB(blocks[0])
	if err != nil {
		return err
	}

	cdb, err = genClassDB(blocks[1])
	return err
}

func genVendorDB(lines string) (map[string]string, error) {
	vdb := make(map[string]string)

	var vID, dID uint64
	for _, l := range strings.Split(lines, "\n") {
		if strings.HasPrefix(l, "#") {
			continue
		}

		if len(l) == 0 {
			continue
		}

		indent := strings.Count(l, "\t")
		flds := strings.Fields(l)

		if indent == 0 {
			id, _ := parseHexStr(flds[0])
			name := strings.Join(flds[1:], " ")
			vdb[fmt.Sprintf(fmtVendor, id)] = name
			vID = id
		}

		if indent == 1 {
			id, _ := parseHexStr(flds[0])
			name := strings.Join(flds[1:], " ")
			vdb[fmt.Sprintf(fmtVendorDevice, vID, id)] = name
			dID = id
		}

		if indent == 2 {
			svID, _ := parseHexStr(flds[0])
			sbID, _ := parseHexStr(flds[1])
			name := strings.Join(flds[2:], " ")
			vdb[fmt.Sprintf(fmtVendorDeviceSubsystem, vID, dID, svID, sbID)] = name
		}
	}

	return vdb, nil
}

func genClassDB(lines string) (map[string]string, error) {
	cdb := make(map[string]string)

	var cID, scID uint64
	for _, l := range strings.Split(lines, "\n") {
		if strings.HasPrefix(l, "#") {
			continue
		}

		if len(l) == 0 {
			continue
		}

		indent := strings.Count(l, "\t")
		flds := strings.Fields(l)

		if indent == 0 {
			id, _ := parseHexStr(flds[1])
			name := strings.Join(flds[2:], " ")
			cdb[fmt.Sprintf(fmtClass, id)] = name
			cID = id
		}

		if indent == 1 {
			id, _ := parseHexStr(flds[0])
			name := strings.Join(flds[1:], " ")
			cdb[fmt.Sprintf(fmtClassSubclass, cID, id)] = name
			scID = id
		}

		if indent == 2 {
			id, _ := parseHexStr(flds[0])
			name := strings.Join(flds[1:], " ")
			cdb[fmt.Sprintf(fmtClassSubclassIntf, cID, scID, id)] = name
		}
	}

	return cdb, nil
}
