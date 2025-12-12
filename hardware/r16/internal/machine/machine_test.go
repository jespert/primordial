package machine_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/machine"
	"github.com/jespert/primordial/internal/quality/approval"
)

func TestMachine_Dump_empty(t *testing.T) {
	verifier := approval.NewTextVerifier(t)
	m := machine.New()
	m.Dump(verifier.Writer())
	verifier.Verify()
}
