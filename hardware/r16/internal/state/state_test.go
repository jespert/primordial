package state_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/approval"
	"github.com/jespert/primordial/hardware/r16/internal/state"
)

func TestState_Dump_empty(t *testing.T) {
	verifier := approval.NewTextVerifier(t)
	m := state.New()
	m.Dump(verifier.Writer())
	verifier.Verify()
}
