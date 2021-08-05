package bonding

import (
	"fmt"

	"github.com/moxspec/moxspec/loglet"
	"github.com/vishvananda/netlink"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("netlink")
}

// findBondSlaves returns slave devs of a bond master
func findBondSlaves(index int) ([]string, error) {

	var slaves []string

	linklist, err := netlink.LinkList()
	if err != nil {
		return slaves, err
	}

	for _, link := range linklist {
		if link.Attrs().MasterIndex == index {
			slaves = append(slaves, link.Attrs().Name)
		}
	}

	return slaves, nil
}

func GetBondDevices() []string {
	var bonds []string

	linklist, err := netlink.LinkList()
	if err != nil {
		log.Debugf("netlink.LinkList() failed: %s", err)
		return bonds
	}

	for _, link := range linklist {
		if link.Type() != "bond" {
			continue
		}
		bonds = append(bonds, link.Attrs().Name)
	}

	return bonds
}

func getBondParameters(devlink netlink.Link) (*netlink.Bond, error) {

	if parameters, ok := devlink.(*netlink.Bond); ok {
		return parameters, nil
	}
	return nil, fmt.Errorf("%s is not a bonding device", devlink.Attrs().Name)

}
