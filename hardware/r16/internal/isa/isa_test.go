package isa_test

import (
	"testing"

	"github.com/jespert/primordial/hardware/r16/internal/isa"
	"github.com/jespert/primordial/internal/quality/expect"
)

func TestDecode(t *testing.T) {
	for _, tc := range encodingTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a bespoke test to print failures in hexadecimal.
			actual := isa.Decode(tc.encoded)
			if tc.decoded != actual {
				t.Errorf("Expected 0x%08x, got 0x%08x", tc.encoded, actual)
			}

			// Check reversibility.
			reversed := isa.Encode(actual)
			expect.Equal(t, tc.encoded, reversed)
		})
	}
}

func TestEncode(t *testing.T) {
	for _, tc := range encodingTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a bespoke test to print failures in hexadecimal.
			actual := isa.Encode(tc.decoded)
			if tc.encoded != actual {
				t.Errorf("Expected 0x%08x, got 0x%08x", tc.encoded, actual)
			}

			// Check reversibility.
			reversed := isa.Decode(actual)
			expect.Equal(t, tc.decoded, reversed)
		})
	}
}

var encodingTestCases = []struct {
	name    string
	decoded isa.DecodedInstruction
	encoded isa.EncodedInstruction
}{
	{
		name: "R",
		decoded: isa.DecodedInstruction{
			Operation: 0x0123,
			Z:         0xa,
			Y:         0xb,
			X:         0xc,
			W:         0xd,
		},
		encoded: 0x0abcd123,
	},
	{
		name: "B",
		decoded: isa.DecodedInstruction{
			Operation: 0x4003,
			Y:         0xb,
			X:         0xc,
			Imm:       0x6789,
		},
		encoded: 0x43bc6789,
	},
	{
		name: "A",
		decoded: isa.DecodedInstruction{
			Operation: 0x8003,
			Z:         0xa,
			X:         0xc,
			Imm:       0x6789,
		},
		encoded: 0x8a3c6789,
	},
	{
		name: "jal",
		decoded: isa.DecodedInstruction{
			Operation: isa.JAL,
			Z:         isa.RP,
			X:         isa.A2,
			Imm:       0x1234,
		},
		encoded: 0x8e1c1234,
	},
}
