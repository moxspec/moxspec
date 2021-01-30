package blk

import (
	"fmt"
	"testing"
)

func TestParseScheduler(t *testing.T) {
	var got string
	tests := []struct {
		in string
		ex string
	}{
		{"[none] mq-deadline kyber", "none"},
		{"none", "none"},
		{"noop anticipatory deadline [cfq]", "cfq"},
		{"", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseScheduler(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestParseSCSIAddress(t *testing.T) {
	var h, c, tg, l uint16
	tests := []struct {
		in string
		h  uint16
		c  uint16
		tg uint16
		l  uint16
	}{
		{"1:0:0:0", 1, 0, 0, 0},
		{"1:-1:0:0", 0, 0, 0, 0},
		{"0:0:0", 0, 0, 0, 0},
		{"mox", 0, 0, 0, 0},
		{"m:o:x", 0, 0, 0, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			h, c, tg, l = parseSCSIAddress(tt.in)
			if h != tt.h || c != tt.c || tg != tt.tg || l != tt.l {
				t.Errorf("test: %+v, got: %d %d %d %d, expect: %d %d %d %d", tt, h, c, tg, l, tt.h, tt.c, tt.tg, tt.l)
			}
		})
	}
}

func TestParseNodeNumber(t *testing.T) {
	var gotMaj, gotMin uint16
	tests := []struct {
		in    string
		exMaj uint16
		exMin uint16
	}{
		{"259:0", 259, 0},
		{"9::0", 0, 0},
		{"mox", 0, 0},
		{"3:mox", 0, 0},
		{"mox:2", 0, 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			gotMaj, gotMin = parseNodeNumber(tt.in)
			if gotMaj != tt.exMaj || gotMin != tt.exMin {
				t.Errorf("test: %+v, got: %d, %d, expect: %d, %d", tt, gotMaj, gotMin, tt.exMaj, tt.exMin)
			}
		})
	}
}
