package nvme

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/moxspec/moxspec/blk"
	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("nvme")
}

// Controller represents a NVMe controller
type Controller struct {
	Path        string
	Name        string
	Model       string
	Serial      string
	FirmwareRev string
	Namespaces  []*Namespace
}

func (c Controller) nvmeDir() string {
	return filepath.Join(c.Path, "nvme", c.Name)
}

// Namespace represents a NVMe namespace
type Namespace struct {
	blk.CommonSpec
	Name string
	Path string
}

// ID returns namespace id
func (n Namespace) ID() int {
	return parseNamespaceID(n.Name)
}

func parseNamespaceID(name string) int {
	if !strings.HasPrefix(name, "nvme") {
		return 0
	}

	elms := strings.Split(strings.Replace(name, "nvme", "", 1), "n")
	if len(elms) != 2 {
		return 0
	}

	i, err := strconv.Atoi(elms[1])
	if err != nil {
		return 0
	}
	return i
}

// NewDecoder creates and initializes a Controller
func NewDecoder(path string) *Controller {
	c := new(Controller)
	c.Path = path
	return c
}

// Decode searches connected disks
func (c *Controller) Decode() error {
	name, err := getDeviceName(c.Path)
	if err != nil {
		return err
	}
	c.Name = name
	c.Model, _ = util.LoadString(filepath.Join(c.nvmeDir(), "model"))
	c.Serial, _ = util.LoadString(filepath.Join(c.nvmeDir(), "serial"))
	c.FirmwareRev, _ = util.LoadString(filepath.Join(c.nvmeDir(), "firmware_rev"))

	// scan namespaces
	for _, nsDir := range util.FilterPrefixedDirs(c.nvmeDir(), c.Name) {
		ns := new(Namespace)
		ns.Name = filepath.Base(nsDir)
		ns.Path = nsDir
		ns.CommonSpec = *blk.NewCommonSpec(nsDir)

		c.Namespaces = append(c.Namespaces, ns)
		log.Debugf("nvme namespace: %+v", ns)
	}

	return nil
}

func getDeviceName(path string) (string, error) {
	nvmeDir := filepath.Join(path, "nvme")
	if !util.Exists(nvmeDir) {
		return "", fmt.Errorf("not found nvme directory in %s", path)
	}

	dirs := util.FilterPrefixedDirs(nvmeDir, "nvme")
	if len(dirs) != 1 {
		return "", fmt.Errorf("unexpected nvme directory topology: %s", path)
	}

	return filepath.Base(dirs[0]), nil
}
