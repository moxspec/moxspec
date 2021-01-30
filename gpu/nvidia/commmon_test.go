package nvidia

import (
	"fmt"
	"testing"
)

func TestConvUtilString(t *testing.T) {
	var got float64
	tests := []struct {
		in  string
		suf string
		ex  float64
	}{
		{"N/A", "%", 0},
		{"mox", "248", 0},
		{"248", "", 248},
		{"", "248", 0},
		{"", "", 0},
		{"248 mox", "mox", 248},
		{"0 %", "%", 0},
		{"24.8 %", "%", 24.8},
		{"-248 C", "C", -248.0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseFloatWithSuffix(tt.in, tt.suf)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %f,  expect: %f", tt, got, tt.ex)
			}
		})
	}
}

func TestConvCountString(t *testing.T) {
	var got int
	tests := []struct {
		in string
		ex int
	}{
		{"N/A", 0},
		{"mox", 0},
		{"0 %", 0},
		{"24.8 %", 0},
		{"-248", -248},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = convCountString(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d,  expect: %d", tt, got, tt.ex)
			}
		})
	}
}
