package util

import (
	"fmt"
	"testing"
)

func TestRound(t *testing.T) {
	var got float64
	tests := []struct {
		in float64
		pt int
		ex float64
	}{
		{2.48, 2, 2.48},
		{0.248, 1, 0.2},
		{0, 1, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = Round(tt.in, tt.pt)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %f, expect: %f", tt, got, tt.ex)
			}
		})
	}

}
