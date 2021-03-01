package raid

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moxspec/moxspec/blk"
	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/scsi"
	"github.com/moxspec/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("raid")
}

// Expander represents an expander card (raid card)
type Expander struct {
	Path         string
	VirtualDisks []*VirtualDisk
}

// VirtualDisk represents a raid volume
// if a disk is configured as a pass-through mode it will have sas_address
// sas_address will used to bind a disk into topology decoded from raidcli outputs
type VirtualDisk struct {
	blk.CommonSpec
	blk.SCSIAddress
	Name         string
	Path         string
	Driver       string
	SASAddress   string // eg. 0x5003048009282600
	WWN          string
	SerialNumber string
}

func (d VirtualDisk) blkDir() string {
	return filepath.Join(d.Path, "block", d.Name)
}

// NewDecoder creates and initializes an Expander
func NewDecoder(path string) *Expander {
	e := new(Expander)
	e.Path = path
	return e
}

// Decode searches connected disks
func (e *Expander) Decode() error {
	bpaths := readExpandedBlockPaths()

	for _, hostDir := range util.FilterPrefixedDirs(e.Path, "host") {
		log.Debugf("host: %s", hostDir)

		for _, diskPath := range scanDiskPath(hostDir, bpaths) {
			log.Debugf("disk: %s", diskPath)

			d := new(VirtualDisk)
			d.Name = blk.GetBlockName(diskPath)
			d.Path = diskPath
			d.CommonSpec = *blk.NewCommonSpec(d.blkDir())
			d.SCSIAddress = *blk.NewSCSIAddress(diskPath)

			d.SASAddress, _ = util.LoadString(filepath.Join(diskPath, "sas_address"))
			if d.SASAddress != "" {
				log.Debugf("this device has sas_address: %s", d.SASAddress)
			}

			drvLink, err := os.Readlink(filepath.Join(diskPath, "driver"))
			if err == nil {
				d.Driver = filepath.Base(drvLink)
			}

			d.WWN = scsi.DecodeWWN(filepath.Join(diskPath, "vpd_pg83"))
			d.SerialNumber = scsi.DecodeSerialNumber(filepath.Join(diskPath, "vpd_pg80"))

			e.VirtualDisks = append(e.VirtualDisks, d)

			log.Debugf("virtual disk: %+v", d)
		}
	}

	return nil
}

func readExpandedBlockPaths() []string {
	var bpaths []string

	syspath := "/sys/class/block"
	files, err := ioutil.ReadDir(syspath)
	if err != nil {
		log.Warnf("could not read dir: %s", syspath)
		return nil
	}

	for _, file := range files {
		lpath := filepath.Join(syspath, file.Name())
		bpath, err := filepath.EvalSymlinks(lpath)
		if err != nil {
			log.Warnf("could not read link: %s", lpath)
			continue
		}

		if util.Exists(filepath.Join(bpath, "partition")) {
			log.Debugf("%s is a partition", bpath)
			continue
		}

		// e.g:
		//   /sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0/host0/target0:2:0/0:2:0:0/block/sda
		//   /sys/devices/pci0000:60/0000:60:03.1/0000:61:00.0/host0/port-0:0/expander-0:0/port-0:0:0/end_device-0:0:0/target0:0:0/0:0:0:0/block/sda
		if !strings.Contains(bpath, "/pci") || !strings.Contains(bpath, "/host") || !strings.Contains(bpath, "/target") {
			log.Debugf("%s is not under an raid expander", bpath)
			continue
		}

		elms := strings.Split(bpath, "/block/")
		if len(elms) != 2 {
			log.Warnf("could not parse %s", bpath)
			continue
		}

		// elms[0] should be SCSI target device representation
		// e.g:
		//   /sys/devices/pci0000:00/0000:00:01.0/0000:01:00.0/host0/target0:2:0/0:2:0:0
		log.Debugf("found virtual disk: %s", elms[0])
		bpaths = append(bpaths, elms[0])
	}

	return bpaths
}

func scanDiskPath(path string, bpaths []string) []string {
	var disks []string

	for _, bpath := range bpaths {
		if !strings.Contains(bpath, path) {
			continue
		}

		log.Debugf("%s has %s", path, bpath)
		disks = append(disks, bpath)
	}

	return disks
}
