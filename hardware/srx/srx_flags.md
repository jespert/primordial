# SRX Computer (flags)

**Warning: this highly experimental and very much a work in progress.**

This is the version of SRX with flags and separate comparison instructions
(like most RISC architectures other than MIPS and RISC-V).
This design is a work in progress and may be ultimately discarded.

## Flags

This design uses the standard NZCV set of flags.
To reduce excessive data dependencies,
only specific operations can read or write flags (similar to ARM64).

## Instruction formats

Instructions can be 16, 32 or 48-bits long.

Immediates can be 4, 8, 16 or 32-bits-long.
This consumes many opcodes but allows a very natural set of instructions.
Compared to RISC-V, it supports much further branches, jumps and calls.

16-bit format: (C for compact)

| Format | Mnemonic            | [12..15] | [8..11] | [4..7]  | [2..3] | [0..1] |
|--------|---------------------|----------|---------|---------|--------|--------|
| C1     | One argument        | Func8    | Func8   | AC/XZ   | Opcode | 0      |
| C2     | Two arguments       | Func4    | B/Y     | AC/XZ   | Opcode | 0      |
| CE     | Eight-bit immediate | imm8     | imm8    | AC/XZ   | Opcode | 0      |
| CF     | Four-bit immediate  | Func4    | imm4    | AC/XZ   | Opcode | 0      |

32-bit format:

| Format | Mnemonic   | [28..31] | [24..27] | [20..23] | [16..19] | [12..15] | [8..11]     | [2..7] | [0..1] |
|--------|------------|----------|----------|----------|----------|----------|-------------|--------|--------|
| R      | Register   | Func8    | Func8    | W        | B/Y      | C/Z      | A/X         | Opcode | 1      |
| E      | Eight      | Func8    | Func8    | imm8     | imm8     | C/Z      | A/X         | Opcode | 1      |
| A      | Assignment | imm16    | imm16    | imm16    | imm16    | C/Z      | A/X         | Opcode | 2      |
| B      | Branch     | imm16    | imm16    | imm16    | imm16    | Cond     | imm20/Func4 | Opcode | 2      |
| S      | Store      | imm16    | imm16    | imm16    | B/Y      | imm16    | A           | Opcode | 2      |

48-bit format: (X for extended)

| Format | Mnemonic   | [20..47] | [16..19] | [12..15]    | [8..11] | [2..7] | [0..1] |
|--------|------------|----------|----------|-------------|---------|--------|--------|
| XA     | Assignment | imm32    | imm32    | Z           | A/X     | Opcode | 3      |
| XB     | Branch     | imm32    | imm32    | imm36/Func4 | Cond    | Opcode | 3      |

## Branches

CMP uses a 16-bit instruction for both integers and pointers.

The branch parts are more interesting. We have two options:

| Type           | Range      | [16..31] | [12..15] | [8..11] | [2..7] | [0..1] |
|----------------|------------|----------|----------|---------|--------|--------|
| Base register  | A ± 64 kiB | imm16    | Cond     | A       | Opcode | 1      |
| Longer operand | IP ± 1 MiB | imm20    | Cond     | imm20   | Opcode | 1      |

There would be also a 48-bit variant for far branches:

| Type           | Range       | [16..47] | [12..15] | [8..11] | [2..7] | [0..1] |
|----------------|-------------|----------|----------|---------|--------|--------|
| Base register  | A ± 4 GiB   | imm32    | Cond     | A       | Opcode | 1      |
| Longer operand | IP ± 64 GiB | imm36    | Cond     | imm20   | Opcode | 1      |

And a 16-bit variant for near branches:

| Type             | Range      | [12..15] | [8..11] | [4..7] | [2..3] | [0..1] |
|------------------|------------|----------|---------|--------|--------|--------|
| Offset           | IP ± 256 B | imm8     | imm8    | Cond   | Opcode | 1      |
| Register         | A          | Cond     | Func4   | A      | Opcode | 1      |
| Register Indexed | A ± 16 B   | Cond     | imm4    | A      | Opcode | 1      |

The 16-bit register indexed variant is probably quite useless, so we might as
well use the other two and free up encoding space.

Having a register-based offset in conditional branches permits the use of
LUI/AIUPAC tricks to jump to any 64-bit address with an additional instruction.
However, we don't expect this to become common any time soon.

While in a way, having four extra bits is tempting, it does add complexity and
besmirches the symmetry of power-of-two immediate values.
