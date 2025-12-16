# SRX Computer (test-and-branch)

**Warning: this highly experimental and very much a work in progress.**

This is the version of SRX with test-and-branch instructions (like RISC-V).
This design is a work in progress and may be ultimately discarded.

Specification and emulator for the SRX computer architecture.
It supports 16, 32, 64 and 128-bit architectures.

SRX has separate data and address registers and variable-length instructions
(S = split register file, X = variable-length instructions).

## Instruction encoding

### Formats

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

| Format | Mnemonic   | [28..31] | [24..27] | [20..23] | [16..19] | [12..15] | [8..11] | [2..7] | [0..1] |
|--------|------------|----------|----------|----------|----------|----------|---------|--------|--------|
| R      | Register   | Func8    | Func8    | W        | Y        | Z        | X       | Opcode | 1      |
| E      | Eight      | Func8    | Func8    | imm8     | imm8     | Z        | X       | Opcode | 1      |
| A      | Assignment | imm16    | imm16    | imm16    | imm16    | Z        | X       | Opcode | 2      |
| B      | Branch     | imm16    | imm16    | imm16    | Y        | imm16    | X       | Opcode | 2      |

48-bit format: (X for extended)

| Format | Mnemonic   | [20..47] | [16..19] | [12..15] | [8..11] | [2..7] | [0..1] |
|--------|------------|----------|----------|----------|---------|--------|--------|
| XA     | Assignment | imm32    | imm32    | C/Z      | A/X     | Opcode | 3      |
| XB     | Branch     | imm32    | B/Y      | imm32    | A/X     | Opcode | 3      |

The format A is the one most likely to suffer from linker reallocation,
so we chose to keep its immediate contiguous to simplify the tooling.

The S format reduces the pressure on opcode allocation versus implementing
shift operations using the A format.

The opcode field determines the format of the instruction.
Within an opcode, the function field (Func) determines the specific instruction.
X and Y are source registers.
Z and W are destination registers.

If an instruction does not use all the fields that its format provides,
those fields must be set to zero.

## 48-bits opcodes

There is a limited number of instructions that add significant value to their
32-bit counterparts, thanks to the longer immediates.

These are:

| Class                         | Variants                       | Count |
|-------------------------------|--------------------------------|-------|
| Load signed data operations   | B, H, W, D                     | 4     |
| Load unsigned data operations | B, H, W, D                     | 4     |
| Load                          | Q, A                           | 2     |
| Store                         | B, H, W, D, Q, A               | 6     |
| Load (FP)                     | H, W, D, Q (2 banks)           | 8     |
| Store (FP)                    | H, W, D, Q (2 banks)           | 8     |
| Branch (data)                 | EQ, NE, LT.S, LT.U, GE.S, GE.U | 6     |
| Branch (address)              | EQ, NE, LT, LT, Z, NZ          | 6     |
| Call                          | A                              | 1     |
| Jump                          | A                              | 1     |
| `add.l`                       | W, D, Q                        | 3     |
| `and.l`                       | W, D, Q                        | 3     |
| `or.l`                        | W, D, Q                        | 3     |
| `xor.l`                       | W, D, Q                        | 3     |
| `slt.l`                       | U, S                           | 2     |
| `aiupc.l`                     | A                              | 1     |
| `lui.ld`                      | D (offset 32)                  | 1     |
| `lui.lq0`                     | Q (offset 64)                  | 1     |
| `lui.lq1`                     | Q (offset 96)                  | 1     |

which consumes exactly the 64 opcodes.

Dropping Q would free 11 opcodes (load unsigned D, plus all the Q).
Then it would consume 53 out of the 64 opcodes available.

Another opportunity comes from realising that it doesn't make sense to
have up to nine instructions that can load into the ZR register.
We could keep one as a prefetch instruction and overlap the others with
instructions that can benefit from the long immediate and a register.
Calls, jumps, LUIs, AUIPCs, etc., might be made to fit this pattern.
That could free up to seven opcodes.
Then it would consume 56 out of the 64 opcodes available.
