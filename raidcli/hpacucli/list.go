package hpacucli

import (
	"strings"
)

func parseCtlList(in string) ([]*Controller, error) {
	var ctls []*Controller
	for _, line := range strings.Split(in, "\n") {
		l := strings.TrimSpace(line)

		if l == "" {
			continue
		}

		if !strings.Contains(l, " in Slot ") {
			continue
		}

		spls := strings.Split(l, " in Slot ")
		if len(spls) != 2 {
			continue
		}

		flds := strings.Fields(spls[1])
		if len(flds) < 1 {
			continue
		}

		c := new(Controller)
		c.ProductName = spls[0]
		c.Slot = flds[0]
		ctls = append(ctls, c)
	}

	return ctls, nil
}
