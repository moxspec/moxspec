package nvme

import (
	"fmt"
	"testing"
)

func TestParseNamespaceID(t *testing.T) {
	var got int
	tests := []struct {
		in string
		ex int
	}{
		{"nvme0n1", 1},
		{"nvme10n1", 1},
		{"nvme0n1p1", 0},
		{"nvme0", 0},
		{"mox", 0},
		{"", 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseNamespaceID(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
