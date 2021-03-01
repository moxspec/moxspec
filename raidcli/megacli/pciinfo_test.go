package megacli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/moxspec/moxspec/raidcli"
	"github.com/kylelemons/godebug/pretty"
)

func TestParsePCIInfoSingle(t *testing.T) {
	tests := []struct {
		path string
		ex   []*Controller
	}{
		{"adpgetpciinfo.empty.input", []*Controller{}},
		{"adpgetpciinfo.single.input", []*Controller{
			&Controller{
				ControllerSpec: raidcli.ControllerSpec{
					Number:   0,
					Bus:      1,
					Device:   0,
					Function: 0,
				},
			},
		}},
		{"adpgetpciinfo.twin.input", []*Controller{
			&Controller{
				ControllerSpec: raidcli.ControllerSpec{
					Number:   0,
					Bus:      1,
					Device:   0,
					Function: 0,
				},
			},
			&Controller{
				ControllerSpec: raidcli.ControllerSpec{
					Number:   1,
					Bus:      2,
					Device:   0xFF,
					Function: 0,
				},
			},
		}},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			got, err := parsePCIInfo(string(in))
			if err != nil {
				t.Errorf("%s", err)
			}

			if len(got) != len(tt.ex) {
				t.Errorf("%s len(got):%d, len(ex):%d", tt.path, len(got), len(tt.ex))
			}

			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i], tt.ex[i]) {
					t.Errorf(pretty.Compare(got[i], tt.ex[i]))
				}
			}
		})
	}
}

func TestParseControllerNumber(t *testing.T) {
	var got int
	tests := []struct {
		in string
		ex int
	}{
		{"mox", -1},
		{"PCI information for Controller 0", 0},
		{"PCI information for Controller -10000", -1},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, _ = parseControllerNumber(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d, expect: %d", tt, got, tt.ex)
			}
		})
	}
}
