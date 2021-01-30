package nvmeadm

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	kelvin        = 273
	nvmeSmartLen  = 512
	ioctlAdminCmd = 0xC0484E41
)

func post(fd *os.File, cmdptr *nvmeAdminCmd) error {
	r1, r2, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd.Fd()),
		uintptr(ioctlAdminCmd),
		uintptr(unsafe.Pointer(cmdptr)),
	)

	log.Debugf("r1: %d, r2: %d, errno: %d", r1, r2, errno)
	if r1 != 0 || errno != 0 {
		return fmt.Errorf("ioctl failed: r1=0x%X, errno=%d", r1, errno)
	}
	return nil
}
