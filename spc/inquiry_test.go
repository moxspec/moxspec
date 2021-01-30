package spc

import (
	"reflect"
	"testing"
)

func TestParseInquiryData(t *testing.T) {
	tests := []struct {
		ex   *inquiryData
		path string
	}{
		{
			&inquiryData{
				"SEAGATE",
				"ST300MM0006",
				"B001",
			},
			"inquiry_seagate_savvio.txt",
		},
		{
			&inquiryData{
				"ATA",
				"Micron_5200_MTFD",
				"U820",
			},
			"inquiry_micron_5200.txt",
		},
		{
			&inquiryData{
				"ATA",
				"SAMSUNG MZ7LM240",
				"204Q",
			},
			"inquiry_samsung_pm863.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got, err := parseInquiryData(bytes)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %+v, got %+v, expect %+v", test, got, test.ex)
			}
		})
	}

}

func TestParseInquiryDataError(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{
			"short data",
			[]byte{
				0x00, 0x00, 0x00,
			},
		},
		{
			"invalid length field",
			[]byte{
				// length field + 4 > buffer len
				0x00, 0x00, 0x06, 0x12, 0xFF, 0x00, 0x00, 0x02,
				0x41, 0x54, 0x41, 0x20, 0x20, 0x20, 0x20, 0x20,
				0x53, 0x41, 0x4D, 0x53, 0x55, 0x4E, 0x47, 0x20,
				0x4D, 0x5A, 0x37, 0x4C, 0x4D, 0x32, 0x34, 0x30,
				0x32, 0x30, 0x34, 0x51, 0x53, 0x32, 0x54, 0x57,
			},
		},
		{
			"empty",
			[]byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := parseInquiryData(test.bytes)
			if err == nil {
				t.Errorf("test: %s got %s, expect error", test.name, err)
			}

			if ret != nil {
				t.Errorf("test: %s got %s, expect nil", test.name, ret)
			}
		})
	}
}

func TestParseInquiryVPDSerialNumber(t *testing.T) {
	tests := []struct {
		ex   string
		path string
	}{
		{
			"S0K5MZZ40000K62872UR",
			"vpd_sn_seagate_savvio.txt",
		},
		{
			"19342374BAF1",
			"vpd_sn_micron_5200.txt",
		},
		{
			"S2TWNX0K200947",
			"vpd_sn_samsung_pm863.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got, err := parseInquiryVPDSerialNumber(bytes)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			if got != test.ex {
				t.Errorf("test: %s got %s, expect %s", test.path, got, test.ex)
			}
		})
	}
}

func TestParseInquiryVPDSerialNumberError(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{
			"short data",
			[]byte{
				0x00, 0x00, 0x00,
			},
		},
		{
			"invalid length field",
			[]byte{
				// length field + 4 > buffer len
				0x00, 0x80, 0x00, 0xFF, 0x30, 0x33, 0x32, 0x56,
				0x55, 0x4B, 0x46, 0x53, 0x4A, 0x43, 0x30, 0x30,
				0x30, 0x36, 0x31, 0x30, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			"empty",
			[]byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := parseInquiryVPDSerialNumber(test.bytes)
			if err == nil {
				t.Errorf("test: %s got %s, expect error", test.name, err)
			}

			if ret != "" {
				t.Errorf("test: %s got %s, expect empty string", test.name, ret)
			}
		})
	}
}

func TestParseInquiryVPDSupportedPages(t *testing.T) {
	tests := []struct {
		path string
		ex   map[byte]struct{}
	}{
		{
			"vpd_sv_seagate_savvio.txt",
			map[byte]struct{}{
				0x00: struct{}{},
				0x80: struct{}{},
				0x83: struct{}{},
				0x86: struct{}{},
				0x87: struct{}{},
				0x88: struct{}{},
				0x8A: struct{}{},
				0x90: struct{}{},
				0xB0: struct{}{},
				0xB1: struct{}{},
				0xB2: struct{}{},
				0xC0: struct{}{},
				0xC1: struct{}{},
				0xC3: struct{}{},
				0xD1: struct{}{},
				0xD2: struct{}{},
			},
		},
		{
			"vpd_sv_micron_5200.txt",
			map[byte]struct{}{
				0x00: struct{}{},
				0x80: struct{}{},
				0x83: struct{}{},
				0x87: struct{}{},
				0x89: struct{}{},
				0xB0: struct{}{},
				0xB1: struct{}{},
				0xB2: struct{}{},
			},
		},
		{
			"vpd_sv_samsung_pm863.txt",
			map[byte]struct{}{
				0x00: struct{}{},
				0x80: struct{}{},
				0x83: struct{}{},
				0x87: struct{}{},
				0x89: struct{}{},
				0x8A: struct{}{},
				0xB0: struct{}{},
				0xB1: struct{}{},
				0xB2: struct{}{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got, err := parseInquiryVPDSupportedPages(bytes)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %+v, got %+v, expect %+v", test, got, test.ex)
			}
		})
	}
}

func TestParseInquiryVPDSupportedPagesError(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{
			"short data",
			[]byte{
				0x00, 0x00, 0x00,
			},
		},
		{
			"invalid length field",
			[]byte{
				// length field + 4 > buffer len
				0x00, 0x00, 0x00, 0xFF, 0x00, 0x80, 0x83, 0x86,
				0x87, 0x88, 0x8A, 0x8D, 0x90, 0x91, 0xB0, 0xB1,
				0xB2, 0xB7,
			},
		},
		{
			"empty",
			[]byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := parseInquiryVPDSupportedPages(test.bytes)
			if err == nil {
				t.Errorf("test: %s got %s, expect error", test.name, err)
			}

			if ret != nil {
				t.Errorf("test: %s got %+v, expect nil", test.name, ret)
			}
		})
	}
}
