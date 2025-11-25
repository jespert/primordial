# R16 Computer

Specification and emulator for an extremely simple 16-bit computer.

You can think of it as a very late RISC competitor to the PDP-11 in which
C and UNIX were born.

## Features

- 16-bit flat address space (64 KiB).
- 16-bit registers.
- 32-bit instructions.
- Integer only.
- Little-endian.
- Memory-mapped I/O.
- Byte addressable memory.

Beyond this, the design has some features that make it suitable both in
education and environments where a toolchain is not yet available.

- Only three instruction formats.
- All instruction fields are aligned to four bits (one hexadecimal digit).
- 16-bit immediate values, so any value or address can be used or
  referenced in a single instruction.
- Absolute (not relative!) addressing, even for branches.
- Immediates appear in the instruction formats as they are interpreted,
  with no implicit shifts or bit relocations.
- First-class support for byte operands in ALU, not just halfwords.
- Trivial hardware interface: can work without an operating system.
- User programs run on the upper half of the memory (max 32 KiB),
  so loops and arrays can use the less error-prone signed counters and indices.
- The operating system (if any) and MMIO use the bottom half of the memory.

All the above make it an easier target for educational and bootstrapping
projects than a more general-purpose architecture like RISC-V could be.
Finally, 32 KiB (64 KiB with the operating system) of memory is small enough
to visually inspect, while big enough to write a simple assembler or compiler,
as proven by history.

## Registers

The convention for saved registers (S) grows downwards to mitigate the risk of
off-by-one errors between the alias and the register number.

| Register | Alias  | Purpose                                        |
|----------|--------|------------------------------------------------|
| 0        | ZR     | Zero (hardcoded)                               |
| 1        | S6     | Saved register                                 |
| 2        | S5     | Saved register                                 |
| 3        | S4     | Saved register                                 |
| 4        | S3     | Saved register                                 |
| 5        | S2     | Saved register                                 |
| 6        | S1     | Saved register                                 |
| 7        | S0, FP | Saved register, Frame pointer                  |
| 8        | T0     | Temporary register, alternative return pointer |
| 9        | T1     | Temporary register                             |
| a        | A0     | Argument register                              |
| b        | A1     | Argument register                              |
| c        | A2     | Argument register                              |
| d        | A3     | Argument register                              |
| e        | RP     | Return pointer                                 |
| f        | SP     | Stack pointer                                  |

## Conventions

- The stack grows downwards.
- The stack is aligned to 16 bits.

## Instruction encoding

### Formats

There are only three formats. Unlike most RISC architectures, we do not need
additional formats to support calls, jumps, or branches because the A and B
formats can already reference any memory position.

| Format | Mnemonic   | [28..31] | [24..27] | [20..23] | [16..19] | [12..15] | [8..11] | [4..7] | [0..3] |
|--------|------------|----------|----------|----------|----------|----------|---------|--------|--------|
| A      | Assignment | Opcode   | Z        | Func     | X        | imm16    | imm16   | imm16  | imm16  |
| B      | Branch     | Opcode   | Func     | Y        | X        | imm16    | imm16   | imm16  | imm16  |
| R      | Register   | Opcode   | Z        | Y        | X        | W        | Func    | Func   | Func   |

The opcode field determines the format of the instruction.
Within an opcode, the function field (Func) determines the specific instruction.
X and Y are source registers.
Z and W are destination registers.

If an instruction does not use all the fields that its format provides,
those fields must be set to zero.

### Opcodes

Opcodes not listed are reserved for future use or custom extensions.

| Hexadecimal | Binary | Format | Usage                                               |
|-------------|--------|--------|-----------------------------------------------------|
| 0           | 0000   | R      | Operations on registers only and special operations |
| 4           | 0100   | B      | Conditional control flow: branch                    |
| 5           | 0101   | B      | Store to memory                                     |
| 8           | 1000   | A      | Unconditional control flow: call, jump, return      |
| 9           | 1001   | A      | Load from memory                                    |
| b           | 1011   | A      | Generic arithmetic with immediates                  |
| e           | 1110   | A      | Byte arithmetic with immediates                     |
| f           | 1111   | A      | Halfword arithmetic with immediates                 |

To simplify the most basic hardware implementations, the two MSBs of the opcode
determine the format in the base architecture:

- 00: R format
- 01: B format
- 10: A format
- 11: A format

For the A and B formats, the two LSBs also determine the type:

- 00: Control flow
- 01: Memory operations
- 10: ALU on bytes
- 11: ALU on halfwords

### Unconditional control flow

There is a single instruction format for unconditional flow control:
`jal %Z, %X, offset` (jump and link). It stores into %Z the address of the next
instruction and jumps to (%Z+offset). We provide a bunch of pseudo-instructions
for clarity on disassembly, but ultimately they are all one and the same.

