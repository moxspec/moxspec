package util

import (
	"fmt"
	"testing"
)

func TestSanitizeString(t *testing.T) {
	var got string
	tests := []struct {
		in string
		ex string
	}{
		{"248 (TM)", "248"},
		{"248 (R)", "248"},
		{"\x00mox\x00", "mox"},
		{"To be filled by O.E.M. mox ", "mox"},
		{"To Be Filled By O.E.M. mox", "mox"},
		{"248", "248"},
		{"(R) mox (TM)", "mox"},
		{"(TM) mox To Be Filled By O.E.M.", "mox"},
		{"", ""},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = SanitizeString(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestShortenVendorName(t *testing.T) {
	var got string
	tests := []struct {
		in string
		ex string
	}{
		{"INSYDE Corp.", "INSYDE"},
		{"Advanced Micro Devices, Inc. [AMD]", "AMD"},
		{"Broadcom Limited", "Broadcom"},
		{"Broadcom and subsidiaries", "Broadcom"},
		{"Broadcom Inc. and subsidiaries", "Broadcom"},
		{"LSI Logic / Symbios Logic", "LSI/Symbios"},
		{"Intel(R) Corporation", "Intel"},
		{"Intel Corporation", "Intel"},
		{"FUJITSU // American Megatrends Inc.", "FUJITSU // AMI"},
		{"Samsung Electronics Co Ltd", "Samsung"},
		{"Matrox Electronics Systems Ltd.", "Matrox"},
		{"Red Hat, Inc", "Red Hat"},
		{"NVIDIA Corporation", "NVIDIA"},
		{"ASPEED Technology, Inc.", "ASPEED"},
		{"Hewlett-Packard Company", "HP"},
		{"Mellanox Technologies", "Mellanox"},
		{"Quanta Computer Inc", "Quanta"},
		{"Hynix Semiconductor", "Hynix"},
		{"PenguinComputing", "Penguin Computing"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = ShortenVendorName(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}

func TestProcName(t *testing.T) {
	var got string
	tests := []struct {
		in string
		ex string
	}{
		{"Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz", "Intel Xeon Gold 6152 2.10GHz"},
		{"Intel(R) Xeon(R) CPU E5-2683 v4 @ 2.10GHz", "Intel Xeon E5-2683 v4 2.10GHz"},
		{"AMD EPYC 7601 32-Core Processor", "AMD EPYC 7601"},
		{"Intel(R) Xeon(R) CPU           L5630  @ 2.13GHz", "Intel Xeon L5630 2.13GHz"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got = ShortenProcName(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}
}
