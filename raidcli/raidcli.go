package raidcli

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("raidcli")
}

// Level represents RAID level
type Level string

// These are raid levels
const (
	Unknown Level = "unknown"
	RAID0   Level = "RAID 0"
	RAID1   Level = "RAID 1"
	RAID5   Level = "RAID 5"
	RAID6   Level = "RAID 6"
	RAID01  Level = "RAID 0+1"
	RAID10  Level = "RAID 1+0"
)

// Base represents base number
type Base float64

// These are base numbers
const (
	Decimal Base = 1000.0
	Binary  Base = 1024.0
)

// Run runs given external command
func Run(clipath string, arg ...string) (string, error) {
	log.Debugf("running %s %s", clipath, strings.Join(arg, " "))
	return util.Exec(clipath, arg...)
}

// SplitKeyVal returns key-value pair separated by given delimier
func SplitKeyVal(line, delim string) (key string, val string, err error) {
	log.Debug(line)
	l := strings.TrimSpace(line)
	if l == "" {
		err = fmt.Errorf("empty line")
		return
	}

	if !strings.Contains(l, delim) {
		err = fmt.Errorf("line has no '%s'", delim)
		return
	}

	flds := strings.SplitN(l, delim, 2)
	if len(flds) != 2 {
		err = fmt.Errorf("invalid format")
		return
	}

	key = strings.TrimSpace(flds[0])
	val = strings.TrimSpace(flds[1])

	return
}

// ParseSize returns unsigned integer value from given string with unit and base
func ParseSize(val string, base Base) uint64 {
	flds := strings.Fields(val)
	if len(flds) == 0 {
		return 0
	}
	size, err := strconv.ParseFloat(flds[0], 64)
	if err != nil {
		return 0
	}

	return uint64(size * GetMultiplier(val, base))
}

// GetMultiplier returns multiplier from given unit and base
func GetMultiplier(val string, base Base) float64 {
	b := float64(base)

	if strings.Contains(val, "KB") {
		return math.Pow(b, 1)
	}
	if strings.Contains(val, "MB") {
		return math.Pow(b, 2)
	}
	if strings.Contains(val, "GB") {
		return math.Pow(b, 3)
	}
	if strings.Contains(val, "TB") {
		return math.Pow(b, 4)
	}
	if strings.Contains(val, "PB") {
		return math.Pow(b, 5)
	}
	if strings.Contains(val, "EB") {
		return math.Pow(b, 6)
	}
	return 0
}

// ShapeSpacedString returns shaped string which is remoeved multiple spaces
// e.g:
//   Before: 7JJ7W7GCHUH721010ALE600                         T281
//   After:  7JJ7W7GCHUH721010ALE600 T281
func ShapeSpacedString(val string) string {
	flds := strings.Fields(val)
	return strings.Join(flds, " ")
}
