package cpu

import "sort"

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
}

// ThreadCount returns a thread count
func (c Core) ThreadCount() int {
	return len(c.Threads)
}
