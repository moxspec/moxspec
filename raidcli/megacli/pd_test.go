package megacli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestParsePDList(t *testing.T) {
	tests := []struct {
		path string
		ex   []*PhyDrive
	}{
		{"pdlist.empty.input", []*PhyDrive{}},
		{"pdlist.2disk_1hdd_1ssd.input", []*PhyDrive{
			&PhyDrive{
				WWN:              "5000CCA06E3B2513",
				EnclosureID:      "252",
				SlotNumber:       "0",
				Group:            0,
				Span:             0,
				Arm:              0,
				Type:             "SAS",
				Model:            "HITACHI HUC109030CSS600 A440W5H1JAZG",
				InquiryRaw:       "HITACHI HUC109030CSS600 A440W5H1JAZG",
				Size:             299999170658, // 279.396 * 1024 * 1024 * 1024
				FirmwareRevision: "A440",
				State:            "Online, Spun Up",
				PhyBlockSize:     512,
				LogBlockSize:     512,
				ConnectedPort:    "0(path0)",
				DriveSpeed:       "6.0Gb/s",
				LinkSpeed:        "6.0Gb/s",
				CurTemp:          32,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  1,
				DeviceID:         0,
			},
			&PhyDrive{
				WWN:              "5000CCA06E3AB89B",
				EnclosureID:      "252",
				SlotNumber:       "1",
				Group:            0,
				Span:             0,
				Arm:              1,
				Type:             "SAS",
				Model:            "HITACHI HUC109030CSS600 A440W5H1941G",
				InquiryRaw:       "HITACHI HUC109030CSS600 A440W5H1941G",
				Size:             299999170658, // 279.396 * 1024 * 1024 * 1024
				FirmwareRevision: "A440",
				State:            "Online, Spun Up",
				PhyBlockSize:     512,
				LogBlockSize:     512,
				ConnectedPort:    "1(path0)",
				DriveSpeed:       "6.0Gb/s",
				LinkSpeed:        "6.0Gb/s",
				CurTemp:          34,
				SolidStateDrive:  true,
				SMARTAlert:       true,
				MediaErrorCount:  0,
				DeviceID:         1,
			}},
		},
		{"pdlist.4disk_hdd.input", []*PhyDrive{
			&PhyDrive{
				WWN:              "5000039508008535",
				EnclosureID:      "252",
				SlotNumber:       "0",
				Group:            0,
				Span:             0,
				Arm:              0,
				Type:             "SAS",
				Model:            "TOSHIBA MBF2300RC 5212EB07PD90AHV1",
				InquiryRaw:       "TOSHIBA MBF2300RC       5212EB07PD90AHV1",
				Size:             299999170658, // 279.396 * 1024 * 1024 * 1024
				FirmwareRevision: "5212",
				State:            "Online, Spun Up",
				PhyBlockSize:     0,
				LogBlockSize:     0,
				ConnectedPort:    "3(path0)",
				DriveSpeed:       "6.0Gb/s",
				LinkSpeed:        "6.0Gb/s",
				CurTemp:          30,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         19,
			},
			&PhyDrive{
				WWN:              "5000039508008701",
				EnclosureID:      "252",
				SlotNumber:       "1",
				Group:            0,
				Span:             0,
				Arm:              1,
				Type:             "SAS",
				Model:            "TOSHIBA MBF2300RC 5212EB07PD90AHWF",
				InquiryRaw:       "TOSHIBA MBF2300RC       5212EB07PD90AHWF",
				Size:             299999170658, // 279.396 * 1024 * 1024 * 1024
				FirmwareRevision: "5212",
				State:            "Online, Spun Up",
				PhyBlockSize:     0,
				LogBlockSize:     0,
				ConnectedPort:    "2(path0)",
				DriveSpeed:       "6.0Gb/s",
				LinkSpeed:        "6.0Gb/s",
				CurTemp:          30,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         18,
			},
			&PhyDrive{
				WWN:              "50000395080086B9",
				EnclosureID:      "252",
				SlotNumber:       "2",
				Group:            0,
				Span:             1,
				Arm:              0,
				Type:             "SAS",
				Model:            "TOSHIBA MBF2300RC 5212EB07PD90AHW9",
				InquiryRaw:       "TOSHIBA MBF2300RC       5212EB07PD90AHW9",
				Size:             299999170658, // 279.396 * 1024 * 1024 * 1024
				FirmwareRevision: "5212",
				State:            "Online, Spun Up",
				PhyBlockSize:     0,
				LogBlockSize:     0,
				ConnectedPort:    "1(path0)",
				DriveSpeed:       "6.0Gb/s",
				LinkSpeed:        "6.0Gb/s",
				CurTemp:          30,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         17,
			},
			&PhyDrive{
				WWN:              "5000039488201035",
				EnclosureID:      "252",
				SlotNumber:       "3",
				Group:            0,
				Span:             1,
				Arm:              1,
				Type:             "SAS",
				Model:            "TOSHIBA MBF2300RC 5212EB07PD108N2S",
				InquiryRaw:       "TOSHIBA MBF2300RC       5212EB07PD108N2S",
				Size:             299999170658, // 279.396 * 1024 * 1024 * 1024
				FirmwareRevision: "5212",
				State:            "Online, Spun Up",
				PhyBlockSize:     0,
				LogBlockSize:     0,
				ConnectedPort:    "0(path0)",
				DriveSpeed:       "6.0Gb/s",
				LinkSpeed:        "6.0Gb/s",
				CurTemp:          29,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         20,
			},
		}},
		{"pdlist.4disk_2hdd_2jbod.input", []*PhyDrive{
			&PhyDrive{
				WWN:              "5000CCA2731373AB",
				EnclosureID:      "32",
				SlotNumber:       "0",
				Group:            0,
				Span:             0,
				Arm:              0,
				Type:             "SAS",
				Model:            "HGST HUH721010AL5200 LS142YGAPMLD",
				InquiryRaw:       "HGST    HUH721010AL5200 LS142YGAPMLD",
				Size:             9796648603484, // 8.910 * 1024 * 1024 * 1024
				FirmwareRevision: "LS14",
				State:            "JBOD",
				PhyBlockSize:     4096,
				LogBlockSize:     512,
				ConnectedPort:    "0(path0)",
				DriveSpeed:       "12.0Gb/s",
				LinkSpeed:        "12.0Gb/s",
				CurTemp:          30,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         0,
			},
			&PhyDrive{
				WWN:              "5000CCA273155BB7",
				EnclosureID:      "32",
				SlotNumber:       "1",
				Group:            0,
				Span:             0,
				Arm:              0,
				Type:             "SAS",
				Model:            "HGST HUH721010AL5200 LS142YGBS45D",
				InquiryRaw:       "HGST    HUH721010AL5200 LS142YGBS45D",
				Size:             9796648603484, // 8.910 * 1024 * 1024 * 1024
				FirmwareRevision: "LS14",
				State:            "Unconfigured(good), Spun Up",
				PhyBlockSize:     4096,
				LogBlockSize:     512,
				ConnectedPort:    "0(path0)",
				DriveSpeed:       "12.0Gb/s",
				LinkSpeed:        "12.0Gb/s",
				CurTemp:          30,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         1,
			},
			&PhyDrive{
				WWN:              "5000CCA26ADC8FDB",
				EnclosureID:      "32",
				SlotNumber:       "2",
				Group:            0,
				Span:             0,
				Arm:              0,
				Type:             "SAS",
				Model:            "HGST HUH721010AL5200 LS142TKX9AJD",
				InquiryRaw:       "HGST    HUH721010AL5200 LS142TKX9AJD",
				Size:             9796648603484, // 8.910 * 1024 * 1024 * 1024
				FirmwareRevision: "LS14",
				State:            "Unconfigured(good), Spun Up",
				PhyBlockSize:     4096,
				LogBlockSize:     512,
				ConnectedPort:    "0(path0)",
				DriveSpeed:       "12.0Gb/s",
				LinkSpeed:        "12.0Gb/s",
				CurTemp:          29,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         2,
			},
			&PhyDrive{
				WWN:              "5000CCA26ADE95BB",
				EnclosureID:      "32",
				SlotNumber:       "3",
				Group:            0,
				Span:             0,
				Arm:              0,
				Type:             "SAS",
				Model:            "HGST HUH721010AL5200 LS142TKYDUPD",
				InquiryRaw:       "HGST    HUH721010AL5200 LS142TKYDUPD",
				Size:             9796648603484, // 8.910 * 1024 * 1024 * 1024
				FirmwareRevision: "LS14",
				State:            "JBOD",
				PhyBlockSize:     4096,
				LogBlockSize:     512,
				ConnectedPort:    "0(path0)",
				DriveSpeed:       "12.0Gb/s",
				LinkSpeed:        "12.0Gb/s",
				CurTemp:          30,
				SolidStateDrive:  false,
				SMARTAlert:       false,
				MediaErrorCount:  0,
				DeviceID:         3,
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

			got, err := parsePDList(strings.Split(string(in), "\n"))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if len(got) != len(tt.ex) {
				t.Fatalf("%s len(got):%d, len(ex):%d", tt.path, len(got), len(tt.ex))
			}

			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i], tt.ex[i]) {
					t.Errorf(pretty.Compare(got[i], tt.ex[i]))
				}
			}
		})
	}
}

func TestDrivePos(t *testing.T) {
	var grp, span, arm int
	var err error
	tests := []struct {
		in   string
		grp  int
		span int
		arm  int
		err  error
	}{
		{"Span: 1, Arm: 1", 0, 0, 0, errors.New("dummy")},
		{"DiskGroup: 0, Span: 1, Arm: 1", 0, 1, 1, nil},
		{"DiskGroup: 1, Span: 1, Arm: 1", 1, 1, 1, nil},
		{"DiskGroup: 2, Span: 1, Arm: 1", 2, 1, 1, nil},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			grp, span, arm, err = parseDrivePos(tt.in)
			if tt.err == nil && err != nil {
				t.Errorf("error should be nil, got: %s", err)
			}

			if tt.err != nil && err == nil {
				t.Errorf("error should NOT be nil")
			}

			if grp != tt.grp || span != tt.span || arm != tt.arm {
				t.Errorf("test: %+v, got: grp=%d, span=%d, arm=%d, expect: grp=%d, span=%d, arm=%d", tt, grp, span, arm, tt.grp, tt.span, tt.arm)
			}
		})
	}
}
