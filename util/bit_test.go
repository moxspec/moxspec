package util

import (
	"fmt"
	"testing"
)

func TestFindMSB(t *testing.T) {
	var err error
	var got int
	tests := []struct {
		in interface{}
		ex int
	}{
		{false, -1},
		{0, -1},
		{2, 1},
		{4, 2},
		{8, 3},
		{248, 7},   // As you know, 248 is the most important number of the world
		{0x248, 9}, // As you know, 248 is the most important number of the world
		{"mox", -1},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err = FindMSB(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d (%s), expect: %d", tt, got, err, tt.ex)
			}
		})
	}

}
