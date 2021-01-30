package scsi

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestDecodePg80(t *testing.T) {
	tests := []struct {
		path       string
		sn         string
		wantsError bool
	}{
		{"vpd_pg80.megaraid.jbod.input", "2YGAPMLD", false},
		{"vpd_pg80.samsung.ssd.input", "S2UJNX0JC00313", false},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {

			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			sn, err := decodePg80(in)

			if tt.wantsError && err == nil {
				t.Errorf("%s wants error but got nil", tt.path)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if sn != tt.sn {
				t.Errorf("%s got: %s, expect: %s", tt.path, sn, tt.sn)
			}
		})
	}
}

func TestDecodePg83(t *testing.T) {
	tests := []struct {
		path       string
		numDesc    int
		wantsError bool
	}{
		{"vpd_pg83.megaraid.jbod.input", 5, false},
		{"vpd_pg83.samsung.ssd.input", 3, false},
		{"vpd_pg83.megaraid.SAS-3.3108.input", 1, false},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			descs, err := decodePg83(in)

			if tt.wantsError && err == nil {
				t.Errorf("%s wants error but got nil", tt.path)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(descs) != tt.numDesc {
				t.Errorf("%s len(descs): %d, expect: %d", tt.path, len(descs), tt.numDesc)
			}
		})
	}
}

func TestDecodeSCSINameString(t *testing.T) {
	tests := []struct {
		in string
		ex string
	}{
		{"naa.5000CCA2731373AB", "5000CCA2731373AB"},
		{"248", ""},
		{"248.248", ""},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.in, func(t *testing.T) {
			got := decodeSCSINameString(tt.in)
			if got != tt.ex {
				t.Errorf("%s got: %s, expect: %s", tt.in, got, tt.ex)
			}
		})
	}
}

func TestDecodeSPC5ID(t *testing.T) {
	tests := []struct {
		in []byte
		ex string
	}{
		{
			[]byte{2, 4, 8, 2, 4, 8, 2, 4, 8, 2, 4, 8, 2, 4, 8, 2, 4, 8},
			"08020408-0204-0802-0408-020408020408",
		},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got := decodeSPC5ID(tt.in)
			if got != tt.ex {
				t.Errorf("in: %v, got: %s, expect: %s", tt.in, got, tt.ex)
			}
		})
	}
}

func TestFormatNAABytes(t *testing.T) {
	tests := []struct {
		in []byte
		l  int
		ex string
	}{
		{[]byte{2, 4, 8}, 248, ""},
		{[]byte{2, 4, 8}, 0, ""},
		{[]byte{2, 4, 8}, -248, ""},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got := formatNAABytes(tt.in, tt.l)
			if got != tt.ex {
				t.Errorf("in: %v l: %d, got: %s, expect: %s", tt.in, tt.l, got, tt.ex)
			}
		})
	}
}
