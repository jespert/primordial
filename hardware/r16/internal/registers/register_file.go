// Package registers implements the register file.
//
// The register file is its own package to enforce the invariant that the
// zero register is always zero.
package registers

import "github.com/jespert/primordial/internal/quality/assert"

type File struct {
	values [NumRegisters - 1]int16
}

func (r *File) Read(register int) int16 {
	r.assertValidRegister(register)

	if register == 0 {
		return 0
	}

	return r.values[register-1]
}

func (r *File) Write(register int, value int16) {
	r.assertValidRegister(register)

	if register != 0 {
		r.values[register-1] = value
	}
}

func (r *File) assertValidRegister(register int) {
	inBounds := register >= 0 && register < NumRegisters
	assert.Truef(inBounds, "register %d is out of bounds", register)
}

const NumRegisters = 16
