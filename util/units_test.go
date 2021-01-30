package util

import (
	"fmt"
	"testing"
)

func TestConvUnit(t *testing.T) {
	var got string
	var err error
	tests := []struct {
		val float64
		mlt float64
		tgt float64
		fit bool
		ex  string
	}{
		{937703088.0 * 512.0, BaseDecimal, GIGA, false, "480.0GB"},
		{3125627568.0 * 512.0, BaseDecimal, GIGA, false, "1600.0GB"},
		{3125627568.0 * 512.0, BaseDecimal, GIGA, true, "1.6TB"},
		{1000 * 1000.0, BaseDecimal, MEGA, true, "1.0MB"},
		{2048.0 * 1024.0, BaseBinary, MEGA, false, "2.0MiB"},
		{2048.0 * 1024.0, BaseBinary, 256.0, false, ""},
		{0, BaseBinary, YOTTA, false, "0B"},
		{1024 * 1000, BaseBinary, YOTTA, false, "0B"},
		{1024 * 1000, BaseBinary, BYTE, false, "1024000B"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err = ConvUnit(tt.val, tt.mlt, tt.tgt, tt.fit)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s (err:%s)", tt, got, tt.ex, err)
			}
		})
	}
}
