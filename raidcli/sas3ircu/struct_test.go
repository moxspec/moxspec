package sas3ircu

import (
	"fmt"
	"testing"
)

func TestIsHealthy(t *testing.T) {
	var got bool
	tests := []struct {
		in string
		ex bool
	}{
		{"248 (OKY)", true},
		{"(OKY)", true},
		{"Inactive, Okay (OKY)", true},
		{"Failed (FLD)", false},
		{"248", false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = isHealthy(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t, expect: %t", tt, got, tt.ex)
			}
		})
	}
}
