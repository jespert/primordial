package state_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/state"
	"github.com/jespert/primordial/internal/quality/approval"
	"github.com/jespert/primordial/internal/quality/require"
)

func TestMemory_B(t *testing.T) {
	var memory state.Memory
	for i := range state.MemorySize {
		address := state.Address(i)

		// Memory is initially zero initialised.
		original, err := memory.ReadB(address)
		require.Success(t, err)
		require.Equal(t, 0, original)

		// We assign a new value to the memory.
		written := byte(i)
		require.Success(t, memory.WriteB(address, written))

		// And we read it back. It should be the one we wrote.
		actual, err := memory.ReadB(address)
		require.Success(t, err)
		require.Equal(t, written, actual)
	}

	// Verify the final state of the memory.
	verifier := approval.NewTextVerifier(t)
	memory.Dump(verifier.Writer())
	verifier.Verify()
}

func TestMemory_H(t *testing.T) {
	var memory state.Memory
	for i := 0; i < state.MemorySize; i += 2 {
		address := state.Address(i)

		// Memory is initially zero initialised.
		original, err := memory.ReadH(address)
		require.Success(t, err)
		require.Equal(t, 0, original)

		// We assign a new value to the memory.
		written := int16(i)
		require.Success(t, memory.WriteH(address, written))

		// And we read it back. It should be the one we wrote.
		actual, err := memory.ReadH(address)
		require.Success(t, err)
		require.Equal(t, written, actual)
	}

	// Verify the final state of the memory.
	verifier := approval.NewTextVerifier(t)
	memory.Dump(verifier.Writer())
	verifier.Verify()
}

func TestMemory_W(t *testing.T) {
	var memory state.Memory
	for i := 0; i < state.MemorySize; i += 4 {
		address := state.Address(i)

		// Memory is initially zero initialised.
		original, err := memory.ReadW(address)
		require.Success(t, err)
		require.Equal(t, 0, original)

		// We assign a new value to the memory.
		written := int32(i)
		require.Success(t, memory.WriteW(address, written))

		// And we read it back. It should be the one we wrote.
		actual, err := memory.ReadW(address)
		require.Success(t, err)
		require.Equal(t, written, actual)
	}

	// Verify the final state of the memory.
	verifier := approval.NewTextVerifier(t)
	memory.Dump(verifier.Writer())
	verifier.Verify()
}

func TestMemory_Dump_initial(t *testing.T) {
	// Initially, the memory is empty.
	var memory state.Memory
	verifier := approval.NewTextVerifier(t)
	memory.Dump(verifier.Writer())
	verifier.Verify()
}

func TestMemory_Dump_empty_lines(t *testing.T) {
	const lineSize = 16

	var memory state.Memory
	require.Success(t, memory.WriteB(0, 1))
	require.Success(t, memory.WriteB(2*lineSize, 2))
	require.Success(t, memory.WriteB(5*lineSize, 3))
	require.Success(t, memory.WriteB(9*lineSize, 3))

	verifier := approval.NewTextVerifier(t)
	memory.Dump(verifier.Writer())
	verifier.Verify()
}
