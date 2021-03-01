package hpacucli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/moxspec/moxspec/raidcli"
	"github.com/kylelemons/godebug/pretty"
)

func TestSplitConfigDetailSections(t *testing.T) {
	tests := []struct {
		path    string
		ctlen   int
		ldpdlen int
	}{
		{"show_config_detail.sa.input", 77, 125},
		{"show_config_detail.dynsa.input", 38, 63},
		{"show_config_detail.hba.input", 107, 329},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			ctLines, ldpdLines, err := splitConfigDetailSections(string(in))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(ctLines) != tt.ctlen {
				t.Errorf("%s len(ctLines): %d, expect: %d", tt.path, len(ctLines), tt.ctlen)
			}

			if len(ldpdLines) != tt.ldpdlen {
				t.Errorf("%s len(ldpdLines): %d, expect: %d", tt.path, len(ldpdLines), tt.ldpdlen)
			}
		})
	}
}

func TestParseCTLines(t *testing.T) {
	tests := []struct {
		path    string
		sn      string
		firm    string
		battery bool
		pciaddr string
	}{
		{"show_config_detail.ctlines.sa.input", "PDSXK0ARH5O18K", "6.68", true, "0000:04:00.0"},
		{"show_config_detail.ctlines.hba.input", "PDNLN0BRH8O054", "6.60", false, ""},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			lines := strings.Split(string(in), "\n")
			sn, firm, battery, pciaddr, err := parseCTLines(lines)
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if sn != tt.sn {
				t.Errorf("%s sn: got: %s expect: %s", tt.path, sn, tt.sn)
			}

			if firm != tt.firm {
				t.Errorf("%s firmware: got: %s expect: %s", tt.path, firm, tt.firm)
			}

			if battery != tt.battery {
				t.Errorf("%s battery: got: %t expect: %t", tt.path, battery, tt.battery)
			}

			if pciaddr != tt.pciaddr {
				t.Errorf("%s pciaddr: got: %s expect: %s", tt.path, pciaddr, tt.pciaddr)
			}
		})
	}
}

func TestSplitArrays(t *testing.T) {
	tests := []struct {
		path   string
		arrlen int
		unlen  int
	}{
		{"show_config_detail.1array_unconf.sa.input", 1, 39},
		{"show_config_detail.2array.raid0_raid0.sa.input", 2, 0},
		{"show_config_detail.2array.raid1_raid5.sa.input", 2, 0},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			arrays, unassigned, err := splitArrays(strings.Split(string(in), "\n"))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(arrays) != tt.arrlen {
				t.Errorf("%s len(arrays): %d, expect: %d", tt.path, len(arrays), tt.arrlen)
			}

			if len(unassigned) != tt.unlen {
				t.Errorf("%s len(unassigned): %d, expect: %d", tt.path, len(unassigned), tt.unlen)
			}
		})
	}
}

func TestSplitLDChunks(t *testing.T) {
	tests := []struct {
		path   string
		ldslen int
	}{
		{"show_config_detail.1array_unconf.sa.input", 1},
		{"show_config_detail.2array.raid0_raid0.sa.input", 2},
		{"show_config_detail.2array.raid1_raid5.sa.input", 2},
	}

	for _, test := range tests {
		tt := test // bind loop var

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			lds, err := splitLDChunks(strings.Split(string(in), "\n"))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(lds) != tt.ldslen {
				t.Errorf("%s len(lds): %d, expect: %d", tt.path, len(lds), tt.ldslen)
			}
		})
	}
}

func TestSplitLDPDLines(t *testing.T) {
	tests := []struct {
		path  string
		ldlen int
		pdlen int
	}{
		{"show_config_detail.ldpd_raid10.sa.input", 23, 96},
		{"show_config_detail.ldpd_raid5.sa.input", 18, 81},
		{"show_config_detail.ldpd_raid1.sa.input", 20, 38},
	}

	for _, tt := range tests {
		in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
		if err != nil {
			t.Errorf("%s %s", tt.path, err)
		}

		ldLines, pdLines, err := splitLDPDSections(strings.Split(string(in), "\n"))
		if err != nil {
			t.Errorf("%s %s", tt.path, err)
		}

		if len(ldLines) != tt.ldlen {
			t.Errorf("%s len(ldLines): %d, expect: %d", tt.path, len(ldLines), tt.ldlen)
		}

		if len(pdLines) != tt.pdlen {
			t.Errorf("%s len(pdLines): %d, expect: %d", tt.path, len(pdLines), tt.pdlen)
		}
	}
}

func TestParseLDLines(t *testing.T) {
	tests := []struct {
		path string
		ex   *LogDrive
	}{
		{"show_config_detail.ld.sa.input", &LogDrive{
			VolumeID:  "1",
			StripSize: 262144, // 256 * 1024
			UUID:      "600508B1001C1BD85F342B5B6BDC3C2C",
			DiskName:  "/dev/sda",
			LogDriveSpec: raidcli.LogDriveSpec{
				Label:  "vol:1",
				RAIDLv: raidcli.RAID10,
				Size:   1600000000000,
				State:  "OK",
			},
			pdIDList: []string{
				"1I:1:1",
				"1I:1:2",
				"1I:1:3",
				"1I:1:4",
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

			got, err := parseLDLines(strings.Split(string(in), "\n"))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if !reflect.DeepEqual(got, tt.ex) {
				t.Errorf(fmt.Sprintf("%s %s", tt.path, pretty.Compare(got, tt.ex)))
			}
		})
	}
}

func TestParsePDLines(t *testing.T) {
	tests := []struct {
		path string
		ex   []*PhyDrive
	}{
		{"show_config_detail.pd.sa.input", []*PhyDrive{
			&PhyDrive{
				Port:         "1I",
				Box:          "1",
				Bay:          "1",
				Protocol:     "SAS",
				Model:        "HP EG0900FCSPN",
				SerialNumber: "Y3E0A0RMFTM11346",
				Size:         900100000000,
				Firmware:     "HPD0",
				Status:       "OK",
				Rotation:     10000,
				NegSpeed:     "6.0Gbps",
				CurTemp:      36,
				MaxTemp:      52,
			},
			&PhyDrive{
				Port:         "1I",
				Box:          "1",
				Bay:          "2",
				Protocol:     "SAS",
				Model:        "HP EG0900FCSPN",
				SerialNumber: "Y3E0A0XKFTM11346",
				Size:         900100000000,
				Firmware:     "HPD0",
				Status:       "OK",
				Rotation:     10000,
				NegSpeed:     "6.0Gbps",
				CurTemp:      38,
				MaxTemp:      56,
			},
			&PhyDrive{
				Port:         "1I",
				Box:          "1",
				Bay:          "3",
				Protocol:     "SAS",
				Model:        "HP EG0900FCSPN",
				SerialNumber: "Y3E0A11BFTM11346",
				Size:         900100000000,
				Firmware:     "HPD0",
				Status:       "OK",
				Rotation:     10000,
				NegSpeed:     "6.0Gbps",
				CurTemp:      35,
				MaxTemp:      51,
			},
			&PhyDrive{
				Port:         "1I",
				Box:          "1",
				Bay:          "4",
				Protocol:     "SAS",
				Model:        "HP EG0900FCSPN",
				SerialNumber: "Y3E0A0RGFTM11346",
				Size:         900100000000,
				Firmware:     "HPD0",
				Status:       "OK",
				Rotation:     10000,
				NegSpeed:     "6.0Gbps",
				CurTemp:      36,
				MaxTemp:      52,
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

			got, err := parsePDLines(strings.Split(string(in), "\n"))
			if err != nil {
				t.Errorf("%s", err)
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
