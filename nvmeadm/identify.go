package nvmeadm

import (
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"

	"github.com/moxspec/moxspec/util"
)

const (
	lbafMax = 15 // up to LBAF15
)

func identifyController(fd *os.File, d *Device) error {
	data := make([]byte, 4096, 4096)

	cmd := nvmeAdminCmd{
		opcode:  0x06,
		nsid:    0, // 5.15.2 Identify Namespace data structure (CNS 00h)
		addr:    *(*uint64)(unsafe.Pointer(&data)),
		dataLen: 4096,
		cdw10:   0x01,
	}

	log.Debugf("nvme admin cmd: %+v", cmd)
	err := post(fd, &cmd)
	if err != nil {
		return err
	}
	log.Debugf("result: %+v", data)

	d.SerialNumber = util.SanitizeString(fmt.Sprintf("%s", data[4:24]))
	d.ModelNumber = util.SanitizeString(fmt.Sprintf("%s", data[24:64]))
	d.FirmwareRevision = util.SanitizeString(fmt.Sprintf("%s", data[64:72]))
	d.WarnTemp = int16(binary.LittleEndian.Uint16(data[266:268]))
	if d.WarnTemp > 0 {
		d.WarnTemp = d.WarnTemp - kelvin
	}
	d.CritTemp = int16(binary.LittleEndian.Uint16(data[268:270]))
	if d.CritTemp > 0 {
		d.CritTemp = d.CritTemp - kelvin
	}
	d.Size = binary.LittleEndian.Uint64(data[280:296])
	d.MaxNamespaces = binary.LittleEndian.Uint32(data[516:520])

	var i uint32
	for i = 1; i <= d.MaxNamespaces; i++ {
		cap, err := identifyNamespaceSize(fd, i)
		if err != nil {
			log.Warn(err.Error())
		}
		d.NamespaceSizes = append(d.NamespaceSizes, cap)
	}

	return nil
}

// 5.15.2 Identify Namespace data structure (CNS 00h)
func identifyNamespaceSize(fd *os.File, id uint32) (uint64, error) {
	log.Debugf("identifying namespace %d", id)

	data := make([]byte, 4096, 4096)
	cmd := nvmeAdminCmd{
		opcode:  0x06,
		nsid:    id,
		addr:    *(*uint64)(unsafe.Pointer(&data)),
		dataLen: 4096,
		cdw10:   0x00,
	}

	log.Debugf("nvme admin cmd: %+v", cmd)
	err := post(fd, &cmd)
	if err != nil {
		return 0, err
	}
	log.Debugf("result: %+v", data)

	return parseNamespaceSize(data)
}

func parseNamespaceSize(data []byte) (uint64, error) {
	nsze := binary.LittleEndian.Uint64(data[0:8])
	ncap := binary.LittleEndian.Uint64(data[8:16])
	log.Debugf("nsze=%d ncap=%d", nsze, ncap)

	lba, err := parseLBASize(data)
	if err != nil {
		return 0, err
	}

	nssize := nsze * lba
	log.Debugf("capacity: %d bytes", nssize)
	return nssize, nil
}

func parseLBASize(data []byte) (uint64, error) {
	lbafs := parseLBAFList(data)
	if len(lbafs) == 0 {
		return 0, fmt.Errorf("could not parse lbaf list")
	}

	lid := data[26] & 0xF
	if lid > byte(len(lbafs))-1 {
		return 0, fmt.Errorf("lbaf id:%d is out of range for lbaf list", lid)
	}

	lsize := lbafs[lid]
	log.Debugf("lbaf id: %d, size: %d", lid, lsize)
	return lsize, nil
}

func parseLBAFList(data []byte) []uint64 {
	log.Debug("parsing lbaf list")
	idx := 128 // LBA Format definition starts from 128

	var res []uint64
	for fid := 0; fid <= lbafMax; fid++ {
		// Figure 110: Identify – LBA Format Data Structure, NVM Command Set Specific
		//
		// idx   : Metadata Size (MS)
		// idx+1 : Metadata Size (MS)
		// idx+2 : LBA Data Size (LBADS)
		// idx+3 : Relative Performance (RP)
		// =====
		// total: 4 bytes

		// LBADS
		//   The value is reported in terms of a power of two (2^n).
		//   A value smaller than 9 (i.e., 512 bytes) is not supported.
		size := uint64(1 << data[idx+2])
		if size == 1 {
			size = 0
		}
		res = append(res, size)

		log.Debugf("LBAF%d size=%d", fid, size)

		idx = idx + 4
		if idx > len(data)-1 {
			break
		}
	}

	return res
}
