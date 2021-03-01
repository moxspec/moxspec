package megacli

import (
	"fmt"
	"strings"

	"github.com/moxspec/moxspec/raidcli"
)

func setAdpInfo(c *Controller) error {
	res, err := raidcli.Run(clipath, "-AdpAllInfo", fmt.Sprintf("-a%d", c.Number), "-NoLog")
	if err != nil {
		return err
	}

	name, sn, bios, fw, bbu, err := parseAdpInfo(res)
	if err != nil {
		return err
	}

	c.ProductName = name
	c.SerialNumber = sn
	c.BIOS = bios
	c.Firmware = fw
	c.Battery = bbu
	return nil
}

func parseAdpInfo(in string) (name, sn, bios, fw string, bbu bool, err error) {
	for _, line := range strings.Split(in, "\n") {
		key, val, err := raidcli.SplitKeyVal(line, ":")
		if err != nil {
			log.Debug(err)
			continue
		}

		switch key {
		case "Product Name":
			name = val
		case "Serial No":
			sn = val
		case "BIOS Version":
			bios = val
		case "FW Version":
			fw = val
		case "BBU":
			if strings.ToLower(val) == "present" {
				bbu = true
			}
		}

		if key == "Current Time" {
			break
		}
	}

	return
}
