package pci

import (
	"fmt"
	"testing"
)

func TestParseLocator(t *testing.T) {
	var dom, bus, dev, fun uint32
	var err error
	tests := []struct {
		in  string
		dom uint32
		bus uint32
		dev uint32
		fun uint32
	}{
		{"0000:d7:12.0", 0, 0xD7, 0x12, 0},
		{"mox:mox:mox.mox", 0, 0, 0, 0},
		{"mox:0:mox.0", 0, 0, 0, 0},
		{"0:mox:0.mox", 0, 0, 0, 0},
		{"10000:d7:12.0", 0x10000, 0xD7, 0x12, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			dom, bus, dev, fun, err = ParseLocater(tt.in)
			if dom != tt.dom || bus != tt.bus || dev != tt.dev || fun != tt.fun {
				t.Errorf("test: %+v, got: dom=%d bus=%d dev=%d fun=%d, expect: dom=%d bus=%d dev=%d fun=%d (err:%s)", tt,
					dom, bus, dev, fun, tt.dom, tt.bus, tt.dev, tt.fun, err)
			}
		})
	}
}
