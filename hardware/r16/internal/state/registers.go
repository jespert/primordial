// Package state implements the state of the r16 machine.
//
// The register file is its own package to enforce the invariant that the
// register zero is always zero.
package state

import "github.com/jespert/primordial/internal/quality/assert"

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

const NumRegisters = 16
