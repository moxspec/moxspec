package acs

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/spc"
)

const (
	bufSize        = 512
	sbufSize       = 32
	sgDxferFromDev = -3 // from target to initiator
	sgIO           = 0x2285
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("spc-pch")
}

// NewDecoder creates and initializes a Device as Decoder
func NewDecoder(path string) *spc.Device {
	return spc.NewDevice(post, path, spc.SATADisk)
}

// http://www.tldp.org/HOWTO/SCSI-Generic-HOWTO/sg_io_hdr_t.html
type sgIoHdr struct {
	interfaceID    int32
	dxferDirection int32
	cmdLen         byte
	mxSbLen        byte
	iovecCount     uint16
	dxferLen       uint32
	dxferPtr       uint64
	cmdPtr         uint64
	sbPtr          uint64
	timeout        uint32
	flags          uint32
	packID         int32
	usrPtr         uint64
	status         byte
	maskedStatus   byte
	msgStatus      byte
	sbLenWr        byte
	hostStatus     uint16
	driverStatus   uint16
	resid          int32
	duration       uint32
	info           uint32
}

func post(fd *os.File, cdb []byte) ([]byte, error) {
	buf := make([]byte, bufSize, bufSize)
	sbuf := make([]byte, sbufSize, sbufSize)

	hdr := sgIoHdr{
		interfaceID:    int32('S'), // fixed
		cmdLen:         byte(len(cdb)),
		mxSbLen:        sbufSize,
		dxferDirection: sgDxferFromDev,
		dxferLen:       bufSize,
		dxferPtr:       *(*uint64)(unsafe.Pointer(&buf)),
		cmdPtr:         *(*uint64)(unsafe.Pointer(&cdb)),
		sbPtr:          *(*uint64)(unsafe.Pointer(&sbuf)),
		timeout:        500,
	}

	log.Debugf("sghdr: %+v", hdr)

	r1, r2, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd.Fd()),
		uintptr(sgIO),
		uintptr(unsafe.Pointer(&hdr)),
	)

	log.Debugf("r1: %d, r2: %d, errno: %d", r1, r2, errno)
	if errno != 0 {
		return nil, fmt.Errorf("ioctl failed: code %d", errno)
	}

	log.Debugf("result (buffer): %+v", buf)
	log.Debugf("result (sense buffer): %+v", sbuf)

	return buf, nil
}
