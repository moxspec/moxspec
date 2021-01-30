package hpacucli

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
		{"ok", false},
		{"OK", true},
		{"OKK", false},
		{" OK ", false},
		{"mox", false},
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
