package megacli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/moxspec/moxspec/raidcli"
)

func TestParseLDList(t *testing.T) {
	tests := []struct {
		path string
		ex   []*LogDrive
	}{
		{
			"ldlist.novol.input", []*LogDrive{},
		},
		{
			"ldlist.3vol.input", []*LogDrive{
				&LogDrive{
					GroupID:  0,
					TargetID: 0,
					LogDriveSpec: raidcli.LogDriveSpec{
						Label:       "grp:0",
						RAIDLv:      raidcli.RAID1,
						Size:        238998823895, // 222.585*1024*1024*1024
						State:       "Optimal",
						StripSize:   262144,
						CachePolicy: "WriteThrough",
					},
				},
				&LogDrive{
					GroupID:  1,
					TargetID: 1,
					LogDriveSpec: raidcli.LogDriveSpec{
						Label:       "grp:1",
						RAIDLv:      raidcli.RAID0,
						Size:        9998958742994, // 9.094*1024*1024*1024*1024
						State:       "Optimal",
						StripSize:   262144,
						CachePolicy: "WriteThrough",
					},
				},
				&LogDrive{
					GroupID:  2,
					TargetID: 2,
					LogDriveSpec: raidcli.LogDriveSpec{
						Label:       "grp:2",
						RAIDLv:      raidcli.RAID01,
						Size:        238998823895, // 222.585*1024*1024*1024
						State:       "Optimal",
						StripSize:   262144,
						CachePolicy: "WriteThrough",
					},
				},
			},
		},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			got, err := parseLDList(string(in))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
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

func TestParseLogID(t *testing.T) {
	var grp, tgt int
	var err error
	tests := []struct {
		in  string
		grp int
		tgt int
		err error
	}{
		{"mox", 0, 0, errors.New("dummy")},
		{"1 (Target Id: 0)", 1, 0, nil},
		{"1 (Target Id: 1)", 1, 1, nil},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			grp, tgt, err = parseLogID(tt.in)
			if tt.err == nil && err != nil {
				t.Errorf("error should be nil, got: %s", err)
			}

			if tt.err != nil && err == nil {
				t.Errorf("error should NOT be nil")
			}

			if grp != tt.grp || tgt != tt.tgt {
				t.Errorf("test: %+v, got: grp=%d, tgt=%d, expect: grp=%d, tgt=%d", tt, grp, tgt, tt.grp, tt.tgt)
			}
		})
	}
}

func TestParseCachePolicy(t *testing.T) {
	var got string
	tests := []struct {
		in string
		ex string
	}{
		{"", "unknown"},
		{"WriteBack, ReadAhead, Direct, No Write Cache if Bad BBU", "WriteBack"},
		{"WriteThrough, ReadAhead, Direct, No Write Cache if Bad BBU", "WriteThrough"},
		{"mox", "unknown"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseCachePolicy(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}