The flexibility of being able to choose the destination register
allows the use of millicode, as in RISC-V. The most common use case is to
reduce the tedium of function prologues and epilogues. The `mcall` and `mret`
pseudo-instructions are used for this purpose.

| Instruction          | Opcode | Func | Semantics                                            |
|----------------------|--------|------|------------------------------------------------------|
| `jal %Z, %X, offset` | 8      | 1    | Universal unconditional flow control                 |
| `call target`        | 8      | 1    | Pseudo-instruction: jal %rp, %zr, target             |
| `mcall target`       | 8      | 1    | Pseudo-instruction: jal %t0, %zr, target (millicode) |
| `rcall %X, offset`   | 8      | 1    | Pseudo-instruction: jal %rp, %X, offset              |
| `jump target`        | 8      | 1    | Pseudo-instruction: jal %zr, %zr, target             |
| `rjump %X, offset`   | 8      | 1    | Pseudo-instruction: jal %zr, %X, offset              |
| `ret`                | 8      | 1    | Pseudo-instruction: jal %zr, %rp, 0                  |
| `mret`               | 8      | 1    | Pseudo-instruction: jal %zr, %t0, 0 (millicode)      |

### Conditional control flow

We use RISC-V style conditional with fused test-and-branch instructions due to
their simplicity and excellent ergonomics. Where necessary, we have enforced
signedness suffixes (s, u) to mitigate accidental misuse.

| Instruction            | Opcode | Func | Binary | Semantics                              |
|------------------------|--------|------|--------|----------------------------------------|
| `beq %Y, %X, target`   | 4      | 0    | 0000   | Branch to target if %X = %Y            |
| `bne %Y, %X, target`   | 4      | 1    | 0001   | Branch to target if %X ≠ %Y            |
| `blt.s %Y, %X, target` | 4      | 8    | 1000   | Branch to target if %X < %Y (signed)   |
| `bge.s %Y, %X, target` | 4      | a    | 1010   | Branch to target if %X ≥ %Y (signed)   |
| `blt.u %Y, %X, target` | 4      | c    | 1100   | Branch to target if %X < %Y (unsigned) |
| `bge.u %Y, %X, target` | 4      | e    | 1110   | Branch to target if %X ≥ %Y (unsigned) |

### Load from memory

There are two instructions to load bytes, depending on whether it should be
treated as a signed or unsigned value.
We have enforced a suffix (s, u) on both to mitigate accidental misuse.
This is unnecessary for halfwords because they match the register size.

| Instruction              | Opcode | Func | Binary | Semantics                                              |
|--------------------------|--------|------|--------|--------------------------------------------------------|
| `load.bs %Z, %X, offset` | 9      | 0    | 0000   | Read byte at (%X+offset), sign extend, and write to %Z |
| `load.h %Z, %X, offset`  | 9      | 1    | 0001   | Read half at (%X+offset) and write to %Z               |
| `load.bu %Z, %X, offset` | 9      | 4    | 0100   | Read byte at (%X+offset), zero extend, and write to %Z |

### Store to memory

Note that signedness is irrelevant for stores, so a single instruction per
operand size suffices.

| Instruction              | Opcode | Func | Binary | Semantics                                |
|--------------------------|--------|------|--------|------------------------------------------|
| `store.b %Y, %X, offset` | 5      | 0    | 0000   | Read byte at %Y and write to (%X+offset) |
| `store.h %Y, %X, offset` | 5      | 1    | 0001   | Read half at %Y and write to (%X+offset) |

### Arithmetic with immediates

Analogous instructions are provided for bytes and halfwords.
When we talk about logical instead of bitwise operations below,
we assume that Booleans are stored as 0 or 1 exactly.

For bytes:

| Instruction          | Opcode | Func | Semantics                                            |
|----------------------|--------|------|------------------------------------------------------|
| `and.ib %Z, %X, imm` | e      | 0    | Bitwise AND / Logical AND                            |
| `or.ib  %Z, %X, imm` | e      | 1    | Bitwise OR / Logical OR                              |
| `xor.ib %Z, %X, imm` | e      | 2    | Bitwise XOR                                          |
| `inv.b  %Z, %X`      | e      | 2    | Pseudo-instruction: bitwise NOT (`xor.b %Z, %X, -1`) |
| `not.b  %Z, %X`      | e      | 2    | Pseudo-instruction: logical NOT (`xor.b %Z, %X, 1`)  |
| `sra.ib %Z, %X, imm` | e      | 3    | Read half at %Y and write to (%X+offset)             |
| `srl.ib %Z, %X, imm` | e      | 4    | Read half at %Y and write to (%X+offset)             |
| `sll.ib %Z, %X, imm` | e      | 5    | Read half at %Y and write to (%X+offset)             |
| `add.ib %Z, %X, imm` | e      | 6    | Read half at %Y and write to (%X+offset)             |

