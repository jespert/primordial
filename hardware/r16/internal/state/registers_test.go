package state_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/state"
	"github.com/jespert/primordial/internal/quality/approval"
	"github.com/jespert/primordial/internal/quality/expect"
)

func TestRegisters_Write_to_zero_is_hardcoded(t *testing.T) {
	var file state.Registers
	file.Write(0, 1)
	expect.Equal(t, 0, file.Read(0))
}

func TestRegisters_Write_and_read_from_general_register(t *testing.T) {
	var file state.Registers
	for i := 1; i < state.NumRegisters; i++ {
		file.Write(i, int16(i))
		expect.Equal(t, int16(i), file.Read(i))
	}
}

func TestRegisters_Read_out_of_bounds_above(t *testing.T) {
	// The register number is too high.
	var file state.Registers
	expect.Panic(t, func() { file.Read(state.NumRegisters) })
}

func TestRegisters_Read_out_of_bounds_under(t *testing.T) {
	// The register number is too low.
	var file state.Registers
	expect.Panic(t, func() { file.Read(-1) })
}

func TestRegisters_DumpNonZero_initial(t *testing.T) {
	// All general-purpose registers are zero initially.
	var registers state.Registers
	verifier := approval.NewTextVerifier(t)
	registers.DumpNonZero(verifier.Writer())
	verifier.Verify()
}

func TestRegisters_DumpNonZero_all_non_zero(t *testing.T) {
	// Set all but ZR to a non-zero value.
	var registers state.Registers
	for i := 1; i < state.NumRegisters; i++ {
		registers.Write(i, int16(i))
	}

	verifier := approval.NewTextVerifier(t)
	registers.DumpNonZero(verifier.Writer())
	verifier.Verify()
}
