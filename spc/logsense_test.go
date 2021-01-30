package spc

import (
	"reflect"
	"testing"
)

func TestParseReadLogSenseSupportPages(t *testing.T) {
	tests := []struct {
		path string
		ex   map[byte]struct{}
	}{
		{
			"logsense_support_log_page_seagate_savvio.txt",
			map[byte]struct{}{
				0x00: struct{}{},
				0x02: struct{}{},
				0x03: struct{}{},
				0x05: struct{}{},
				0x06: struct{}{},
				0x0D: struct{}{},
				0x0E: struct{}{},
				0x0F: struct{}{},
				0x10: struct{}{},
				0x15: struct{}{},
				0x18: struct{}{},
				0x1A: struct{}{},
				0x2F: struct{}{},
				0x37: struct{}{},
				0x38: struct{}{},
				0x3E: struct{}{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got, err := parseReadLogSenseSupportPages(bytes)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %s got %+v, expect %+v", test.path, got, test.ex)
			}
		})
	}
}

func TestParseReadLogSenseSupportPagesError(t *testing.T) {
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
				0x00, 0x00, 0x00, 0xFF, 0x00, 0x02, 0x03, 0x05,
				0x06, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x15, 0x18,
				0x19, 0x2f, 0x32, 0x34, 0x35,
			},
		},
		{
			"empty",
			[]byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := parseReadLogSenseSupportPages(test.bytes)
			if err == nil {
				t.Errorf("test: %s got %s, expect error", test.name, err)
			}

			if ret != nil {
				t.Errorf("test: %s got %+v, expect nil", test.name, ret)
			}
		})
	}
}

func TestParseLogSenseErrorCounterPages(t *testing.T) {
	tests := []struct {
		path string
		ex   map[int]uint64
	}{
		{
			"logsense_read_error_counter_seagate_savvio.txt",
			map[int]uint64{
				0x0000: 2848169391,
				0x0001: 0,
				0x0002: 0,
				0x0003: 2848169391,
				0x0004: 0,
				0x0005: 68674406166528,
				0x0006: 0,
			},
		},
		{
			"logsense_write_error_counter_seagate_savvio.txt",
			map[int]uint64{
				0x0001: 0,
				0x0002: 0,
				0x0003: 0,
				0x0004: 0,
				0x0005: 20233417674752,
				0x0006: 0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			bytes, err := loadTestData(test.path)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			got, err := parseLogSenseErrorCounterPages(bytes)
			if err != nil {
				t.Errorf("%s %s", test.path, err)
			}

			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %s got %+v, expect %+v", test.path, got, test.ex)
			}
		})
	}
}

func TestParseLogSenseErrorCounterPagesError(t *testing.T) {
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
				0x00, 0x00, 0x00, 0xFF, 0x00, 0x02, 0x03, 0x05,
				0x06, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x15, 0x18,
				0x19, 0x2f, 0x32, 0x34, 0x35,
			},
		},
		{
			"empty",
			[]byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret, err := parseLogSenseErrorCounterPages(test.bytes)
			if err == nil {
				t.Errorf("test: %s got %s, expect error", test.name, err)
			}

			if ret != nil {
				t.Errorf("test: %s got %+v, expect nil", test.name, ret)
			}
		})
	}
}

func TestParseLogSenseErrorCounterParamValue(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		ex   uint64
	}{
		{
			"empty",
			[]byte{},
			0,
		},
		{
			"normal case",
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			0x0102030405060708,
		},
		{
			"short data",
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
			0x010203040506,
		},
		{
			"long data",
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseLogSenseErrorCounterParamValue(test.in)
			if got != test.ex {
				t.Errorf("test: %s, got %+v, expect %+v", test.name, got, test.ex)
			}
		})
	}
}
