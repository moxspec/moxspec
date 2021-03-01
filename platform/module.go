package platform

import (
	"strings"

	"github.com/moxspec/moxspec/util"
)

// IsLoadedModule returns if a module is loaded
func IsLoadedModule(name string) bool {
	list, err := util.LoadString("/proc/modules")
	if err != nil {
		return false
	}

	return scanModules(list, name)
}

func scanModules(list, name string) bool {
	if list == "" || name == "" {
		return false
	}

	for _, l := range strings.Split(list, "\n") {
		flds := strings.Fields(l)
		if len(flds) == 0 {
			continue
		}

		if flds[0] == name {
			return true
		}
	}

	return false
}
