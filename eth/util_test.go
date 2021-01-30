package eth

import (
	"fmt"
	"testing"
)

func TestIfrnName(t *testing.T) {
	var got [16]byte
	tests := []struct {
		in string
		ex [16]byte
	}{
		{"", [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"mox", [16]byte{'m', 'o', 'x', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"248", [16]byte{'2', '4', '8', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"too-long-input-text-over-sixteen", [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"not■ascii", [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = ifrnName(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestHasBit(t *testing.T) {
	var got bool
	tests := []struct {
		in  uint64
		pos uint
		ex  bool
	}{
		{0, 248, false},
		{248, 1, false},
		{248, 4, true},
		{248, 7, true},
		{248, 8, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = hasBit(tt.in, tt.pos)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t, expect: %t", tt, got, tt.ex)
			}
		})
	}
}
