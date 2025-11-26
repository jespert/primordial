package state_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/approval"
	"github.com/jespert/primordial/hardware/r16/internal/state"
)

func TestState_Dump_empty(t *testing.T) {
	w := &bytes.Buffer{}
	m := state.New()
	m.Dump(w)
	approval.Verify(t, w.String())
}
