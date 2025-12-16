// Package state implements the state of the r16 machine.
//
// The register file is its own package to enforce the invariant that the
// register zero is always zero.
package state

import (
	"fmt"
	"io"

	"github.com/jespert/primordial/internal/quality/assert"
)

type Registers struct {
	values [NumRegisters - 1]int16
}

func (r *Registers) Read(register int) int16 {
	r.assertValidRegister(register)

	if register == 0 {
		return 0
	}

	return r.values[register-1]
}

func (r *Registers) Write(register int, value int16) {
	r.assertValidRegister(register)

	if register != 0 {
		r.values[register-1] = value
	}
}

func (r *Registers) assertValidRegister(register int) {
	inBounds := register >= 0 && register < NumRegisters
	assert.Truef(inBounds, "register %d is out of bounds", register)
}

func (r *Registers) DumpNonZero(w io.Writer) {
	// There is nothing we can do on IO failure, so we just ignore errors.
	allZero := true
	for i := range NumRegisters {
		v := r.Read(i)
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

const NumRegisters = 16
