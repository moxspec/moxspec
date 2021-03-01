package cpuid

import "github.com/moxspec/moxspec/loglet"

type parser func(*Processor) error

var log *loglet.Logger

func init() {
	log = loglet.NewLogger("cpuid")
}

// NewDecoder creates and initializes a Processor as Decoder
func NewDecoder() *Processor {
	return new(Processor)
}

// Decode makes Processor satisfy the mox.Decoder interface
func (cpu *Processor) Decode() error {
	var err error

	var parsers = []parser{
		parseVendorID,
		parseFlags,
		parseTLB,
		parseCache,
		parseFrequency,
		parseTopology,
		parseBrandString,
	}

	for _, p := range parsers {
		err = p(cpu)
		if err != nil {
			return err
		}
	}

	return nil
}
