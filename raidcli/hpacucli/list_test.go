package hpacucli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/moxspec/moxspec/raidcli"
)

func TestParseCtlList(t *testing.T) {
	tests := []struct {
		path string
		ex   []*Controller
	}{
		{"ctrl_all_show.sa.input", []*Controller{
			&Controller{
				Slot: "2",
				ControllerSpec: raidcli.ControllerSpec{
					ProductName: "Smart Array P420",
				},
			},
		}},
		{"ctrl_all_show.hba.input", []*Controller{
			&Controller{
				Slot: "0",
				ControllerSpec: raidcli.ControllerSpec{
					ProductName: "Smart HBA H240ar",
				},
			},
		}},
		{"ctrl_all_show.dynsa.input", []*Controller{
			&Controller{
				Slot: "0",
				ControllerSpec: raidcli.ControllerSpec{
					ProductName: "Dynamic Smart Array B120i RAID",
				},
			},
			&Controller{
				Slot: "1",
				ControllerSpec: raidcli.ControllerSpec{
					ProductName: "Smart Array P420i",
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

			got, err := parseCtlList(string(in))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(got) != len(tt.ex) {
				t.Errorf("%s len(got):%d, len(ex):%d", tt.path, len(got), len(tt.ex))
			}

			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i], tt.ex[i]) {
					t.Errorf(fmt.Sprintf("%s %s", tt.path, pretty.Compare(got[i], tt.ex[i])))
				}
			}
		})
	}
}
