package util

import (
	"fmt"
	"testing"
)

func TestHasPrefixIn(t *testing.T) {
	var got bool
	tests := []struct {
		in  string
		pre []string
		ex  bool
	}{
		{"", []string{"mox"}, false},
		{"fail.pattern", []string{"mox"}, false},
		{"pass.pattern", []string{"2", "4", "8", "pass"}, true},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = HasPrefixIn(tt.in, tt.pre...)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t, expect: %t", tt, got, tt.ex)
			}
		})
	}
}
