package cpuid

import (
	"fmt"
	"testing"
)

func TestUnit32sBackward(t *testing.T) {
	var got []byte
	tests := []struct {
		in []uint32
		ex string
	}{
		{[]uint32{0x756e6547, 0x49656e69, 0x6c65746e}, "GenuineIntel"},
		{[]uint32{0x68747541, 0x69746e65, 0x444d4163}, "AuthenticAMD"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = readUint32sBackward(tt.in...)
			if string(got) != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestIsValidEaxIn(t *testing.T) {
	var got bool
	tests := []struct {
		eax uint32
		std uint32
		ext uint32
		ex  bool
	}{
		{0x00000001, 0x0000000F, 0x8000000F, true},
		{0x0000000F, 0x0000000F, 0x8000000F, true},
		{0x000000FF, 0x0000000F, 0x8000000F, false},
		{0x80000001, 0x0000000F, 0x8000000F, true},
		{0x8000000F, 0x0000000F, 0x8000000F, true},
		{0x800000FF, 0x0000000F, 0x8000000F, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, _ = isValidEaxIn(tt.eax, tt.std, tt.ext)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t, expect: %t", tt, got, tt.ex)
			}
		})
	}
}
