package eth

import (
	"fmt"
	"sort"
)

const (
	portTP    = 0x00
	portAUI   = 0x01
	portMII   = 0x02
	portFIBRE = 0x03
	portBNC   = 0x04
	portDA    = 0x05
	portNONE  = 0xef
	portOTHER = 0xff
)

// https://elixir.bootlin.com/linux/v5.1-rc3/source/include/uapi/linux/ethtool.h
const (
	linkMode10baseTHalf = iota
	linkMode10baseTFull
	linkMode100baseTHalf
	linkMode100baseTFull
	linkMode1000baseTHalf
	linkMode1000baseTFull // 5
	linkModeAutoneg
	linkModeTP
	linkModeAUI
	linkModeMII
	linkModeFIBRE // 10
	linkModeBNC
	linkMode10000baseTFull
	linkModePause
	linkModeAsymPause
	linkMode2500baseXFull // 15
	linkModeBackplane
	linkMode1000baseKXFull
	linkMode10000baseKX4Full
	linkMode10000baseKRFull
	linkMode10000baseRFEC // 20
	linkMode20000baseMLD2Full
	linkMode20000baseKR2Full
	linkMode40000baseKR4Full
	linkMode40000baseCR4Full
	linkMode40000baseSR4Full // 25
	linkMode40000baseLR4Full
	linkMode56000baseKR4Full
	linkMode56000baseCR4Full
	linkMode56000baseSR4Full
	linkMode56000baseLR4Full // 30
	linkMode25000baseCRFull
	linkMode25000baseKRFull
	linkMode25000baseSRFull
	linkMode50000baseCR2Full
	linkMode50000baseKR2Full // 35
	linkMode100000baseKR4Full
	linkMode100000baseSR4Full
	linkMode100000baseCR4Full
	linkMode100000baseLR4ER4Full
	linkMode50000baseSR2Full // 40
	linkMode1000baseXFull
	linkMode10000baseCRFull
	linkMode10000baseSRFull
	linkMode10000baseLRFull
	linkMode10000baseLRMFull // 45
	linkMode10000baseERFull
	linkMode2500baseTFull
	linkMode5000baseTFull
	linkModeFECNONE
	linkModeFECRS // 50
	linkModeFECBASER
	linkMode50000baseKRFull
	linkMode50000baseSRFull
	linkMode50000baseCRFull
	linkMode50000baseLRERFRFull // 55
	linkMode50000baseDRFull
	linkMode100000baseKR2Full
	linkMode100000baseSR2Full
	linkMode100000baseCR2Full
	linkMode100000baseLR2ER2FR2Full // 60
	linkMode100000baseDR2Full
	linkMode200000baseKR4Full
	linkMode200000baseSR4Full // 63 ... uint64/int64
	linkMode200000baseLR4ER4FR4Full
	linkMode200000baseDR4Full // 65
	linkMode200000baseCR4Full
)

