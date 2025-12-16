package state

import (
	"bytes"
	"fmt"
	"io"
)

type Address uint16

type Memory struct {
	// Some memory ranges will not be used in practice due to MMIO,
	// but it is easier to allocate the whole flat range.
	data [MemorySize]byte
}

func (m *Memory) ReadB(address Address) (byte, error) {
	return m.data[address], nil
}

func (m *Memory) WriteB(address Address, value byte) error {
	m.data[address] = value
	return nil
}

func (m *Memory) ReadH(address Address) (uint16, error) {
	return uint16(m.data[address]) | uint16(m.data[address+1])<<8, nil
}

func (m *Memory) WriteH(address Address, value uint16) error {
	m.data[address] = byte(value)
	m.data[address+1] = byte(value >> 8)
	return nil
}

func (m *Memory) Dump(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	const bytesPerLine = 16
	const halfLine = bytesPerLine / 2

	var zeroLine [16]byte
	var numEmpty int
	for i := 0; i < len(m.data)/bytesPerLine; i++ {
		baseAddress := i * bytesPerLine
		line := m.data[baseAddress : baseAddress+bytesPerLine]

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
			_, _ = fmt.Fprintf(w, "%02x ", m.data[baseAddress+i])
		}

		_, _ = fmt.Fprint(w, " ")

		for i := halfLine; i < bytesPerLine; i++ {
			_, _ = fmt.Fprintf(w, "%02x ", m.data[baseAddress+i])
		}

		_, _ = fmt.Fprint(w, " |")

		for i := 0; i < bytesPerLine; i++ {
			v := m.data[baseAddress+i]
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

const MemorySize = 64 * 1024
