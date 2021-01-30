package eth

import (
	"fmt"
	"testing"
)

func TestScanSFPCableLen(t *testing.T) {
	var got byte
	tests := []struct {
		in []byte
		ex byte
	}{
		{[]byte{2, 4, 8}, 8},
		{[]byte{8, 4, 2}, 8},
		{[]byte{4, 8, 2}, 8},
		{[]byte{4}, 4},
		{[]byte{0}, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = scanSFPCableLen(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
