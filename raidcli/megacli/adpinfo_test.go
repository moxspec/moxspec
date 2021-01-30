package megacli

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseAdpInfo(t *testing.T) {
	tests := []struct {
		path       string
		name       string
		sn         string
		bios       string
		fw         string
		bbu        bool
		wantsError bool
	}{
		{
			"adpallinfo.MegaRAID_SAS_9270-8i.input",
			"LSI MegaRAID SAS 9270-8i",
			"SV54952556",
			"5.50.03.0_4.17.08.00_0x06110200",
			"3.460.05-4565",
			true,
			false,
		},
		{
			"adpallinfo.PERC_H730P_Mini.input",
			"PERC H730P Mini",
			"85B01NJ",
			"6.33.01.0_4.16.07.00_0x06120304",
			"4.290.00-8334",
			true,
			false,
		},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			name, sn, bios, fw, bbu, err := parseAdpInfo(string(in))

			if tt.wantsError && err == nil {
				t.Errorf("%s wants error but got nil", tt.path)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if name != tt.name {
				t.Errorf("%s name: got: %s, expect: %s", tt.path, name, tt.name)
			}

			if fw != tt.fw {
				t.Errorf("%s firm: got: %s, expect: %s", tt.path, fw, tt.fw)
			}

			if bios != tt.bios {
				t.Errorf("%s bios: got: %s, expect: %s", tt.path, bios, tt.bios)
			}

			if sn != tt.sn {
				t.Errorf("%s sn: got: %s, expect: %s", tt.path, sn, tt.sn)
			}

			if bbu != tt.bbu {
				t.Errorf("%s bbu: got: %t, expect: %t", tt.path, bbu, tt.bbu)
			}
		})
	}
}
