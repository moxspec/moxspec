package pci

import (
	"strconv"
	"strings"
)

// Class code
const (
	MassStorageController   = 0x01
	NetworkController       = 0x02
	DisplayController       = 0x03 // for GPU
	CommunicationController = 0x07 // for FPGA
	ProcessingAccelerator   = 0x12 // for FPGA
)

func parseHexStr(hexStr string) (uint64, error) {
	cleaned := strings.Replace(strings.TrimSpace(hexStr), "0x", "", -1)
	i64, err := strconv.ParseUint(cleaned, 16, 64)
	return i64, err
}
