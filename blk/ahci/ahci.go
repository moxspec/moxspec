package ahci

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/actapio/moxspec/blk"
	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("ahci")
}

// HostBusAdapter represents a AHCI Host Bus Adapter
type HostBusAdapter struct {
	Path  string
	Disks []*Disk
}

// Disk represents a SATA disk
type Disk struct {
	blk.CommonSpec
	blk.SCSIAddress
	Name       string
	Path       string
	Driver     string
	Rotational bool
}

func (d Disk) blkDir() string {
	return filepath.Join(d.Path, "block", d.Name)
}

// NewDecoder creates and initializes a HostBusAdapter
func NewDecoder(path string) *HostBusAdapter {
	h := new(HostBusAdapter)
	h.Path = path
	return h
}

// Decode searches connected disks
func (h *HostBusAdapter) Decode() error {
	for _, ataDir := range util.FilterPrefixedDirs(h.Path, "ata") {
		log.Debugf("ahci port: %s", ataDir)

		diskPath := scanDiskPathFromATADir(ataDir)
		if diskPath == "" {
			continue
		}

		log.Debugf("ahci disk: %s", diskPath)
		d := newDisk(diskPath)

		h.Disks = append(h.Disks, d)
		log.Debugf("sata disk: %+v", d)
	}

	if len(h.Disks) > 0 {
		return nil
	}

	// fall-back
	for _, hostDir := range util.FilterPrefixedDirs(h.Path, "host") {
		log.Debugf("host: %s", hostDir)

		diskPath := scanDiskPathFromHostDir(hostDir)
		if diskPath == "" {
			continue
		}

		log.Debugf("ahci disk: %s", diskPath)
		d := newDisk(diskPath)

		h.Disks = append(h.Disks, d)
		log.Debugf("sata disk: %+v", d)
	}

	return nil
}

func newDisk(diskPath string) *Disk {
	d := new(Disk)
	d.Name = blk.GetBlockName(diskPath)
	d.Path = diskPath
	d.CommonSpec = *blk.NewCommonSpec(d.blkDir())
	d.SCSIAddress = *blk.NewSCSIAddress(diskPath)

	drvLink, err := os.Readlink(filepath.Join(diskPath, "driver"))
	if err == nil {
		d.Driver = filepath.Base(drvLink)
	}

	if r, _ := util.LoadUint64(filepath.Join(d.blkDir(), "queue", "rotational")); r == 1 {
		d.Rotational = true
	}

	return d
}

func scanDiskPathFromATADir(path string) string {
	fs := util.FilterPrefixedDirs(path, "host")
	if len(fs) != 1 {
		return ""
	}
	hostPath := fs[0]

	return scanDiskPathFromHostDir(hostPath)
}

func scanDiskPathFromHostDir(path string) string {
	fs := util.FilterPrefixedDirs(path, "target")
	if len(fs) != 1 {
		return ""
	}
	tgtPath := fs[0]

	// FROM: /sys/devices/pci0000:00/0000:00:11.5/ata3/host2/target2:0:0
	// TO:   2:0:0
	scsiPrefix := strings.Replace(filepath.Base(tgtPath), "target", "", 1)

	fs = util.FilterPrefixedDirs(tgtPath, scsiPrefix)
	if len(fs) != 1 {
		return ""
	}
	scsiPath := fs[0]

	return scsiPath
}
