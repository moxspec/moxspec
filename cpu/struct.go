package cpu

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/moxspec/moxspec/util"
)

// Topology represents a processor topology
type Topology struct {
	path     string
	packages map[uint16]*Package
}

func (t Topology) hasPackage(id uint16) bool {
	p, ok := t.packages[id]
	return (ok && p != nil)
}

// Packages returns sorted package list
func (t Topology) Packages() []*Package {
	var l []*Package
	for _, p := range t.packages {
		l = append(l, p)
	}

	sort.Slice(l, func(i, j int) bool {
		return (l[i].ID < l[j].ID)
	})

	return l
}

// Package represents a processor package (socket)
// NOTE: MCM processor may have multiple nodes per package
//       e.g: Zen: up to four nodes per package
type Package struct {
	ID            uint16
	ThrottleCount uint16
	nodes         map[uint16]*Node
}

func (p Package) hasNode(id uint16) bool {
	n, ok := p.nodes[id]
	return (ok && n != nil)
}

// NodeCount returns a node count
func (p Package) NodeCount() int {
	return len(p.nodes)
}

// CoreCount returns a core count
func (p Package) CoreCount() int {
	var sum int
	for _, n := range p.nodes {
		sum = sum + n.CoreCount()
	}
	return sum
}

// ThreadCount returns a thread count
func (p Package) ThreadCount() int {
	var sum int
	for _, n := range p.nodes {
		sum = sum + n.ThreadCount()
	}
	return sum
}

// Nodes returns sorted node list
func (p Package) Nodes() []*Node {
	var l []*Node
	for _, n := range p.nodes {
		l = append(l, n)
	}

	sort.Slice(l, func(i, j int) bool {
		return (l[i].ID < l[j].ID)
	})

	return l
}

// Node represents a numa node
type Node struct {
	ID    uint16
	cores map[uint16]*Core
}

func (n Node) hasCore(id uint16) bool {
	c, ok := n.cores[id]
	return (ok && c != nil)
}

// CoreCount returns a core count
func (n Node) CoreCount() int {
	return len(n.cores)
}

// ThreadCount returns a thread count
func (n Node) ThreadCount() int {
	var sum int
	for _, c := range n.cores {
		sum = sum + c.ThreadCount()
	}
	return sum
}

// Cores returns sorted core list
func (n Node) Cores() []*Core {
	var l []*Core
	for _, n := range n.cores {
		l = append(l, n)
	}

	sort.Slice(l, func(i, j int) bool {
		return (l[i].ID < l[j].ID)
	})

	return l
}

// Core represents a physical core
type Core struct {
	ID            uint16
	ThrottleCount uint16
	Threads       []uint16
	BaseFreq      uint64
	MaxFreq       uint64
	MinFreq       uint64
	Scaling       struct {
		Driver             string
		Governor           string
		AvailableGovernors []string
		CurFreq            uint64
		MaxFreq            uint64
		MinFreq            uint64
	}
}

// ThreadCount returns a thread count
func (c Core) ThreadCount() int {
	return len(c.Threads)
}

func (c *Core) decode(cpudir string) error {
	list := filepath.Join(filepath.Join(cpudir, "topology"), "thread_siblings_list")
	ls, err := util.LoadString(list)
	if err != nil {
		return fmt.Errorf("could not load %s", list)
	}
	c.Threads = parseListString(ls)

	c.ThrottleCount, _ = util.LoadUint16(filepath.Join(cpudir, "thermal_throttle", "core_throttle_count"))
	c.BaseFreq, _ = util.LoadUint64(filepath.Join(cpudir, "cpufreq", "base_frequency"))
	c.MaxFreq, _ = util.LoadUint64(filepath.Join(cpudir, "cpufreq", "cpuinfo_max_freq"))
	c.MinFreq, _ = util.LoadUint64(filepath.Join(cpudir, "cpufreq", "cpuinfo_min_freq"))
	c.Scaling.CurFreq, _ = util.LoadUint64(filepath.Join(cpudir, "cpufreq", "scaling_cur_freq"))
	c.Scaling.MaxFreq, _ = util.LoadUint64(filepath.Join(cpudir, "cpufreq", "scaling_max_freq"))
	c.Scaling.MinFreq, _ = util.LoadUint64(filepath.Join(cpudir, "cpufreq", "scaling_min_freq"))
	c.Scaling.Governor, _ = util.LoadString(filepath.Join(cpudir, "cpufreq", "scaling_governor"))
	c.Scaling.Driver, _ = util.LoadString(filepath.Join(cpudir, "cpufreq", "scaling_driver"))
	govs, _ := util.LoadString(filepath.Join(cpudir, "cpufreq", "scaling_available_governors"))
	c.Scaling.AvailableGovernors = strings.Fields(govs)

	return nil
}
