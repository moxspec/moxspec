package cpu

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseListString(t *testing.T) {
	var got []uint16
	tests := []struct {
		in string
		ex []uint16
	}{
		{"mox", nil},
		{"1", []uint16{1}},
		{"1-3", []uint16{1, 2, 3}},
		{"0,1-3", []uint16{0, 1, 2, 3}},
		{"0,1,2,3", []uint16{0, 1, 2, 3}},
		{"0,1-3,5", []uint16{0, 1, 2, 3, 5}},
		{"0,mox,5", []uint16{0, 5}},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseListString(tt.in)
			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %+v, expect: %+v", tt, got, tt.ex)
			}
		})
	}
}

func TestParseRangeString(t *testing.T) {
	var got []uint16
	tests := []struct {
		in string
		ex []uint16
	}{
		{"mox", nil},
		{"1", nil},
		{"1-mox", nil},
		{"1-3", []uint16{1, 2, 3}},
		{"3-1", []uint16{1, 2, 3}},
		{"3-3", []uint16{3}},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseRangeString(tt.in)
			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf("test: %+v, got: %+v, expect: %+v", tt, got, tt.ex)
			}
		})
	}
}
