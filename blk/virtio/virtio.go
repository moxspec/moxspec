package virtio

import (
	"os"
	"path/filepath"

	"github.com/moxspec/moxspec/blk"
	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("virtio")
}

// Controller represents a controller
type Controller struct {
	Path  string
	Disks []*Disk
}

// Disk represents a SATA disk
type Disk struct {
	blk.CommonSpec
	Name   string
	Path   string
	Driver string
}

func (d Disk) blkDir() string {
	return filepath.Join(d.Path, "block", d.Name)
}

// NewDecoder creates and initializes an Controller
func NewDecoder(path string) *Controller {
	c := new(Controller)
	c.Path = path
	return c
}

// Decode searches connected disks
func (c *Controller) Decode() error {
	for _, vioDir := range util.FilterPrefixedDirs(c.Path, "virtio") {
		d := new(Disk)
		d.Name = blk.GetBlockName(vioDir)
		d.Path = vioDir
		d.CommonSpec = *blk.NewCommonSpec(d.blkDir())

		drvLink, err := os.Readlink(filepath.Join(vioDir, "driver"))
		if err == nil {
			d.Driver = filepath.Base(drvLink)
		}

		c.Disks = append(c.Disks, d)
		log.Debugf("virtio disk: %+v", d)
	}

	return nil
}
