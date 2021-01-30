package eth

import "fmt"

// Module represents the module
type Module struct {
	FormFactor   string
	VendorName   string
	ProductName  string
	SerialNumber string
	Connector    string
	CableLength  uint8
}

type moduleInfo struct {
	cmd        uint32
	kind       uint32
	eepromSize uint32
}

func (m moduleInfo) code() uint32 {
	return m.cmd
}

func (m moduleInfo) dump() string {
	return fmt.Sprintf("%+v", m)
}

func newModuleInfo() *moduleInfo {
	m := new(moduleInfo)
	m.cmd = ethtoolGetModuleInfo
	return m
}

func getModuleInfo(ndev *netdev) (*moduleInfo, error) {
	log.Debug("getting module info")

	mi := newModuleInfo()
	errno, err := ndev.post(mi)
	if err != nil {
		return nil, err
	}

	if errno < 0 {
		return nil, fmt.Errorf("err: %d", errno)
	}

	if !isValidROMSize(mi.eepromSize) {
		log.Debugf("invalid eeprom size: %d bytes", mi.eepromSize)
		return nil, fmt.Errorf("no eeprom available")
	}

	return mi, nil
}

func getModule(ndev *netdev) (*Module, error) {
	mi, err := getModuleInfo(ndev)
	if err != nil {
		return nil, err
	}

	ep, err := dumpEEPROM(ndev, mi.kind, mi.eepromSize)
	if err != nil {
		return nil, err
	}

	var md *Module
	log.Debug("detecting eeprom type")
	switch mi.kind {
	case sff8079:
		log.Debug("SFF-8079 (SFP)")
		md, err = decodeSFF8079(ep.data)
	case sff8472:
		log.Debug("SFF-8472 (SFP+/SFP28)")
		md, err = decodeSFF8472(ep.data)
	case sff8636:
		log.Debug("SFF-8636 (QSFP+/QSFP28)")
		log.Debug("nothing to do")
	case sff8436:
		log.Debug("SFF-8436 (QSFP+)")
		log.Debug("nothing to do")
	}

	return md, err
}
