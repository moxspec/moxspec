package nvidia

import (
	"encoding/xml"

	"github.com/moxspec/moxspec/loglet"
)

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("nvidia")
}

// NewDecoder creates and initializes a Device as Decoder
func NewDecoder() *Devices {
	ds := new(Devices)
	ds.dict = make(map[string]*GPU)
	return ds
}

func decodeLog(in string) (*SmiLog, error) {
	x := new(SmiLog)
	err := xml.Unmarshal([]byte(in), x)
	return x, err
}
