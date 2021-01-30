package megaraid

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/spc"
)

var log *loglet.Logger

const (
	iovSize         = 512
	iocFirmware     = 0xC1944D01
	iocPacketSize   = 404
	cmdStatusOffset = 22
	smartAttrSize   = 12
	ioctlnodeFile   = "/dev/megaraid_sas_ioctl_node"
)

func init() {
	log = loglet.NewLogger("spc-megaraid")
}

// NewDecoder creates and initializes a Device as Decoder
func NewDecoder(ctl, dev int, diskType spc.DiskType) *spc.Device {
	post := func(fd *os.File, cdb []byte) ([]byte, error) {
		return post(fd, ctl, dev, cdb)
	}

	return spc.NewDevice(post, ioctlnodeFile, diskType)
}

type iovBuf = [iovSize]byte
type pktBuf = [iocPacketSize]byte

func cmdStatus(p *pktBuf) int {
	return int(p[cmdStatusOffset])
}

type iovec struct {
	iovBase uint64
	iovSize uint64
}

type iocFrame struct {
	cmd               uint8
	senseLen          uint8
	cmdStatus         uint8
	scsiStatus        uint8
	targetID          uint8
	lun               uint8
	cdbLen            uint8
	sgeCount          uint8
	context           uint32
	padding1          uint32
	flags             uint16
	timeout           uint16
	dataXferLen       uint32
	sensuBufPhyAddrLo uint32
	sensuBufPhyAddrHi uint32
	cdb               [16]uint8
	sgePhysAddr       uint64
	sgeLength         uint32
	padding2          [68]uint8
}

type iocPacket struct {
	hostNo      uint16
	padding     uint16
	sglOffset   uint32
	sgeCount    uint32
	senseOffset uint32
	senseLen    uint32
	frame       iocFrame
	sgl         [16]iovec
}

func newIocPkt(ctl, dev int, cdb []byte) (*pktBuf, *iovBuf, error) {
	if len(cdb) > 16 {
		return nil, nil, fmt.Errorf("cdb is too long")
	}

	if len(cdb) == 0 {
		return nil, nil, fmt.Errorf("cdb is empty")
	}

	ibuf := iovBuf{}
	packet := iocPacket{
		hostNo:    uint16(ctl),
		sglOffset: 48,
		sgeCount:  1,
		frame: iocFrame{
			cmd:         0x04,
			cmdStatus:   0xFF,
			targetID:    uint8(dev),
			cdbLen:      uint8(len(cdb)),
			sgeCount:    1,
			flags:       0x0010,
			dataXferLen: uint32(len(ibuf)),
			sgePhysAddr: uint64(uintptr(unsafe.Pointer(&ibuf))),
			sgeLength:   uint32(len(ibuf)),
		},
		sgl: [16]iovec{
			{
				iovBase: uint64(uintptr(unsafe.Pointer(&ibuf))),
				iovSize: uint64(len(ibuf)),
			},
		},
	}

	copy(packet.frame.cdb[:], cdb)

	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, packet)
	if err != nil {
		return nil, nil, err
	}

	pbuf := pktBuf{}
	copy(pbuf[:], b.Bytes())

	return &pbuf, &ibuf, nil
}

func post(fd *os.File, ctl, dev int, cdb []byte) ([]byte, error) {
	pbuf, ibuf, err := newIocPkt(ctl, dev, cdb)
	if err != nil {
		return nil, err
	}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd.Fd()),
		uintptr(iocFirmware),
		uintptr(unsafe.Pointer(pbuf)))

	if errno != 0 {
		return nil, fmt.Errorf("ioctl failure (errno = %d)", errno)
	}

	cmdStatus := cmdStatus(pbuf)
	if cmdStatus != 0 {
		return nil, fmt.Errorf("ioctl failure (cmd status = %d)", cmdStatus)
	}

	return ibuf[:], nil
}
