package raidcli

import (
	"errors"
	"fmt"
	"testing"
)

func TestSplitKeyVal(t *testing.T) {
	var key, val string
	var err error
	tests := []struct {
		in    string
		delim string
		key   string
		val   string
		err   error
	}{
		{"Bus Number : 5", ":", "Bus Number", "5", nil},
		{"Number of enclosures on adapter 0 -- 1", ":", "", "", errors.New("dummy")},
		{"Drive's position: DiskGroup: 0, Span: 0, Arm: 0", ":", "Drive's position", "DiskGroup: 0, Span: 0, Arm: 0", nil},
		{"mox:248:2:4:8", ":", "mox", "248:2:4:8", nil},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			key, val, err = SplitKeyVal(tt.in, tt.delim)
			if tt.err == nil && err != nil {
				t.Errorf("test:%+v, error should be nil, got: %s", tt, err)
			}

			if tt.err != nil && err == nil {
				t.Errorf("test:%+v, error should NOT be nil", tt)
			}

			if key != tt.key || val != tt.val {
				t.Errorf("test: %+v, got: key=%s, val=%s expect: key=%s, val=%s", tt, key, val, tt.key, tt.val)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	var got uint64
	tests := []struct {
		in string
		ex uint64
	}{
		{"2 4 8", 0},
		{"2.48 KB", 2539},
		{"2.48 MB", 2600468},
		{"2.48 GB", 2662879723},
		{"2.48 TB", 2726788836884},
		{"2.48 PB", 2792231768969707},
		{"2.48 EB", 2859245331424980480},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = ParseSize(tt.in, Binary)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
