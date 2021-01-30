package platform

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/actapio/moxspec/util"
)

var releases = []string{
	"/etc/oracle-release",
	"/etc/centos-release",
	"/etc/redhat-release",
}

// GetDistroName returns distro pretty name
func GetDistroName() string {
	var name string

	// http://0pointer.de/blog/projects/os-release
	log.Debug("parsing /etc/os-release")
	fd, err := os.Open("/etc/os-release")
	if err == nil {
		name = parseOSRelease(fd)
		if name != "" {
			return name
		}
	}

	rel, err := util.ScanPathList(releases)
	if err != nil {
		return ""
	}

	log.Debugf("found %s", rel)
	name, _ = util.LoadString(rel)
	if name != "" {
		return name
	}

	return ""
}

func parseOSRelease(fd io.Reader) string {
	bs, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Debug(err)
		return ""
	}

	var name, ver string
	for _, line := range strings.Split(string(bs), "\n") {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "#") {
			continue
		}
		if strings.HasPrefix(l, "=") {
			continue
		}

		flds := strings.Split(l, "=")
		if len(flds) != 2 {
			continue
		}

		key := flds[0]
		val := strings.TrimSuffix(strings.TrimPrefix(flds[1], "\""), "\"")

		log.Debugf("key: %s / val: %s", key, val)

		switch key {
		case "NAME":
			name = val
		case "VERSION":
			ver = val
		case "PRETTY_NAME":
			return val
		default:
			continue
		}
	}

	if name != "" && ver != "" {
		return fmt.Sprintf("%s %s", name, ver)
	}

	return ""
}
