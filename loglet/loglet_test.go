package loglet

import (
	"bytes"
	"testing"
)

func TestOutput(t *testing.T) {
	out := new(bytes.Buffer)

	SetLevel(DEBUG)
	SetOutput(out)
	log := NewLogger("mox")
	log.Debug("mox")

	ex := "[debug][mox] mox\n"
	got := out.String()

	if got != ex {
		t.Errorf("got: %s, expect: %s", got, ex)
	}
}
