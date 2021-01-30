package spc

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestMakeReadLogSenseCmd(t *testing.T) {
	tests := []struct {
		page    byte
		subpage byte
		length  byte
		ex      []byte
	}{
		{
			0xFF,
			0xEE,
			0xDD,
			[]byte{
				0x4D, 0x00, 0xFF, 0xEE, 0x00,
				0x00, 0x00, 0x00, 0xDD, 0x00,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			got := makeReadLogSenseCmd(test.page, test.subpage, test.length)
			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %+v, got: %+v, expect: %+v", test, got, test.ex)
			}
		})
	}
}

func TestMakeInquiryCmd(t *testing.T) {
	tests := []struct {
		page   byte
		evpd   bool
		length int
		ex     []byte
	}{
		{0xFF, true, 0x1234, []byte{0x12, 0x01, 0xFF, 0x12, 0x34, 0x00}},
		{0xEE, false, 0x1100, []byte{0x12, 0x00, 0xEE, 0x11, 0x00, 0x00}},
		{0x00, true, -1, nil},
		{0x00, true, math.MaxUint16 + 1, nil},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			got := makeInquiryCmd(test.page, test.evpd, test.length)
			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %+v, got: %+v, expect: %+v", test, got, test.ex)
			}
		})
	}
}

func TestMakeATAPThruCmd(t *testing.T) {
	tests := []struct {
		feature byte
		lbaLow  byte
		lbaMid  byte
		lbaHigh byte
		cmd     byte
		ex      []byte
	}{
		{
			0xFF,
			0xEE,
			0xDD,
			0xCC,
			0xBB,
			[]byte{
				0x85, 0x08, 0x0E, 0x00, 0xFF, 0x00, 0x01, 0x00,
				0xEE, 0x00, 0xDD, 0x00, 0xCC, 0x00, 0xBB, 0x00,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test), func(t *testing.T) {
			got := makeATAPThruCmd(test.feature, test.lbaLow, test.lbaMid, test.lbaHigh, test.cmd)
			if !reflect.DeepEqual(got, test.ex) {
				t.Errorf("test: %+v, got: %+v, expect: %+v", test, got, test.ex)
			}
		})
	}
}
