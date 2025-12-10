// Package state contains the state of the machine (registers and memory)
//
// The reason the State is in a separate package is to enforce invariants,
// the most important of which is not overwriting ZR.
package state

import (
	"bytes"
	"fmt"
	"io"
)

// State of the machine (registers and memory).
type State struct {
	// Some memory ranges will not be used in practice due to MMIO,
	// but it is easier to allocate the whole flat range.
	memory [64 * 1024]byte

	// The register at index 0 won't be used, and it is tempting to store
	// the instruction pointer there, but probably not worth the confusion.
	registers [16]int16

	// Instruction pointer.
	ip uint16
}

// New creates a new State.
func New() *State {
	return &State{}
}

// Dump the state in human-friendly string representation to the given writer.
func (m *State) Dump(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	_, _ = fmt.Fprint(w, "IP: ", m.ip)
	_, _ = fmt.Fprint(w, "\n\nNon-zero registers:\n")
	m.dumpNonZeroRegisters(w)

	_, _ = fmt.Fprint(w, "\nMemory:\n")
	m.dumpMemory(w)
}

func (m *State) dumpNonZeroRegisters(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	allZero := true
	for i, v := range m.registers {
		if v == 0 {
			continue
		}

		allZero = false
		_, _ = fmt.Fprintf(
			w,
			"%1X: 0x%04x S:%d U:%d\n",
			i,
			v,
			v,
			uint16(v),
		)
	}

	if allZero {
		_, _ = fmt.Fprint(w, "(none)\n")
	}
}

func (m *State) dumpMemory(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	const bytesPerLine = 16
	const halfLine = bytesPerLine / 2

	var zeroLine [16]byte
	var numEmpty int
	for i := 0; i < len(m.memory)/bytesPerLine; i++ {
		baseAddress := i * bytesPerLine
		line := m.memory[baseAddress : baseAddress+bytesPerLine]

		if bytes.Compare(line, zeroLine[:]) == 0 {
			numEmpty++
			continue
		} else if numEmpty != 0 {
			if numEmpty == 1 {
				_, _ = fmt.Fprint(w, "(1 empty line)\n")
			} else {
				_, _ = fmt.Fprintf(w, "(%d empty lines)\n", numEmpty)
			}
			numEmpty = 0
		}

		_, _ = fmt.Fprintf(w, "%04x  ", baseAddress)

		for i := 0; i < halfLine; i++ {
			_, _ = fmt.Fprintf(w, "%02x ", m.memory[baseAddress+i])
		}

		_, _ = fmt.Fprint(w, " ")

		for i := halfLine; i < bytesPerLine; i++ {
			_, _ = fmt.Fprintf(w, "%02x ", m.memory[baseAddress+i])
		}

		_, _ = fmt.Fprint(w, " |")

		for i := 0; i < bytesPerLine; i++ {
			v := m.memory[baseAddress+i]
			if v >= 32 && v <= 126 {
				_, _ = fmt.Fprintf(w, "%c", v)
			} else {
				_, _ = fmt.Fprint(w, ".")
			}
		}

		_, _ = fmt.Fprint(w, "|\n")
	}

	if numEmpty > 0 {
		_, _ = fmt.Fprintf(w, "(%d empty lines)\n", numEmpty)
		numEmpty = 0
	}
}
