package platform

import (
	"fmt"
	"testing"
)

func TestScanModules(t *testing.T) {
	var got bool
	tests := []struct {
		list string
		name string
		ex   bool
	}{
		{"", "mox", false},
		{"mox", "", false},
		{`dm_mirror 22289 0 - Live 0xffffffffc02a7000
		dm_mod 123941 2 dm_mirror,dm_log, Live 0xffffffffc0273000`, "dm_mod", true},
		{`dm_mirror 22289 0 - Live 0xffffffffc02a7000
		dm_mod 123941 2 dm_mirror,dm_log, Live 0xffffffffc0273000`, "mox", false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = scanModules(tt.list, tt.name)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t, expect: %t", tt, got, tt.ex)
			}
		})
	}
}
