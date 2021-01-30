package msr

import (
	"fmt"
	"testing"
)

func TestReadTemp(t *testing.T) {
	intelRdr := func(p int64) uint64 {
		switch p {
		case 0x19C:
			return 2286028800
		case 0x1A2:
			return 6425088
		}
		return 0
	}

	var got int16
	tests := []struct {
		rdr msrReader
		ven vendor
		ex  int16
	}{
		{intelRdr, INTEL, 32},
		{intelRdr, AMD, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = readTemp(tt.rdr, tt.ven)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}

}
