package msr

import (
	"fmt"
	"os"

	"github.com/moxspec/moxspec/loglet"
	"github.com/moxspec/moxspec/util"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("msr")
}

type vendor int

// Vendors
const (
	INTEL vendor = iota
	AMD
)

type msrReader func(int64) uint64

// Register represents a msr
type Register struct {
	ID     uint16
	Vendor vendor
	Temp   int16
}

// NewDecoder creates and initializes a Register as Decoder
func NewDecoder(id uint16, v vendor) *Register {
	r := new(Register)
	r.ID = id
	r.Vendor = v
	return r
}

// Decode makes Register satisfy the mox.Decoder interface
func (r *Register) Decode() error {
	path := fmt.Sprintf("/dev/cpu/%d/msr", r.ID)
	fd, err := os.OpenFile(path, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return err
	}
	defer fd.Close()

	rdr := func(ptr int64) uint64 {
		d := make([]byte, 8, 8)
		fd.ReadAt(d, ptr)
		ret := util.BytesToUint64(d)
		log.Debugf("msr addr: 0x%x => got: %v (%d)", ptr, d, ret)
		return ret
	}

	r.Temp = readTemp(rdr, r.Vendor)

	return nil
}

func readTemp(rdr msrReader, v vendor) int16 {
	var res int16
	if v == INTEL {
		do := int16((rdr(0x19C) >> 16) & 0x7F)
		tt := int16((rdr(0x1A2) >> 16) & 0xFF)
		res = tt - do
	}
	return res
}
