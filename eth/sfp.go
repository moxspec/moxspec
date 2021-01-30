package eth

import (
	"fmt"
	"strings"
)

const (
	conUnknown     = 0x00
	conSC          = 0x01
	conFC1         = 0x02
	conFC2         = 0x03
	conBNCTNC      = 0x04
	conFCCOAX      = 0x05
	conFiberjack   = 0x06
	conLC          = 0x07
	conMTRJ        = 0x08
	conMU          = 0x09
	conSG          = 0x0A
	conOptDAC      = 0x0B
	conMPO         = 0x0C
	conMPO2        = 0x0D
	conHSDCII      = 0x20
	conCopperDAC   = 0x21
	conRJ45        = 0x22
	conNOSEPARABLE = 0x23
	conMXC2x16     = 0x24
)

const (
	idUnknown    = 0x00
	idGBIC       = 0x01
	idSoldered   = 0x02
	idSFP        = 0x03
	id300pinXBI  = 0x04
	idXENPAK     = 0x05
	idXFP        = 0x06
	idXFF        = 0x07
	idXFPE       = 0x08
	idXPAK       = 0x09
	idX2         = 0x0A
	idDWDMSFP    = 0x0B
	idQSFP       = 0x0C
	idQSFPP      = 0x0D
	idCXP        = 0x0E
	idHD4X       = 0x0F
	idHD8X       = 0x10
	idQSFP28     = 0x11
	idCXP2       = 0x12
	idCDFP       = 0x13
	idHD4XFanout = 0x14
	idHD8XFanout = 0x15
	idCDFPS3     = 0x16
	idMQSFP      = 0x17
)

func decodeSFF8472(data [maxRomSize]byte) (*Module, error) {
	return decodeSFF8079(data)
}

func decodeSFF8079(data [maxRomSize]byte) (*Module, error) {
	ff := decodeSFF8024ID(data[0])
	cn := decodeSFF8024Con(data[2])
	cl := scanSFPCableLen(data[15:20])
	vn := strings.TrimSpace(fmt.Sprintf("%s", data[20:36]))
	pn := strings.TrimSpace(fmt.Sprintf("%s", data[40:56]))
	sn := strings.TrimSpace(fmt.Sprintf("%s", data[68:84]))

	log.Debugf("formfactor: %s", ff)
	log.Debugf("connector: %s", cn)
	log.Debugf("cable length: %dm", cl)
	log.Debugf("vendor name: %s", vn)
	log.Debugf("product name: %s", pn)
	log.Debugf("serial number: %s", sn)

	md := new(Module)
	md.FormFactor = ff
	md.Connector = cn
	md.CableLength = cl
	md.VendorName = vn
	md.ProductName = pn
	md.SerialNumber = sn

	return md, nil
}

func scanSFPCableLen(data []byte) byte {
	var l byte
	for _, d := range data {
		if d > l {
			l = d
		}
	}
	return l
}

func decodeSFF8024ID(in byte) string {
	switch in {
	case idUnknown:
		return "unknown"
	case idGBIC:
		return "GBIC"
	case idSoldered:
		return "Soldered Module"
	case idSFP:
		return "SFP"
	case id300pinXBI:
		return "300 pin XBI"
	case idXENPAK:
		return "XENPAK"
	case idXFP:
		return "XFP"
	case idXFF:
		return "XFF"
	case idXFPE:
		return "XFP-E"
	case idXPAK:
		return "XPAK"
	case idX2:
		return "X2"
	case idDWDMSFP:
		return "DWDM-SFP"
	case idQSFP:
		return "QSFP"
	case idQSFPP:
		return "QSFP+"
	case idCXP:
		return "CXP"
	case idHD4X:
		return "MiniSAS-4x"
	case idHD8X:
		return "MiniSAS-8x"
	case idQSFP28:
		return "QSFP28"
	case idCXP2:
		return "CXP2"
	case idCDFP:
		return "CDFP"
	case idHD4XFanout:
		return "MiniSAS-4x Fanout"
	case idHD8XFanout:
		return "MiniSAS-8x Fanout"
	case idCDFPS3:
		return "CDFP Style 3"
	case idMQSFP:
		return "MicroQSFP"
	}
	return "unknown"
}

func decodeSFF8024Con(in byte) string {
	switch in {
	case conUnknown:
		return "unknown"
	case conSC:
		return "SC"
	case conFC1:
		return "FC Style 1 copper"
	case conFC2:
		return "FC Style 2 copper"
	case conBNCTNC:
		return "BNC-TNC"
	case conFCCOAX:
		return "FC COAX"
	case conFiberjack:
		return "Fibrejack"
	case conLC:
		return "LC"
	case conMTRJ:
		return "MTRJ"
	case conMU:
		return "MU"
	case conSG:
		return "SG"
	case conOptDAC:
		return "Optical DAC"
	case conMPO:
		return "MPO"
	case conMPO2:
		return "MPO 2x16"
	case conHSDCII:
		return "HSSDC II"
	case conCopperDAC:
		return "Copper DAC"
	case conRJ45:
		return "RJ45"
	case conNOSEPARABLE:
		return "No separable"
	case conMXC2x16:
		return "MXC 2x16"
	}
	return "unknown"
}
