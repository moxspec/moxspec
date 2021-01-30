package nvmeadm

import (
	"encoding/binary"
	"os"
	"unsafe"
)

// 5.14.1.2 SMART / Health Information (Log Identifier 02h)
// https://nvmexpress.org/wp-content/uploads/NVM-Express-1_3c-2018.05.24-Ratified.pdf
func readSmartLog(fd *os.File, d *Device) error {
	data := make([]byte, nvmeSmartLen, nvmeSmartLen)

	cmd := nvmeAdminCmd{
		opcode:  0x02,
		nsid:    0xFFFFFFFF,
		addr:    *(*uint64)(unsafe.Pointer(&data)),
		dataLen: nvmeSmartLen,

		// bit 31:28 reserved = 0h
		// bit 27:16 NUMBER OF DWORDS = 07Fh(512B)
		// bit 15:08 RESERVED = 00h
		// bit 7:0 LOG PAGE IDENTIFIER
		cdw10: 0x007F0002,
	}

	log.Debugf("nvme admin cmd: %+v", cmd)
	err := post(fd, &cmd)
	if err != nil {
		return err
	}
	log.Debugf("result: %+v", data)

	d.CurTemp = int16(binary.LittleEndian.Uint16(data[1:3]))
	if d.CurTemp > 0 {
		d.CurTemp = d.CurTemp - kelvin
	}
	d.SpareSpace = data[3]
	d.Used = data[5]
	d.UnitsRead = binary.LittleEndian.Uint64(data[32:48])
	d.UnitsWritten = binary.LittleEndian.Uint64(data[48:64])
	d.ByteRead = d.UnitsRead * 512 * 1000
	d.ByteWritten = d.UnitsWritten * 512 * 1000
	d.PowerCycleCount = binary.LittleEndian.Uint64(data[112:128])
	d.PowerOnHours = binary.LittleEndian.Uint64(data[128:144])
	d.UnsafeShutdownCount = binary.LittleEndian.Uint64(data[144:160])
	d.UnrecoveredError = binary.LittleEndian.Uint64(data[160:176]) // Media and Data Integrity Errors
	d.CritWarnings = parseCritWarnings(data[0])

	return nil
}

func parseCritWarnings(b byte) []string {
	ws := []string{
		"the available spare capacity has fallen below the threshold.",
		"a temperature is above an over temperature threshold or below an under temperature threshold.",
		"the NVM subsystem reliability has been degraded due to significant media related errors or any internal error that degrades NVM subsystem reliability.",
		"the media has been placed in read only mode.",
		"the volatile memory backup device has failed. This field is only valid if the controller has a volatile memory backup solution.",
	}

	var res []string
	var i byte
	for i = 0; i < byte(len(ws)); i++ {
		if b&1<<i == 0 {
			continue
		}
		res = append(res, ws[i])
	}
	return res
}
