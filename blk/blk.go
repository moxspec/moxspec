package blk

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("blk")
}

// CommonSpec represents a general disk spec
type CommonSpec struct {
	Major        uint16
	Minor        uint16
	Model        string
	Blocks       uint64
	LogBlockSize uint16
	PhyBlockSize uint16
	Scheduler    string
}

// Size returns the byte size of the device
func (c CommonSpec) Size() uint64 {
	return uint64(c.Blocks) * uint64(c.LogBlockSize)
}

// NewCommonSpec creates and initializes a common spec
func NewCommonSpec(path string) *CommonSpec {
	c := new(CommonSpec)
	c.Blocks, _ = util.LoadUint64(filepath.Join(path, "size"))
	c.Model, _ = util.LoadString(filepath.Join(path, "device", "model"))
	c.LogBlockSize, _ = util.LoadUint16(filepath.Join(path, "queue", "logical_block_size"))
	c.PhyBlockSize, _ = util.LoadUint16(filepath.Join(path, "queue", "physical_block_size"))

	num, err := util.LoadString(filepath.Join(path, "dev"))
	if err == nil {
		c.Major, c.Minor = parseNodeNumber(num)
	}

	sched, err := util.LoadString(filepath.Join(path, "queue", "scheduler"))
	if err == nil {
		c.Scheduler = parseScheduler(sched)
	}

	return c
}

// SCSIAddress represents a SCSI address
type SCSIAddress struct {
	Host    uint16
	Channel uint16
	Target  uint16
	Lun     uint16
}

// NewSCSIAddress creates and initializes a SCSIAddress
func NewSCSIAddress(path string) *SCSIAddress {
	s := new(SCSIAddress)

	addr := filepath.Base(path)
	s.Host, s.Channel, s.Target, s.Lun = parseSCSIAddress(addr)

	return s
}

func parseSCSIAddress(addr string) (host, ch, tgt, lun uint16) {
	if !strings.Contains(addr, ":") {
		return
	}

	flds := strings.Split(addr, ":")
	if len(flds) != 4 {
		return
	}

	h, err := strconv.Atoi(flds[0])
	if err != nil || h < 0 {
		return
	}

	c, err := strconv.Atoi(flds[1])
	if err != nil || c < 0 {
		return
	}

	t, err := strconv.Atoi(flds[2])
	if err != nil || t < 0 {
		return
	}

	l, err := strconv.Atoi(flds[3])
	if err != nil || l < 0 {
		return
	}

	host = uint16(h)
	ch = uint16(c)
	tgt = uint16(t)
	lun = uint16(l)
	return
}

// parseNodeNumber returns major/minor number from given string
func parseNodeNumber(str string) (uint16, uint16) {
	elm := strings.Split(str, ":")
	if len(elm) != 2 {
		log.Warnf("invalid format was given. %s.", str)
		return 0, 0
	}

	major, err := strconv.ParseInt(elm[0], 10, 16)
	if err != nil {
		log.Warn(err.Error())
		return 0, 0
	}
	minor, err := strconv.ParseInt(elm[1], 10, 16)
	if err != nil {
		log.Warn(err.Error())
		return 0, 0
	}

	return uint16(major), uint16(minor)
}

// GetBlockName scans the block device name from given path and returns it
func GetBlockName(path string) string {
	bdir := filepath.Join(path, "block")
	if !util.Exists(bdir) {
		return ""
	}

	files, err := ioutil.ReadDir(bdir)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		return file.Name()
	}

	return ""
}

// parseScheduler returns selected scheduler from given kernel representation
func parseScheduler(list string) string {
	if list == "" {
		return ""
	}

	scheds := strings.Fields(list)
	if len(scheds) == 1 && !strings.Contains(scheds[0], "[") && !strings.Contains(scheds[0], "]") {
		return scheds[0]
	}

	for _, s := range scheds {
		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			sched := s
			sched = strings.Replace(sched, "[", "", 1)
			sched = strings.Replace(sched, "]", "", 1)
			return sched
		}
	}

	return ""
}
