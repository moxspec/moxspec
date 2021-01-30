package spc

import (
	"reflect"
	"testing"
)

func TestParseLogExt(t *testing.T) {
	tests := []struct {
		path string
		ex   *logExtData
	}{
		{
			"logext_log_4_page_1_samsung_pm863.txt",
			&logExtData{
				0x00000000000077,
				0x000000000040CE,
				0x00000243F40D17,
				0x000000961B7627,
			},
		},
		{
			"logext_log_4_page_1_micron_5200.txt",
			&logExtData{
				0x0000000000001B,
				0x00000000001F7B,
				0x00000001C70A80,
				0x0000000005AF29,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got, err := parseLogExt(bytes)
			if err != nil {
				t.Errorf("test: %s got error, expect no error", test.path)
			}

			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %s got %+v, expect %+v", test.path, got, test.ex)
			}
		})
	}
}

func TestParseLogExtError(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{"empty", []byte{}},
		{"short data", []byte{0x00, 0x00}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseLogExt(test.bytes)
			if err == nil {
				t.Errorf("test: %s got no error, expect error", test.name)
			}
		})
	}
}
