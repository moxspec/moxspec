package edac

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("edac")
}

// NewDecoder creates and initializes a Topology as Decoder
func NewDecoder() *Topology {
	return newTopology("/sys/devices/system/edac/mc")
}

// Topology represents a edac topology
type Topology struct {
	path        string
	Controllers []*MemoryController
}

// Decode makes Topology satisfy the mox.Decoder interface
func (t *Topology) Decode() error {
	for _, mcdir := range util.FilterPrefixedDirs(t.path, "mc") {
		log.Debugf("scanning %s", mcdir)
		mc := newMemoryController(mcdir)
		err := mc.decode()
		if err != nil {
			log.Debug(err)
		}
		t.Controllers = append(t.Controllers, mc)
	}
	return nil
}

func newTopology(path string) *Topology {
	t := new(Topology)
	t.path = path
	return t
}

// MemoryController represents a memory controller
//
// Abbrev:
//   CE: Correctable Error
//       This count is very important to examine.
//       CEs provide early indications that a DIMM is beginning to fail.
//   UE: Uncorrectable Error
//
// NoInfo:
//  no information as to which DIMM slot is having errors
type MemoryController struct {
	Path          string
	Size          uint64
	Name          string
	CECount       uint64
	CENoInfoCount uint64
	UECount       uint64
	UENoInfoCount uint64
	CSRows        []*ChipSelectRow
}

func (m *MemoryController) decode() error {
	var err error

	mcname, err := util.LoadString(filepath.Join(m.Path, "mc_name"))
	if err != nil {
		return err
	}
	m.Name = mcname

	sizemb, err := util.LoadUint64(filepath.Join(m.Path, "size_mb"))
	if err != nil {
		return err
	}
	m.Size = sizemb * 1000 * 1000

	m.CECount, _ = util.LoadUint64(filepath.Join(m.Path, "ce_count"))
	m.CENoInfoCount, _ = util.LoadUint64(filepath.Join(m.Path, "ce_noinfo_count"))
	m.UECount, _ = util.LoadUint64(filepath.Join(m.Path, "ue_count"))
	m.UENoInfoCount, _ = util.LoadUint64(filepath.Join(m.Path, "ue_noinfo_count"))

	for _, rowdir := range util.FilterPrefixedDirs(m.Path, "csrow") {
		csrow := newChipSelectRow(rowdir)
		err := csrow.decode()
		if err != nil {
			return err
		}
		m.CSRows = append(m.CSRows, csrow)
	}

	return nil
}

func newMemoryController(path string) *MemoryController {
	m := new(MemoryController)
	m.Path = path
	return m
}

// ChipSelectRow represents a Chip-Select Row (csrowX)
type ChipSelectRow struct {
	Path     string
	Name     string
	Size     uint64
	CECount  uint64
	UECount  uint64
	Channels []*Channel
}

func (c *ChipSelectRow) decode() error {
	c.Name = filepath.Base(c.Path)

	sizemb, err := util.LoadUint64(filepath.Join(c.Path, "size_mb"))
	if err != nil {
		return err
	}
	c.Size = sizemb * 1000 * 1000
	c.CECount, _ = util.LoadUint64(filepath.Join(c.Path, "ce_count"))
	c.UECount, _ = util.LoadUint64(filepath.Join(c.Path, "ue_count"))

	log.Debugf("%s: size=%d, ce=%d, ue=%d", c.Name, c.Size, c.CECount, c.UECount)

	chs := make(map[string][]string)
	for _, ch := range util.FilterPrefixedFiles(c.Path, "ch") {
		fname := filepath.Base(ch)

		flds := strings.Split(fname, "_")
		if len(flds) != 3 {
			log.Debugf("%s is ignored", ch)
			continue
		}

		chs[flds[0]] = append(chs[flds[0]], ch)
	}

	for key, files := range chs {
		ch := new(Channel)
		ch.Name = key

		for _, file := range files {
			if strings.HasSuffix(file, "_ce_count") {
				ch.CECount, _ = util.LoadUint64(file)
			} else if strings.HasSuffix(file, "_label") {
				ch.Label, _ = util.LoadString(file)
			}
		}

		log.Debugf("%s: label=%s, ce=%d", ch.Name, ch.Label, ch.CECount)
		c.Channels = append(c.Channels, ch)
	}

	sort.Slice(c.Channels, func(i, j int) bool {
		return c.Channels[i].Name < c.Channels[j].Name
	})

	return nil
}

func newChipSelectRow(path string) *ChipSelectRow {
	c := new(ChipSelectRow)
	c.Path = path
	return c
}

// Channel represents a Channel table (chX)
type Channel struct {
	Name    string
	Label   string
	CECount uint64
}
