package eth

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDefKeys(t *testing.T) {
	var got []uint
	tests := []struct {
		in map[uint]string
		ex []uint
	}{
		{map[uint]string{2: "ni", 4: "shi", 8: "ya"}, []uint{2, 4, 8}},
		{map[uint]string{4: "shi", 2: "ni", 8: "ya"}, []uint{2, 4, 8}},
		{map[uint]string{8: "ya", 4: "shi", 2: "ni"}, []uint{2, 4, 8}},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = defKeys(tt.in)
			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %v, expect: %v", tt, got, tt.ex)
			}
		})
	}
}

func TestScanSpeedBits(t *testing.T) {
	var got []string
	tests := []struct {
		in uint64
		ex []string
	}{
		{64, nil},
		{248, []string{"100base-T/Full", "1000base-T/Half", "1000base-T/Full"}},
		{2148212800, []string{"1000base-KX/Full", "10000base-KR/Full", "25000base-CR/Full"}},
		{9223372036854775808, []string{"200000base-SR4/Full"}},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = scanSpeedBits(tt.in)
			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %v, expect: %v", tt, got, tt.ex)
			}
		})
	}
}

func TestDecodePortName(t *testing.T) {
	var got string
	tests := []struct {
		in uint8
		ex string
	}{
		{0x00, "Twisted Pair"},
		{0x01, "AUI"},
		{0x02, "MII"},
		{0x03, "Fibre"},
		{0x04, "BNC"},
		{0x05, "DAC"},
		{0x08, "Unknown"},
		{0xef, "None"},
		{0xff, "Other"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = decodePortName(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}
