package platform

import (
	"fmt"
	"syscall"

	"github.com/moxspec/moxspec/util"
)

// Uname returns system name, release, nodename and machine
func Uname() (sysname, release, nodename, machine string) {
	uts := syscall.Utsname{}
	syscall.Uname(&uts)

	sysname = parseUtsField(uts.Sysname)
	release = parseUtsField(uts.Release)
	nodename = parseUtsField(uts.Nodename)
	machine = parseUtsField(uts.Machine)

	return
}

func parseUtsField(is [65]int8) string {
	s := make([]byte, len(is))
	for i, j := range is {
		if j == 0 {
			break
		}
		s[i] = byte(j)
	}
	return util.SanitizeString(fmt.Sprintf("%s", s))
}
