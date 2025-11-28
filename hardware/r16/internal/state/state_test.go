package state_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/state"
	"github.com/jespert/primordial/internal/quality/approval"
	"github.com/jespert/primordial/internal/quality/expect"
)

func TestState_Dump_empty(t *testing.T) {
	verifier := approval.NewTextVerifier(t)
	m := state.New()
	m.Dump(verifier.Writer())
	verifier.Verify()
}

func TestState_SetIP(t *testing.T) {
	m := state.New()
	const newIP = 0x8000
	m.SetIP(newIP)
	expect.Equal(t, newIP, m.IP())
}
