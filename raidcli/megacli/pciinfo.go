package megacli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actapio/moxspec/raidcli"
)

func parsePCIInfo(info string) ([]*Controller, error) {
	if info == "" {
		return nil, fmt.Errorf("empty outputs")
	}

	var ctls []*Controller

	var ctl *Controller
	for _, line := range strings.Split(info, "\n") {
		l := strings.TrimSpace(line)

		if strings.HasPrefix(l, "PCI information for Controller") {
			if ctl != nil {
				ctls = append(ctls, ctl)
			}

			ctl = new(Controller)
			num, err := parseControllerNumber(l)
			if err != nil {
				log.Debug(err)
				continue
			}
			ctl.Number = num
		}

		if ctl == nil {
			continue
		}

		key, val, err := raidcli.SplitKeyVal(line, ":")
		if err != nil {
			log.Debug(err)
			continue
		}

		num, err := strconv.ParseInt(val, 16, 64)
		if err != nil {
			continue
		}
		n := uint32(num)

		switch key {
		case "Bus Number":
			ctl.Bus = n
		case "Device Number":
			ctl.Device = n
		case "Function Number":
			ctl.Function = n
		}
	}

	// finalize
	if ctl != nil {
		ctls = append(ctls, ctl)
	}

	return ctls, nil
}

func parseControllerNumber(line string) (int, error) {
	flds := strings.Fields(line)
	if len(flds) < 1 {
		return -1, fmt.Errorf("invalid format")
	}

	num, err := strconv.Atoi(flds[len(flds)-1])
	if err != nil {
		return -1, err
	}

	if num < 0 {
		return -1, fmt.Errorf("invalid number")
	}

	return num, nil
}
