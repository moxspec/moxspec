package pci

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseBitDefs(t *testing.T) {
	var defs = map[byte]string{
		1:  "mox",
		3:  "248",
		4:  "m",
		6:  "248",
		8:  "o",
		9:  "248",
		10: "x",
		11: "mox.248",
		12: "248.mox",
		16: "mox",
	}
	var got []string
	tests := []struct {
		in uint32
		ex []string
	}{
		{0x00, nil},
		{0x02, []string{"mox"}},
		{0x15, []string{"m"}},
		{0x100, []string{"o"}},
		{0x248, []string{"248", "248", "248"}},
		{0x400, []string{"x"}},
		{0x800, []string{"mox.248"}},
		{0x1000, []string{"248.mox"}},
		{0x10000, []string{"mox"}},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseBitDefs(tt.in, defs)
			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %v, expect: %v", tt, got, tt.ex)
			}
		})
	}
}

func TestReadSysfsValue(t *testing.T) {
	tests := []struct {
		object  string
		bit     int
		ex      uint64
		isError bool
	}{
		{"sysfs_value_0x0248", 16, 0x248, false},
		{"sysfs_value_0x0248", 32, 0x248, false},
		{"sysfs_value_0x12345678", 16, 0, true},
		{"sysfs_value_0x12345678", 32, 0x12345678, false},
		{"sysfs_value_invalid_1", 16, 0, true},
		{"sysfs_value_invalid_2", 16, 0, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			got, err := readSysfsValue("./testdata", test.object, test.bit)

			if err == nil && test.isError {
				t.Errorf("test: %+v, got: no error, expect: error", test)
			}
			if err != nil && !test.isError {
				t.Errorf("test: %+v, got: error, expect: no error", test)
			}

			if got != test.ex {
				t.Errorf("test: %+v, got: %v, expect: %v", test, got, test.ex)
			}
		})
	}
}