func defKeys(data map[uint]string) []uint {
	var keys []uint
	for k := range data {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	return keys
}

var portBits = map[uint]string{
	linkModeTP:        "TP",
	linkModeAUI:       "AUI",
	linkModeBNC:       "BNC",
	linkModeMII:       "MII",
	linkModeFIBRE:     "Fibre",
	linkModeBackplane: "Backplane",
}

var speedBits = map[uint]string{
	linkMode10baseTHalf:             "10base-T/Half",
	linkMode10baseTFull:             "10base-T/Full",
	linkMode100baseTHalf:            "100base-T/Half",
	linkMode100baseTFull:            "100base-T/Full",
	linkMode1000baseTHalf:           "1000base-T/Half",
	linkMode1000baseTFull:           "1000base-T/Full",
	linkMode10000baseTFull:          "10000base-T/Full",
	linkMode1000baseKXFull:          "1000base-KX/Full",
	linkMode10000baseKX4Full:        "10000base-KX4/Full",
	linkMode10000baseKRFull:         "10000base-KR/Full",
	linkMode20000baseMLD2Full:       "20000base-MLD2/Full",
	linkMode20000baseKR2Full:        "20000base-KRR2/Full",
	linkMode40000baseKR4Full:        "40000base-KR4/Full",
	linkMode40000baseCR4Full:        "40000base-CR4/Full",
	linkMode40000baseSR4Full:        "40000base-SR4/Full",
	linkMode40000baseLR4Full:        "40000base-LR4/Full",
	linkMode56000baseKR4Full:        "56000base-KR4/Full",
	linkMode56000baseCR4Full:        "56000base-CR4/Full",
	linkMode56000baseSR4Full:        "56000base-SR4/Full",
	linkMode56000baseLR4Full:        "56000base-LR4/Full",
	linkMode25000baseCRFull:         "25000base-CR/Full",
	linkMode25000baseKRFull:         "25000base-KR/Full",
	linkMode25000baseSRFull:         "25000base-SR/Full",
	linkMode50000baseCR2Full:        "50000base-CR2/Full",
	linkMode50000baseKR2Full:        "50000base-KR2/Full",
	linkMode100000baseKR4Full:       "100000base-KR4/Full",
	linkMode100000baseSR4Full:       "100000base-SR4/Full",
	linkMode100000baseCR4Full:       "100000base-CR4/Full",
	linkMode100000baseLR4ER4Full:    "100000base-LR4ER4/Full",
	linkMode50000baseSR2Full:        "50000base-SR2/Full",
	linkMode1000baseXFull:           "1000base-X/Full",
	linkMode10000baseCRFull:         "10000base-CR/Full",
	linkMode10000baseSRFull:         "10000base-SR/Full",
	linkMode10000baseLRFull:         "10000base-LR/Full",
	linkMode10000baseLRMFull:        "10000base-LRM/Full",
	linkMode10000baseERFull:         "10000base-ER/Full",
	linkMode2500baseTFull:           "2500base-T/Full",
	linkMode5000baseTFull:           "5000base-T/Full",
	linkMode50000baseKRFull:         "50000base-KR/Full",
	linkMode50000baseSRFull:         "50000base-SR/Full",
	linkMode50000baseCRFull:         "50000base-CR/Full",
	linkMode50000baseLRERFRFull:     "50000base-LRERFR/Full",
	linkMode50000baseDRFull:         "50000base-DR/Full",
	linkMode100000baseKR2Full:       "100000base-KR2/Full",
	linkMode100000baseSR2Full:       "100000base-SR2/Full",
	linkMode100000baseCR2Full:       "100000base-CR2/Full",
	linkMode100000baseLR2ER2FR2Full: "100000base-LR2ER2FR2/Full",
	linkMode100000baseDR2Full:       "100000base-DR2/Full",
	linkMode200000baseKR4Full:       "200000base-KR4/Full",
	linkMode200000baseSR4Full:       "200000base-SR4/Full",
	linkMode200000baseLR4ER4FR4Full: "200000base-LR4ER4FR4/Full",
	linkMode200000baseDR4Full:       "200000base-DR4/Full",
	linkMode200000baseCR4Full:       "200000base-CR4/Full",
}

const (
	scharMax         = 127
	linkModeDataSize = 3 * scharMax
)

type linkSet struct {
	cmd                 uint32
	speed               uint32
	duplex              uint8
	port                uint8
	phyAddr             uint8
	autoNeg             uint8
	mdioSupport         uint8
	ethTpMdix           uint8
	ethTpMdixCtrl       uint8
	linkModeMasksNwords int8
	reserved            [8]uint32
	linkModeData        [linkModeDataSize]uint32
}

func (l linkSet) code() uint32 {
	return l.cmd
}

func (l linkSet) dump() string {
	return fmt.Sprintf("%+v", l)
}

func (l linkSet) portName() string {
	return decodePortName(l.port)
}

func (l linkSet) supportedSpeed() []string {
	log.Debug("decoding supported speed ")
	sp := readLinkModeData(0, l.linkModeData)
	return scanSpeedBits(sp)
}

func (l linkSet) advertisingSpeed() []string {
	log.Debug("decoding advertising speed ")
	ad := readLinkModeData(l.linkModeMasksNwords, l.linkModeData)
	return scanSpeedBits(ad)
}

func scanSpeedBits(bitmap uint64) []string {
	var ret []string
	for _, pos := range defKeys(speedBits) {
		if hasBit(bitmap, pos) {
			s := speedBits[pos]
			ret = append(ret, s)
		}
	}
	return ret
}

func newLinkSet() *linkSet {
	lset := new(linkSet)
	lset.cmd = ethtoolGetLinkSettings
	return lset
}

func getLinkSettings(ndev *netdev) (*linkSet, error) {
	log.Debug("getting link settings")

	// getLinkSettings uses a two-way handshake

	// First, determine link mode mask words
	lset := newLinkSet()
	errno, err := ndev.post(lset)
	if err != nil {
		return nil, err
	}

	if errno < 0 {
		return nil, fmt.Errorf("err: %d", errno)
	}

	if lset.cmd != ethtoolGetLinkSettings {
		return nil, fmt.Errorf("cmd:%d returned", lset.cmd)
	}

	if lset.linkModeMasksNwords >= 0 {
		return nil, fmt.Errorf("mask size should be less than 0 (got: %d)", lset.linkModeMasksNwords)
	}

	// Then, send actual request
	nwords := lset.linkModeMasksNwords * -1
	lset = newLinkSet()
	lset.linkModeMasksNwords = nwords
	errno, err = ndev.post(lset)
	if err != nil {
		return nil, err
	}

	if errno < 0 {
		return nil, fmt.Errorf("err: %d", errno)
	}

	if lset.cmd != ethtoolGetLinkSettings {
		return nil, fmt.Errorf("cmd:%d returned", lset.cmd)
	}

	if lset.linkModeMasksNwords <= 0 {
		return nil, fmt.Errorf("mask size should be greater than 0 (got: %d)", lset.linkModeMasksNwords)
	}

	return lset, nil
}

func readLinkModeData(offset int8, data [linkModeDataSize]uint32) uint64 {
	var res uint64
	res = uint64(data[offset+1])<<32 | uint64(data[offset])
	return res
}

func decodePortName(port uint8) string {
	switch port {
	case portTP:
		return "Twisted Pair"
	case portAUI:
		return "AUI"
	case portMII:
		return "MII"
	case portFIBRE:
		return "Fibre"
	case portBNC:
		return "BNC"
	case portDA:
		return "DAC"
	case portNONE:
		return "None"
	case portOTHER:
		return "Other"
	}

	return "Unknown"
}
