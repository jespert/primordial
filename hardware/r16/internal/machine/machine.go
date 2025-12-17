// Package machine represents an r16 machine.
package machine

import (
	"fmt"
	"io"

	"github.com/jespert/primordial/hardware/r16/internal/isa"
	"github.com/jespert/primordial/hardware/r16/internal/state"
)

// Machine of the machine (registers and memory).
type Machine struct {
	memory    state.Memory
	registers state.Registers

	// Instruction pointer.
	ip state.Address
}

// New creates a new Machine.
func New() *Machine {
	return &Machine{
		ip: ProgramBase,
	}
}

// Dump the state in human-friendly string representation to the given writer.
func (m *Machine) Dump(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	_, _ = fmt.Fprintf(w, "IP: 0x%04x\n", m.ip)
	_, _ = fmt.Fprint(w, "\nNon-zero registers:\n")
	m.registers.DumpNonZero(w)

	_, _ = fmt.Fprint(w, "\nMemory:\n")
	m.memory.Dump(w)
}

func (m *Machine) Step() error {
	encodedInstruction, err := m.fetchNextInstruction()
	if err != nil {
		return err
	}

	nextIP := m.ip + 2
	instruction := isa.Decode(encodedInstruction)

	switch instruction.Operation {
	case isa.JAL:
		returnPointer := nextIP
		x := m.registers.Read(instruction.X)
		nextIP = state.Address(x) + state.Address(instruction.Imm)
		m.registers.Write(instruction.Z, uint16(returnPointer))
	}

	m.ip = nextIP
	return nil
}

func (m *Machine) LoadProgram(base state.Address, data []byte) error {
	if len(data)+int(base) > state.MemorySize {
		return fmt.Errorf("program too large for memory: %d bytes", len(data))
	}

	m.memory.WriteRaw(ProgramBase, data)
	return nil
}

func (m *Machine) fetchNextInstruction() (isa.EncodedInstruction, error) {
	if m.ip%2 != 0 {
		return 0, fmt.Errorf("unaligned IP: %04x", m.ip)
	}

	v, err := m.memory.ReadW(m.ip)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch instruction at %04x: %w", m.ip, err)
	}

	return isa.EncodedInstruction(v), nil
}

const ProgramBase = 0x8000
