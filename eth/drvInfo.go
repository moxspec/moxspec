package eth

import (
	"fmt"
)

type drvInfo struct {
	cmd         uint32
	driver      [32]byte
	version     [32]byte
	fwVersion   [32]byte
	busInfo     [32]byte
	reserved1   [32]byte
	reserved2   [32]byte
	nPrivFlags  uint32
	nStats      uint32
	testinfoLen uint32
	eedumpLen   uint32
	regdumpLen  uint32
}

func (d drvInfo) code() uint32 {
	return d.cmd
}

func (d drvInfo) dump() string {
	return fmt.Sprintf("%+v", d)
}

func (d drvInfo) firmwareVersion() string {
	return parseFirmwareVersion(d.fwVersion)
}

func parseFirmwareVersion(buf [32]byte) string {
	last := len(buf)
	for i, c := range buf {
		if c == 0x00 {
			last = i
			break
		}
	}

	return string(buf[:last])
}

func getDrvInfo(ndev *netdev) (*drvInfo, error) {
	log.Debug("getting driver info")

	d := new(drvInfo)
	d.cmd = ethtoolGetDriver

	errno, err := ndev.post(d)
	if err != nil {
		return nil, err
	}

	if errno < 0 {
		return nil, fmt.Errorf("err: %d", errno)
	}

	return d, nil
}
