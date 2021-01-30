package eth

import "fmt"

const (
	// NOTE:
	// SFF-8079 ... 256
	// SFF-8472 ... 512
	// SFF-8636 ... 256
	// SFF-8436 ... 256
	maxRomSize = 512
)

const (
	sff8079 = 0x01
	sff8472 = 0x02
	sff8636 = 0x03
	sff8436 = 0x04
)

type eeprom struct {
	cmd    uint32
	magic  uint32
	offset uint32
	size   uint32
	data   [maxRomSize]byte
}

func (e eeprom) code() uint32 {
	return e.cmd
}

func (e eeprom) dump() string {
	return fmt.Sprintf("%+v", e)
}

func newEEPROM() *eeprom {
	e := new(eeprom)
	e.cmd = ethtoolGetEeprom
	return e
}

func isValidROMSize(s uint32) bool {
	return (s != 0 && s <= maxRomSize)
}

func dumpEEPROM(ndev *netdev, kind uint32, size uint32) (*eeprom, error) {
	log.Debug("getting eeprom")

	e := newEEPROM()
	e.size = size
	errno, err := ndev.post(e)

	if err != nil {
		return nil, err
	}

	if errno < 0 {
		return nil, fmt.Errorf("err: %d", errno)
	}

	if e.size == 0 {
		return nil, fmt.Errorf("can not dump eeprom")
	}

	return e, nil

}
