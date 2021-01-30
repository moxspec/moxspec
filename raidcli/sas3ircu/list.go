package sas3ircu

import (
	"fmt"
	"strconv"
	"strings"
)

func parseCtlList(in string) ([]*Controller, error) {
	var ctls []*Controller

	for _, line := range strings.Split(in, "\n") {
		log.Debug(line)
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		flds := strings.Fields(l)
		if len(flds) != 7 {
			continue
		}

		num, err := strconv.Atoi(flds[0])
		if err != nil {
			continue
		}

		dom, bus, dev, fun, err := parsePCIAddr(flds[4])
		if err != nil {
			continue
		}

		c := new(Controller)
		c.Number = num
		c.Domain = dom
		c.Bus = bus
		c.Device = dev
		c.Function = fun
		c.ptAddrMap = make(map[string]*PhyDrive)

		ctls = append(ctls, c)
	}

	return ctls, nil
}

func parsePCIAddr(addr string) (dom, bus, dev, fun uint32, err error) {
	flds := strings.Split(strings.ToLower(addr), ":")
	if len(flds) != 4 {
		err = fmt.Errorf("%s is not pci address", addr)
		return
	}

	do, err := strconv.ParseInt(strings.TrimSuffix(flds[0], "h"), 16, 64)
	if err != nil {
		return
	}
	dom = uint32(do)

	bu, err := strconv.ParseInt(strings.TrimSuffix(flds[1], "h"), 16, 64)
	if err != nil {
		return
	}
	bus = uint32(bu)

	de, err := strconv.ParseInt(strings.TrimSuffix(flds[2], "h"), 16, 64)
	if err != nil {
		return
	}
	dev = uint32(de)

	fu, err := strconv.ParseInt(strings.TrimSuffix(flds[3], "h"), 16, 64)
	if err != nil {
		return
	}
	fun = uint32(fu)
	return
}
