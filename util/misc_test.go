package util

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestDumpBinary(t *testing.T) {
	var got string
	tests := []struct {
		in interface{}
		ex string
	}{
		{uint8(0x00), "0000 0000"},
		{uint8(0x0F), "0000 1111"},
		{uint8(0xFF), "1111 1111"},
		{uint16(0xFFFF), "1111 1111 1111 1111"},
		{uint32(0x0000FFFF), "0000 0000 0000 0000 1111 1111 1111 1111"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = DumpBinary(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestBlkLabelAscSorter(t *testing.T) {
	in := []string{
		"sdaz",
		"sdb",
		"sdaa",
		"sdc",
		"sdab",
		"sda",
		"sdz",
		"sdad",
		"sdac",
	}

	ex := []string{
		"sda",
		"sdb",
		"sdc",
		"sdz",
		"sdaa",
		"sdab",
		"sdac",
		"sdad",
		"sdaz",
	}

	sort.Slice(in, func(i, j int) bool {
		return BlkLabelAscSorter(in[i], in[j])
	})

	if !reflect.DeepEqual(in, ex) {
		t.Errorf("got: %v, expect: %v", in, ex)
	}
}

func TestIPv4MaskSize(t *testing.T) {
	var got int
	tests := []struct {
		in string
		ex int
	}{
		{"255.0.0.0", 8},
		{"255.128.0.0", 9},
		{"255.255.0.0", 16},
		{"255.255.128.0", 17},
		{"255.255.255.0", 24},
		{"255.255.255.128", 25},
		{"255.255.255.255", 32},
		{"255.0.255.255", 0},
		{"mox", 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = IPv4MaskSize(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
