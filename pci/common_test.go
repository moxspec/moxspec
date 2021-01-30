package pci

import (
	"fmt"
	"testing"
)

func TestParseHexStr(t *testing.T) {
	var got uint64
	var err error
	tests := []struct {
		in string
		ex uint64
	}{
		{"0", 0},
		{"0xFF", 255},
		{"FF", 255},
		{"mox", 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err = parseHexStr(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d (err:%s)", tt, got, tt.ex, err)
			}
		})
	}
}
