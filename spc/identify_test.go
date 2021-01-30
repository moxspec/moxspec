package spc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseIdentify(t *testing.T) {
	tests := []struct {
		path string
		ex   *identifyData
	}{
		{
			"identify_samsung_pm863.txt",
			&identifyData{
				"S2TWNX0K200947",
				"GXT5204Q",
				"SAMSUNG MZ7LM240HMHQ-00005",
				"6.0Gb/s",
				"6.0Gb/s",
				"2.5\"",
				1,
				"SATA 3.1",
				true,
				true,
			},
		},
		{
			"identify_micron_5200.txt",
			&identifyData{
				"19342374BAF1",
				"D1MU820",
				"Micron_5200_MTFDDAK7T6TDC",
				"6.0Gb/s",
				"6.0Gb/s",
				"2.5\"",
				1,
				"SATA 3.2",
				true,
				true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got := parseIdentify(bytes)
			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %s got %+v, expect %+v", test.path, got, test.ex)
			}
		})
	}
}

func TestParseIdentifyError(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{"empty", []byte{}},
		{"short data", []byte{0x00, 0x00}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseIdentify(test.bytes)
			if got != nil {
				t.Errorf("test: %s got %+v, expect nil", test.name, got)
			}
		})
	}
}

func TestDecodeTransport(t *testing.T) {
	var got string
	tests := []struct {
		in uint16
		ex string
	}{
		{0, ""},
		{0x1, "ATA8-AST"},
		{0x200, "SATA 3.4"},
		{0x400, ""},
		{0x800, ""},
		{0xFFFF, ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = decodeTransport(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: cur=%s", tt, got, tt.ex)
			}
		})
	}
}

func TestReadWord(t *testing.T) {
	var got uint16
	tests := []struct {
		buf []byte
		at  int
		ex  uint16
	}{
		{[]byte{0, 0}, 1, 0},
		{[]byte{0, 0}, 8, 0},
		{[]byte{0, 0}, -1, 0},
		{[]byte{0x21, 0x07}, 0, 0x0721},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = readWord(tt.buf, tt.at)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestReadWordsAsATAString(t *testing.T) {
	var got string
	tests := []struct {
		buf  []byte
		from int
		to   int
		ex   string
	}{
		{[]byte{'i', 't', 'g', 'o', ' ', 'a', 'a', 'p', 's', 's'}, 0, 0, "ti"},         // tioga pass
		{[]byte{'i', 't', 'g', 'o', ' ', 'a', 'a', 'p', 's', 's'}, 4, 4, "ss"},         // tioga pass
		{[]byte{'4', '2', '2', '8', '8', '4'}, 0, 2, "248248"},                         // 248248
		{[]byte{'4', '2', '2', '8', '8', '4'}, 0, -1, ""},                              // 248248
		{[]byte{'4', '2', '2', '8', '8', '4'}, -1, 1, ""},                              // 248248
		{[]byte{'4', '2', '2', '8', '8', '4'}, -1, 3, ""},                              // 248248
		{[]byte{'4', '2', '2', '8', '8', '4'}, -248, 248, ""},                          // 248248
		{[]byte{'e', 'w', 'a', 'n', 'c', 't', 'e', 'h', 'e'}, 0, 4, ""},                // wenatchee
		{[]byte{'e', 'w', 'a', 'n', 'c', 't', 'e', 'h', 0x00, 'e'}, 0, 4, "wenatchee"}, // wenatchee
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = readWordsAsATAString(tt.buf, tt.from, tt.to)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}
