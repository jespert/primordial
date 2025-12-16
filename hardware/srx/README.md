# SRX Computer

Specification and emulator for the SRX computer architecture.
It supports 16, 32, 64 and 128-bit architectures.

SRX has separate data and address registers and variable-length instructions
(S = split register file, X = variable-length instructions).

## Registers

The convention for saved registers (S) grows downwards to mitigate the risk of
off-by-one errors between the alias and the register number.

| Type | Register | Alias | Purpose                                          |
|------|----------|-------|--------------------------------------------------|
| D    | 0        | ZR    | Zero (hardcoded)                                 |
| D    | 1        |       | Argument register 4                              |
| D    | 2        |       | Argument register 3                              |
| D    | 3        |       | Argument register 2                              |
| D    | 4        |       | Argument register 1                              |
| D    | 5        |       | Argument register 0                              |
| D    | 6        |       | Temporary register 1                             |
| D    | 7        |       | Temporary register 0                             |
| D    | 8        |       | Saved register 0                                 |
| D    | 9        |       | Saved register 1                                 |
| D    | a        |       | Saved register 2                                 |
| D    | b        |       | Saved register 3                                 |
| D    | c        |       | Saved register 4                                 |
| D    | d        |       | Saved register 5                                 |
| D    | e        |       | Saved register 6                                 |
| D    | f        |       | Saved register 7                                 |
| A    | 0        | SP    | Stack pointer                                    |
| A    | 1        | BP    | Base pointer                                     |
| A    | 2        | TP    | Thread pointer                                   |
| A    | 3        |       | Argument register 2                              |
| A    | 4        |       | Argument register 1                              |
| A    | 5        |       | Argument register 0                              |
| A    | 6        |       | Temporary register 1                             |
| A    | 7        |       | Temporary register 0, alternative return pointer |
| A    | 8        | FP    | Saved register 0, Frame pointer                  |
| A    | 9        |       | Saved register 1                                 |
| A    | a        |       | Saved register 2                                 |
| A    | b        |       | Saved register 3                                 |
| A    | c        |       | Saved register 4                                 |
| A    | d        |       | Saved register 5                                 |
| A    | e        |       | Saved register 6                                 |
| A    | f        | RP    | Return pointer                                   |

## Instruction encoding

### Formats

Instructions can be 16, 32 or 48-bits long.

Immediates can be 8, 16 or 32-bits-long.
This consumes many opcodes but allows a very natural set of instructions.
Compared to RISC-V, it supports much further branches, jumps and calls.

16-bit format:

| Format | Mnemonic  | [12..15] | [8..11] | [4..7] | [2..3] | [0..1] |
|--------|-----------|----------|---------|--------|--------|--------|
| HR     | Register  | Z        | X       | Func4  | Opcode | 0      |
| HF     | Function  | ZX       | Func8   | Func8  | Opcode | 0      |
| HI     | Immediate | ZX       | imm8    | imm8   | Opcode | 0      |

32-bit format:

| Format | Mnemonic   | [28..31] | [24..27] | [20..23] | [16..19] | [12..15] | [8..11] | [2..7] | [0..1] |
|--------|------------|----------|----------|----------|----------|----------|---------|--------|--------|
| WR     | Register   | Func8    | Func8    | W        | Y        | Z        | X       | Opcode | 1      |
| WS     | Shift      | Func8    | Func8    | imm8     | imm8     | Z        | X       | Opcode | 1      |
| WA     | Assignment | imm16    | imm16    | imm16    | imm16    | Z        | X       | Opcode | 2      |
| WB     | Branch     | imm16    | imm16    | imm16    | Y        | imm16    | X       | Opcode | 2      |


48-bit format:

| Format | Mnemonic   | [20..47] | [16..19] | [12..15] | [8..11] | [2..7] | [0..1] |
|--------|------------|----------|----------|----------|---------|--------|--------|
| DA     | Assignment | imm32    | imm32    | X        | Z       | Opcode | 3      |
| DB     | Branch     | imm32    | Y        | X        | imm32   | Opcode | 3      |

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
