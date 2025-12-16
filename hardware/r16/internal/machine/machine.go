// Package machine represents an r16 machine.
package machine

import (
	"fmt"
	"io"

	"github.com/jespert/primordial/hardware/r16/internal/state"
)

// Machine of the machine (registers and memory).
type Machine struct {
	memory    state.Memory
	registers state.Registers

	// Instruction pointer.
	ip uint16
}

// New creates a new Machine.
func New() *Machine {
	return &Machine{}
}

// Dump the state in human-friendly string representation to the given writer.
func (m *Machine) Dump(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	_, _ = fmt.Fprint(w, "IP: ", m.ip)
	_, _ = fmt.Fprint(w, "\n\nNon-zero registers:\n")
	m.registers.DumpNonZero(w)

	_, _ = fmt.Fprint(w, "\nMemory:\n")
	m.memory.Dump(w)
}
