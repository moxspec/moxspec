package netlink

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseRtnlLinkStats64(t *testing.T) {
	tests := []struct {
		path string
		ex   RtnlLinkStats64
	}{
		{
			path: "statsDump1.input",
			ex:   RtnlLinkStats64{},
		},
		{
			path: "statsDump2.input",
			ex: RtnlLinkStats64{
				RxPackets:         702538,
				TxPackets:         679505,
				RxBytes:           184340108,
				TxBytes:           203219567,
				RxErrors:          0,
				TxErrors:          0,
				RxDropped:         20,
				TxDropped:         0,
				Multicast:         704724292,
				Collisions:        0,
				RxLengthErrors:    0,
				RxOverErrors:      0,
				RxCrcErrors:       0,
				RxFrameErrors:     0,
				RxFifoErrors:      0,
				RxMissedErrors:    0,
				TxAbortedErrors:   0,
				TxCarrierErrors:   0,
				TxFifoErrors:      0,
				TxHeartbeatErrors: 0,
				TxWindowErrors:    0,
				RxCompressed:      0,
				TxCompressed:      0,
				RxNohandler:       0,
			},
		},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			got, err := parseRtnlLinkStats64(in)
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if !reflect.DeepEqual(*got, tt.ex) {
				t.Errorf("test: %+v, got: %v, expect: %v", tt, got, tt.ex)
			}
		})
	}
}
