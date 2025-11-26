package state_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/state"
)

//go:embed testdata/state_dump_empty.expected.txt
var ExpectedStateDumpEmpty string

func TestState_Dump_empty(t *testing.T) {
	w := &bytes.Buffer{}
	m := state.New()
	m.Dump(w)
	verify(t, ExpectedStateDumpEmpty, w.String())
}

func verify(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("FAIL: Unexpected dump\n\n")
		t.Errorf("Expected:\n%s\nEOF\n\n", ExpectedStateDumpEmpty)
		t.Errorf("Actual:\n%s\nEOF\n\n", actual)
	}
}
