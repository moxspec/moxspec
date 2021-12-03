package ipmi

import (
	"fmt"
	"testing"
)

func TestParseMcInfoRaw(t *testing.T) {
	var rev string
	tests := []struct {
		in string
		ex string
	}{
		{" 01 81 18 08 02 af db 07 00 01 00 19 00 10 00", "24.08"},
		{" 01 81 02 48 02 af db 07 00 01 00 19 00 10 00", "2.48"},
		{" 20 01 01 17 02 bf af 2b 00 02 02 12 00 17 00", "1.17"},
		{" 01 81 00 00 02 af db 07 00 01 00 19 00 10 00", "0.00"},
		{" 20 81 0c 03 02 bf af 2b 00 05 00 00 00 00 00", "12.03"},
		{" 20 81 ff 03 02 bf af 2b 00 05 00 00 00 00 00", "127.03"},
		{" 01 81 02 92 02 af db 07 00 01 00 19 00 10 00", "2.92"},
		{" 01 81 02", ""},
		{" 01 81 02 48", "2.48"},
		{"248", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			rev = parseMcInfoRaw(tt.in)
			if rev != tt.ex {
				t.Errorf("got: %s, expected: %s", rev, tt.ex)
			}
		})
	}
}

func TestParseMacAddressRaw(t *testing.T) {
	var rev string
	tests := []struct {
		in string
		ex string
	}{
		{" 11 5c 54 6d 0b 54 6b", "5c:54:6d:0b:54:6b"},
		{" 11 0c 42 a1 4d 91 c4", "0c:42:a1:4d:91:c4"},
		{" 11 5c 54 6d 0b 54 6b 24", ""},
		{" 24 82 48 ", ""},
		{"248", ""},
		{"", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			rev = parseMACAddressRaw(tt.in)
			if rev != tt.ex {
				t.Errorf("got: %s, expected: %s", rev, tt.ex)
			}
		})
	}
}

func TestParseIPAddressRaw(t *testing.T) {
	var rev string
	tests := []struct {
		in string
		ex string
	}{
		{" 11 00 00 00 00", "0.0.0.0"},
		{" 11 ff 00 00 00", "255.0.0.0"},
		{" 11 0a 17 fd 8b", "10.23.253.139"},
		{" 24 82 48 ", ""},
		{"248", ""},
		{"", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			rev = parseIPAddressRaw(tt.in)
			if rev != tt.ex {
				t.Errorf("got: %s, expected: %s", rev, tt.ex)
			}
		})
	}
}

func TestParseVLANID(t *testing.T) {
	var rev uint16
	tests := []struct {
		in string
		ex uint16
	}{
		{" 11", 0},
		{" 11 c8 80", 200},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			rev = parseVLANID(tt.in)
			if rev != tt.ex {
				t.Errorf("got: %d, expected: %d", rev, tt.ex)
			}
		})
	}
}

func TestParseAddrSrc(t *testing.T) {
	var rev addrSrcType
	tests := []struct {
		in string
		ex addrSrcType
	}{
		{"", srcUnspecified},
		{" 11", srcUnspecified},
		{" 11 00", srcUnspecified},
		{" 11 00 00", srcUnspecified},
		{" 11 01", srcStatic},
		{" 11 02", srcDHCP},
		{" 11 03", srcBIOS},
		{" 11 04", srcOther},
		{" 11 05", srcUnspecified},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			rev = parseAddressSource(tt.in)
			if rev != tt.ex {
				t.Errorf("got: %s, expected: %s", rev, tt.ex)
			}
		})
	}
}
