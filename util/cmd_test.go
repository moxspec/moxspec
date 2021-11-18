package util

import "testing"

func TestScanPathList(t *testing.T) {
	var err error

	_, err = ScanPathList([]string{
		"/sys/248/is/the/most/important/number/in/the/world/",
		"/proc/mox/is/pretty/cool/",
	})

	if err == nil {
		t.Errorf("error should be not nil")
	}

}
