package smbios

import (
	"fmt"
	"testing"

	gosmbios "github.com/digitalocean/go-smbios/smbios"
)

//
// type Structure struct {
// 	Header    Header
// 	Formatted []byte
// 	Strings   []string
// }
//
// ref.  https://github.com/digitalocean/go-smbios/blob/master/smbios/structure.go
//

func TestGetByte(t *testing.T) {
	var st = &gosmbios.Structure{
		Header: gosmbios.Header{},
		Formatted: []byte{
			0x86, 0x80,
		},
		Strings: nil,
	}

	var got uint8
	tests := []struct {
		in int
		ex uint8
	}{
		{headerSize, 0x86},
		{headerSize + 1, 0x80},
		{headerSize + 2, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = getByte(st, tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestGetWord(t *testing.T) {
	var st = &gosmbios.Structure{
		Header: gosmbios.Header{},
		Formatted: []byte{
			0x86, 0x80, 0xFF, 0x00,
		},
		Strings: nil,
	}

	var got uint16
	tests := []struct {
		in int
		ex uint16
	}{
		{headerSize, 0x8086},
		{headerSize + 1, 0xFF80},
		{headerSize + 2, 0x00FF},
		{headerSize + 3, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = getWord(st, tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestGetDWord(t *testing.T) {
	var st = &gosmbios.Structure{
		Header: gosmbios.Header{},
		Formatted: []byte{
			0x86, 0x80, 0xFF, 0x00,
		},
		Strings: nil,
	}

	var got uint32
	tests := []struct {
		in int
		ex uint32
	}{
		{headerSize, (0x00FF << 16) | 0x8086},
		{headerSize + 1, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = getDWord(st, tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestGetQWord(t *testing.T) {
	var st = &gosmbios.Structure{
		Header: gosmbios.Header{},
		Formatted: []byte{
			0x86, 0x80, 0xFF, 0x00, 0xE5, 0x19, 0xFF, 0x00,
		},
		Strings: nil,
	}

	var got uint64
	tests := []struct {
		in int
		ex uint64
	}{
		{headerSize, (((0x00FF << 16) | 0x19E5) << 32) | ((0x00FF << 16) | 0x8086)},
		{headerSize + 1, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = getQWord(st, tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
