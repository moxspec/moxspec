package nvidia

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestDecodeLog(t *testing.T) {
	tests := []struct {
		path     string
		exDrv    string
		exNumGPU int
	}{
		{"2gpu-m40.xml", "390.46", 2},
	}

	for _, test := range tests {
		tt := test

		t.Run(tt.path, func(t *testing.T) {
			in, err := ioutil.ReadFile(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			l, err := decodeLog(string(in))
			if err != nil {
				t.Errorf("%s %s", tt.path, err)
			}

			if l.DriverVersion != tt.exDrv {
				t.Errorf("%s driver_version got:%s ex:%s", tt.path, l.DriverVersion, tt.exDrv)
			}

			if len(l.GPUs) != tt.exNumGPU {
				t.Errorf("%s num_gpu got:%d ex:%d", tt.path, len(l.GPUs), tt.exNumGPU)
			}
		})
	}
}
