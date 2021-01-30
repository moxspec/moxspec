package util

import (
	"fmt"
	"testing"
)

func TestBytesToUint32(t *testing.T) {
	var got uint32
	tests := []struct {
		in []byte
		ex uint32
	}{
		{[]byte{}, 0},
		{[]byte{'h', 'k', 'w', 'a', 'i'}, 0},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF}, 0xFFFFFFFF},
		{[]byte{0x01, 0x02, 0x03, 0x04}, 0x04030201},
		{[]byte{0x01, 0x02, 0x03}, 0x00030201},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = BytesToUint32(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}

}

func TestBytesToUint64(t *testing.T) {
	var got uint64
	tests := []struct {
		in []byte
		ex uint64
	}{
		{[]byte{}, 0},
		{[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}, 0},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 0xFFFFFFFFFFFFFFFF},
		{[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, 0x0807060504030201},
		{[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 0x0007060504030201},
		{[]byte{0x01, 0x02, 0x03}, 0x0000000000030201},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = BytesToUint64(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}

}
