// Package isa implements the R16 instruction set.
package isa

import "github.com/jespert/primordial/internal/quality/assert"

// Register is a register number.
type Register uint8

const (
	ZR Register = 0
	S6 Register = 1
	S5 Register = 2
	S4 Register = 3
	S3 Register = 4
	S2 Register = 5
	S1 Register = 6
	S0 Register = 7
	T0 Register = 8
	T1 Register = 9
	A0 Register = 10
	A1 Register = 11
	A2 Register = 12
	A3 Register = 13
	RP Register = 14
	SP Register = 15
)

// Operation code (opcode + function).
type Operation uint16

const (
	JAL Operation = 0x8001
)

type DecodedInstruction struct {
	Operation Operation
	Z         Register
	Y         Register
	X         Register
	W         Register
	Imm       uint16
}

type EncodedInstruction uint32

// Decode instruction.
func Decode(e EncodedInstruction) DecodedInstruction {
	// The four MSBs of the encoded instruction will be the four
	// MSBs of the operation. The LSBs of the operation will be
	// filled with the function field, if any.
	opcode := Operation((e >> 28) << 12)

	// Precalculate fields that are common to at least two formats.
	z := Register(e>>offsetZ) & 0xf
	y := Register(e>>offsetY) & 0xf
	x := Register(e>>offsetX) & 0xf
	imm := uint16(e)

	// The two MSBs encode determine the instruction format.
	switch fmt := e >> 30; fmt {
	case 0:
		// R-type
		return DecodedInstruction{
			Operation: opcode | (Operation(e) & 0x0fff),
			Z:         z,
			Y:         y,
			X:         x,
			W:         Register(e>>offsetW) & 0xf,
			Imm:       0,
		}

	case 1:
		// B-type
		return DecodedInstruction{
			Operation: opcode | Operation(z),
			Z:         0,
			Y:         y,
			X:         x,
			W:         0,
			Imm:       imm,
		}

	default:
		// A-type
		return DecodedInstruction{
			Operation: opcode | Operation(y),
			Z:         z,
			Y:         0,
			X:         x,
			W:         0,
			Imm:       imm,
		}
	}
}

// Encode instruction.
func Encode(d DecodedInstruction) EncodedInstruction {
	// The four MSBs of the operation will be the four MSBs of the
	// encoded instruction. The LSBs of the operation will be filled
	// with the function field, if any.
	opcode := EncodedInstruction(d.Operation&0xf000) << 16

	// The OR operation provides correct results even when a field is
	// unused because it would be zero.
	z := EncodedInstruction(d.Z) << offsetZ
	y := EncodedInstruction(d.Y) << offsetY
	x := EncodedInstruction(d.X) << offsetX
	w := EncodedInstruction(d.W) << offsetW
	imm := EncodedInstruction(d.Imm)

	// Encode everything but the function field.
	allButFunction := opcode | z | y | x | w | imm

	// The two MSBs encode determine the instruction format.
	var function EncodedInstruction
	switch fmt := d.Operation >> 14; fmt {
	case 0:
		// R-type
		assert.Equal(0, imm) // TODO replace with Zero.
		function = EncodedInstruction(d.Operation & 0xfff)

	case 1:
		// B-type
		assert.Equal(0, z)                  // TODO replace with Zero.
		assert.Equal(0, d.Operation&0x0ff0) // TODO replace with Zero.
		function = EncodedInstruction(d.Operation&0xf) << offsetZ

	default:
		// A-type
		assert.Equal(0, y)                  // TODO replace with Zero.
		assert.Equal(0, d.Operation&0x0ff0) // TODO replace with Zero.
		function = EncodedInstruction(d.Operation&0xf) << offsetY
	}

	return allButFunction | function
}

const (
	offsetZ = 24
	offsetY = 20
	offsetX = 16
	offsetW = 12
)
