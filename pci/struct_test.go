package pci

import (
	"fmt"
	"testing"
)

func TestConfigReadByte(t *testing.T) {
	var conf = Config{
		path: "",
		br: []byte{
			0x86, 0x80,
		},
	}

	var got byte
	tests := []struct {
		in uint16
		ex byte
	}{
		{0, 0x86},
		{1, 0x80},
		{2, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = conf.ReadByteFrom(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestConfigReadWord(t *testing.T) {
	var conf = Config{
		path: "",
		br: []byte{
			0x86, 0x80, 0xFF, 0x00,
		},
	}

	var got uint16
	tests := []struct {
		in uint16
		ex uint16
	}{
		{0, 0x8086},
		{1, 0xFF80},
		{2, 0x00FF},
		{3, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = conf.ReadWordFrom(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestConfigReadDWord(t *testing.T) {
	var conf = Config{
		path: "",
		br: []byte{
			0x86, 0x80, 0xFF, 0x00, 0xE5, 0x19,
		},
	}

	var got uint32
	tests := []struct {
		in uint16
		ex uint32
	}{
		{0, (0x00FF << 16) | 0x8086},
		{2, (0x19E5 << 16) | 0x00FF},
		{3, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = conf.ReadDWordFrom(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
