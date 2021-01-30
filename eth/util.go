package eth

import (
	"fmt"
	"syscall"
	"unicode"
	"unsafe"
)

// ethtool.h
const (
	ethtoolIO              = 0x8946
	ethtoolGetLink         = 0x0a
	ethtoolGetLinkSettings = 0x4c
	ethtoolGetModuleInfo   = 0x42
	ethtoolGetEeprom       = 0x43
	ethtoolGetDriver       = 0x03
)

type request interface {
	code() uint32
	dump() string
}

type netdev struct {
	name [16]byte
	fd   int
}

func (n *netdev) close() error {
	return syscall.Close(n.fd)
}

func (n netdev) post(in request) (int, error) {
	return post(n.fd, n.name, in)
}

func newNetdev(name string) (*netdev, error) {
	n := new(netdev)
	n.name = ifrnName(name)

	fd, err := openSocket()
	if err != nil {
		return nil, err
	}
	n.fd = fd

	return n, nil
}

func post(fd int, name [16]byte, in request) (int, error) {
	var errno int

	log.Debugf("request code: %d", in.code())

	switch in.(type) {
	case *linkSet:
		e := struct {
			ifrnName [16]byte
			ifruData *linkSet
		}{
			ifrnName: name,
			ifruData: in.(*linkSet),
		}
		errno = ioctl(fd, uintptr(unsafe.Pointer(&e)))
	case *linkStat:
		e := struct {
			ifrnName [16]byte
			ifruData *linkStat
		}{
			ifrnName: name,
			ifruData: in.(*linkStat),
		}
		errno = ioctl(fd, uintptr(unsafe.Pointer(&e)))
	case *moduleInfo:
		e := struct {
			ifrnName [16]byte
			ifruData *moduleInfo
		}{
			ifrnName: name,
			ifruData: in.(*moduleInfo),
		}
		errno = ioctl(fd, uintptr(unsafe.Pointer(&e)))
	case *eeprom:
		e := struct {
			ifrnName [16]byte
			ifruData *eeprom
		}{
			ifrnName: name,
			ifruData: in.(*eeprom),
		}
		errno = ioctl(fd, uintptr(unsafe.Pointer(&e)))
	case *drvInfo:
		e := struct {
			ifrnName [16]byte
			ifruData *drvInfo
		}{
			ifrnName: name,
			ifruData: in.(*drvInfo),
		}
		errno = ioctl(fd, uintptr(unsafe.Pointer(&e)))
	default:
		return errno, fmt.Errorf("invalid ioctl input given")
	}

	log.Debugf("responce: %s", in.dump())
	log.Debugf("errno: %d", errno)

	return errno, nil
}

func ioctl(fd int, ptr uintptr) int {
	r1, r2, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(ethtoolIO),
		ptr,
	)
	log.Debugf("r1: %d, r2: %d, errno: %d", r1, r2, errno)
	return int(errno)
}

func ifrnName(name string) [16]byte {
	log.Debugf("generating ifrnName from '%s'", name)

	var buf, ret [16]byte
	if name == "" {
		return ret
	}

	if len(([]byte)(name)) > 16 {
		return ret
	}

	for i, c := range ([]byte)(name) {
		if i >= 16 {
			break
		}
		log.Debugf("%d: %c", i, c)

		if c > unicode.MaxASCII {
			log.Warnf("%c is greater than unicode.MaxASCII(%d)", c, unicode.MaxASCII)
			return ret
		}

		buf[i] = c
	}

	ret = buf
	return ret
}

func openSocket() (int, error) {
	log.Debug("opening socket")
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Debugf("failed: %s", err)
		return fd, err
	}
	log.Debugf("fd = %d", fd)
	return fd, nil
}

func hasBit(data uint64, pos uint) bool {
	return (((data >> pos) & 0x1) == 1)
}
