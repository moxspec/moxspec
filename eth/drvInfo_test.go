package eth

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseFirmwareVersion(t *testing.T) {
	tests := []struct {
		in [32]byte
		ex string
	}{
		{[32]byte{}, ""},
		{[32]byte{0x41, 0x42, 0x43}, "ABC"},
		{[32]byte{0x41, 0x42, 0x43, 0x00, 0x41}, "ABC"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got := parseFirmwareVersion(tt.in)

			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %v, expect: %v", tt, got, tt.ex)
			}
		})
	}
}
