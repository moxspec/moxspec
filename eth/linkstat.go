package eth

import "fmt"

type linkStat struct {
	cmd  uint32
	data uint32
}

func (l linkStat) code() uint32 {
	return l.cmd
}

func (l linkStat) dump() string {
	return fmt.Sprintf("%+v", l)
}

func getLinkStat(ndev *netdev) (bool, error) {
	log.Debug("getting link status")

	l := new(linkStat)
	l.cmd = ethtoolGetLink

	errno, err := ndev.post(l)
	if err != nil {
		return false, err
	}

	if errno < 0 {
		return false, fmt.Errorf("err: %d", errno)
	}

	if l.data > 0 {
		return true, nil
	}

	return false, nil
}
