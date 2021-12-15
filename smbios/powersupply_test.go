package smbios

import (
	"fmt"
	"testing"
)

func TestParsePowerSupplyStats(t *testing.T) {
	var pl, pr, ho bool
	tests := []struct {
		in uint16
		pl bool
		pr bool
		ho bool
	}{
		{0x00, true, false, false},
		{0x01, true, false, true},
		{0x02, true, true, false},
		{0x03, true, true, true},
		{0x04, false, false, false},
		{0x05, false, false, true},
		{0x06, false, true, false},
		{0x07, false, true, true},
		{0x248, true, false, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			pl, pr, ho = parsePowerSupplyStats(tt.in)
			if pl != tt.pl || pr != tt.pr || ho != tt.ho {
				t.Errorf("test: %+v, got: %t %t %t, expect: %t %t %t", tt, pl, pr, ho, tt.pl, tt.pr, tt.ho)
			}
		})
	}

}
