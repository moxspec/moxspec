package eth

import (
	"fmt"
	"testing"
)

func TestIsValidROMSize(t *testing.T) {
	var got bool
	tests := []struct {
		in uint32
		ex bool
	}{
		{0, false},
		{maxRomSize, true},
		{maxRomSize + 1, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = isValidROMSize(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t, expect: %t", tt, got, tt.ex)
			}
		})
	}
}
