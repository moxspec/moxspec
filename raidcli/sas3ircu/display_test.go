package sas3ircu

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/actapio/moxspec/raidcli"
	"github.com/kylelemons/godebug/pretty"
)

func TestSplitSections(t *testing.T) {
	tests := []struct {
		path  string
		ctlen int
		ldlen int
		pdlen int
	}{
		{"display.SAS3008.input", 15, 8, 12},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			ctLines, ldLines, pdLines, err := splitSections(string(in))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(ctLines) != tt.ctlen {
				t.Errorf("%s len(ctLines): %d, expect: %d", tt.path, len(ctLines), tt.ctlen)
			}

			if len(ldLines) != tt.ldlen {
				t.Errorf("%s len(ldLines): %d, expect: %d", tt.path, len(ldLines), tt.ldlen)
			}

			if len(pdLines) != tt.pdlen {
				t.Errorf("%s len(pdLines): %d, expect: %d", tt.path, len(pdLines), tt.pdlen)
			}
		})
	}
}

func TestParseCTLines(t *testing.T) {
	tests := []struct {
		path       string
		firm       string
		bios       string
		wantsError bool
	}{
		{"display.SAS3008.ctlines.input", "6.00.00.00", "8.13.00.00", false},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			firm, bios, err := parseCTLines(strings.Split(string(in), "\n"))

			if tt.wantsError && err == nil {
				t.Errorf("%s wants error but got nil", tt.path)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if firm != tt.firm {
				t.Errorf("%s got: %s expect: %s", tt.path, firm, tt.firm)
			}

			if bios != tt.bios {
				t.Errorf("%s got: %s expect: %s", tt.path, bios, tt.bios)
			}
		})
	}
}

func TestParseLDLines(t *testing.T) {
	tests := []struct {
		path       string
		ex         []*LogDrive
		wantsError bool
	}{
		{"display.SAS3008.ldlines.2array.input", []*LogDrive{
			&LogDrive{
				VolumeID: 322,
				WWID:     "07f652e74f1fbcfd",
				LogDriveSpec: raidcli.LogDriveSpec{
					Label:  "vol:322",
					RAIDLv: raidcli.RAID0,
					Size:   1595999780864, // 1522064 * 1024 * 1024
					State:  "Inactive, Okay (OKY)",
				},
				pdIDList: []string{
					"0:0",
					"0:0",
					"1:1",
					"0:0",
				},
			},
			&LogDrive{
				VolumeID: 323,
				WWID:     "038476c328bd908d",
				LogDriveSpec: raidcli.LogDriveSpec{
					Label:  "vol:323",
					RAIDLv: raidcli.RAID0,
					Size:   1595999780864, // 1522064 * 1024 * 1024
					State:  "Failed (FLD)",
				},
				pdIDList: []string{
					"1:0",
					"0:0",
					"1:2",
					"1:4",
				},
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

			got, err := parseLDLines(strings.Split(string(in), "\n"))

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
				if !reflect.DeepEqual(got[i], tt.ex[i]) {
					t.Errorf(pretty.Compare(got[i], tt.ex[i]))
				}
			}
		})
	}
}

func TestParsePDLines(t *testing.T) {
	tests := []struct {
		path       string
		ex         []*PhyDrive
		wantsError bool
	}{
		{"display.SAS3008.pdlines.4disk_2fail.input", []*PhyDrive{
			&PhyDrive{
				EnclosureID:     "1",
				SlotNumber:      "0",
				SASAddress:      "4433221-1-0000-0000",
				Protocol:        "SATA",
				Size:            400088367104, // 381554 * 1024 * 1024
				Model:           "INTEL SSDSC2BX40",
				SerialNumber:    "BTHC528202X0400VGN",
				Firmware:        "0110",
				State:           "Optimal (OPT)",
				DriveType:       "SATA_SSD",
				SolidStateDrive: true,
			},
			&PhyDrive{
				EnclosureID:     "1",
				SlotNumber:      "4",
				SASAddress:      "4433221-1-0400-0000",
				Protocol:        "SATA",
				Size:            400088367104, // 381554 * 1024 * 1024
				Model:           "INTEL SSDSC2BX40",
				SerialNumber:    "BTHC52820151400VGN",
				Firmware:        "0110",
				State:           "Optimal (OPT)",
				DriveType:       "SATA_SSD",
				SolidStateDrive: true,
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

			got, err := parsePDLines(strings.Split(string(in), "\n"))

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
				if !reflect.DeepEqual(got[i], tt.ex[i]) {
					t.Errorf(pretty.Compare(got[i], tt.ex[i]))
				}
			}
		})
	}
}

func TestFormatRAIDLv(t *testing.T) {
	var got raidcli.Level

	tests := []struct {
		in string
		ex raidcli.Level
	}{
		{"248", raidcli.Unknown},
		{"RAID0", raidcli.RAID0},
		{"RAID1", raidcli.RAID1},
		{"RAID1E", raidcli.RAID1},
		{"RAID10", raidcli.RAID10},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = formatRAIDLv(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestIsValidPD(t *testing.T) {
	var got bool

	tests := []struct {
		in *PhyDrive
		ex bool
	}{
		{nil, false},
		{&PhyDrive{
			DriveType: "Enclosure services device",
		}, false},
		{&PhyDrive{
			DriveType: "Undetermined",
		}, false},
		{&PhyDrive{
			Size: 0,
		}, false},
		{&PhyDrive{
			State: "Missing (MIS)",
		}, false},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = isValidPD(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %t expect: %t", tt, got, tt.ex)
			}
		})
	}
}

func TestParseSizeIn(t *testing.T) {
	var got uint64

	tests := []struct {
		key string
		val string
		ex  uint64
	}{
		{"Size (in MB)", "248", 260046848},
		{"248", "248", 0},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseSizeIn(tt.key, tt.val)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestParseMBSize(t *testing.T) {
	var got uint64

	tests := []struct {
		in string
		ex uint64
	}{
		{"248", 0},
		{"0/248", 0},
		{"248/248", 260046848},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = parseMBSize(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %d expect: %d", tt, got, tt.ex)
			}
		})
	}
}

func TestGenPDMap(t *testing.T) {
	pd := &PhyDrive{
		EnclosureID:     "1",
		SlotNumber:      "4",
		SASAddress:      "4433221-1-0400-0000",
		Protocol:        "SATA",
		Size:            400088367104, // 381554 * 1024 * 1024
		Model:           "INTEL SSDSC2BX40",
		SerialNumber:    "BTHC52820151400VGN",
		Firmware:        "0110",
		State:           "Optimal (OPT)",
		DriveType:       "SATA_SSD",
		SolidStateDrive: true,
	}

	in := []*PhyDrive{
		nil,
		pd,
		nil,
	}

	ex := map[string]*PhyDrive{
		"1:4": pd,
	}

	got, err := genPDMap(in)
	if err != nil {
		t.Errorf("%s", err)
	}

	if !reflect.DeepEqual(got, ex) {
		t.Errorf(pretty.Compare(got, ex))
	}
}
