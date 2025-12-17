package machine

import (
	"bytes"
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/isa"
	"github.com/jespert/primordial/internal/quality/approval"
	"github.com/jespert/primordial/internal/quality/require"
)

func TestMachine_Dump_empty(t *testing.T) {
	m := New()
	verify(t, m)
}

func TestMachine_Step_JAL(t *testing.T) {
	m := withProgram(
		t,
		isa.DecodedInstruction{
			Operation: isa.JAL,
			Z:         isa.RP,
			X:         isa.A0,
			Imm:       0x1234,
		},
	)
	m.registers.Write(isa.A0, 0x9999)

	require.Success(t, m.Step())
	verify(t, m)
}

func verify(t *testing.T, m *Machine) {
	t.Helper()
	verifier := approval.NewTextVerifier(t)
	m.Dump(verifier.Writer())
	verifier.Verify()
}

func withProgram(t *testing.T, instructions ...isa.DecodedInstruction) *Machine {
	var buffer bytes.Buffer
	for _, instruction := range instructions {
		encoded := isa.Encode(instruction)
		buffer.WriteByte(byte(encoded))
		buffer.WriteByte(byte(encoded >> 8))
		buffer.WriteByte(byte(encoded >> 16))
		buffer.WriteByte(byte(encoded >> 24))
	}

	m := New()
	require.Success(t, m.LoadProgram(0, buffer.Bytes()))
	return m
}