For halfwords:

| Instruction          | Opcode | Func | Semantics                                            |
|----------------------|--------|------|------------------------------------------------------|
| `and.ih %Z, %X, imm` | f      | 0    | Bitwise AND / Logical AND                            |
| `or.ih  %Z, %X, imm` | f      | 1    | Bitwise OR / Logical OR                              |
| `xor.ih %Z, %X, imm` | f      | 2    | Bitwise XOR                                          |
| `inv.h  %Z, %X`      | f      | 2    | Pseudo-instruction: bitwise NOT (`xor.h %Z, %X, -1`) |
| `not.h  %Z, %X`      | f      | 2    | Pseudo-instruction: logical NOT (`xor.h %Z, %X, 1`)  |
| `sra.ih %Z, %X, imm` | f      | 3    | Shift right (arithmetic)                             |
| `srl.ih %Z, %X, imm` | f      | 4    | Shift right (logic)                                  |
| `sll.ih %Z, %X, imm` | f      | 5    | Shift left (logic)                                   |
| `add.ih %Z, %X, imm` | f      | 6    | Addition                                             |

Additionally, the below instructions work on full registers but are suitable
for any operand size.

| Instruction          | Opcode | Func | Semantics                                  |
|----------------------|--------|------|--------------------------------------------|
| `slt.is %Z, %X, imm` | b      | 0    | Set %Z to 1 if %X < imm (signed), else 0   |
| `slt.iu %Z, %X, imm` | b      | 1    | Set %Z to 1 if %X < imm (unsigned), else 0 |

### Operations on registers only

Special operations:

| Instruction         | Opcode | Func | Semantics                                     |
|---------------------|--------|------|-----------------------------------------------|
| Illegal instruction | 0      | 0    | Traps on execution of zero-initialised memory |

Analogous instructions are provided for bytes and halfwords.
When we talk about logical instead of bitwise operations below,
we assume that Booleans are stored as 0 or 1 exactly.

For bytes:

| Instruction        | Opcode | Func | Semantics                 |
|--------------------|--------|------|---------------------------|
| `and.b %Z, %Y, %X` | 0      | 100  | Bitwise AND / Logical AND |
| `or.b  %Z, %Y, %X` | 0      | 101  | Bitwise OR / Logical OR   |
| `xor.b %Z, %Y, %X` | 0      | 102  | Bitwise XOR               |
| `sra.b %Z, %Y, %X` | 0      | 103  | Shift right (arithmetic)  |
| `srl.b %Z, %Y, %X` | 0      | 104  | Shift right (logical)     |
| `sll.b %Z, %Y, %X` | 0      | 105  | Shift left (logical)      |
| `add.b %Z, %Y, %X` | 0      | 106  | Addition                  |
| `sub.b %Z, %Y, %X` | 0      | 107  | Subtraction               |

For halfwords:

| Instruction        | Opcode | Func | Semantics                 |
|--------------------|--------|------|---------------------------|
| `and.h %Z, %Y, %X` | 0      | 110  | Bitwise AND / Logical AND |
| `or.h  %Z, %Y, %X` | 0      | 111  | Bitwise OR / Logical OR   |
| `xor.h %Z, %Y, %X` | 0      | 112  | Bitwise XOR               |
| `sra.h %Z, %Y, %X` | 0      | 113  | Shift right (arithmetic)  |
| `srl.h %Z, %Y, %X` | 0      | 114  | Shift right (logical)     |
| `sll.h %Z, %Y, %X` | 0      | 115  | Shift left (logical)      |
| `add.h %Z, %Y, %X` | 0      | 116  | Addition                  |
| `sub.h %Z, %Y, %X` | 0      | 117  | Subtraction               |

Additionally, the below instructions work on full registers but are suitable
for any operand size.

| Instruction        | Opcode | Func | Semantics                                 |
|--------------------|--------|------|-------------------------------------------|
| `slt.s %Z, %Y, %X` | 0      | 200  | Set %Z to 1 if %Y < %X (signed), else 0   |
| `slt.u %Z, %Y, %X` | 0      | 201  | Set %Z to 1 if %Y < %X (unsigned), else 0 |

## Extensibility

Even if it is not a primary goal, the design is quite extensible.
We do not foresee running out of the R format minor opcode space
(4096 operations), so we will probably only ever need one of those.
Across the other two formats, A and B, we can define up to
(15 opcodes) * (16 functions) = 240 operations that support 16-bit immediates.
