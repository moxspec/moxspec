package sas3ircu

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/moxspec/moxspec/raidcli"
	"github.com/kylelemons/godebug/pretty"
)

func TestParseCtlListSingle(t *testing.T) {
	tests := []struct {
		path       string
		ex         []raidcli.ControllerSpec
		wantsError bool
	}{
		{"list.SAS3008.single.input", []raidcli.ControllerSpec{
			raidcli.ControllerSpec{
				Number:   0,
				Domain:   0,
				Bus:      0x61,
				Device:   0,
				Function: 0,
			},
		}, false},
		{"list.SAS3008.twin.input", []raidcli.ControllerSpec{
			raidcli.ControllerSpec{
				Number:   0,
				Domain:   0,
				Bus:      0x83,
				Device:   0,
				Function: 0,
			},
			raidcli.ControllerSpec{
				Number:   1,
				Domain:   0,
				Bus:      0x84,
				Device:   0,
				Function: 0,
			},
		}, false},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			got, err := parseCtlList(string(in))

			if tt.wantsError && err == nil {
				t.Errorf("%s wants error but got nil", tt.path)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(got) != len(tt.ex) {
				t.Errorf("%s len(got):%d, len(ex):%d", tt.path, len(got), len(tt.ex))
			}

			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i].ControllerSpec, tt.ex[i]) {
					t.Errorf(pretty.Compare(got[i].ControllerSpec, tt.ex[i]))
				}
			}
		})
	}
}

func TestParsePCIAddr(t *testing.T) {
	var dom, bus, dev, fun uint32
	var err error

	tests := []struct {
		in         string
		dom        uint32
		bus        uint32
		dev        uint32
		fun        uint32
		wantsError bool
	}{
		{"248", 0, 0, 0, 0, true},
		{"00h:61h:00h:00h", 0, 0x61, 0, 0, false},
		{"00h:3bh:00h:00h", 0, 0x3B, 0, 0, false},
		{"01h:0bh:01h:02h", 0x01, 0x0B, 0x01, 0x02, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.in, func(t *testing.T) {
			dom, bus, dev, fun, err = parsePCIAddr(tt.in)
			if tt.wantsError && err == nil {
				t.Errorf("%s wants error but got nil", tt.in)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("%s %s", tt.in, err)
			}

			if dom != tt.dom || bus != tt.bus || dev != tt.dev || fun != tt.fun {
				t.Errorf("test: %+v, got: dom=%d, bus=%d, dev=%d, fun=%d expect: dom=%d, bus=%d, dev=%d, fun=%d",
					tt, dom, bus, dev, fun, tt.dom, tt.bus, tt.dev, tt.fun)
			}
		})
	}
}
