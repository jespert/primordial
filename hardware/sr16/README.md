# SR16 Computer

Specification and emulator for a 16-bit computer with a split register file.

Answers the crucial question: "what if RISC-V and M68k had a 16-bit baby?"

## Features

- 16 integer registers and 16 pointer registers.
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
- All code is position-independent, via using base pointer (BP) register.
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

## Pointer arithmetic

The reason why splitting the register file into data and address registers
works is because many fewer operations are possible with the address registers.
These are dictated by the operations supported by pointer arithmetic:

- Add an immediate integer to a pointer (result: pointer).
- Copy a pointer to another pointer (pseudo-instruction: `add.ai %C, %A, 0`)
- Add an integer register to a pointer (result: pointer).
- Subtract an integer register from a pointer (result: pointer).
- Subtract two pointers (result: integer).
- Convert a pointer to an integer.
- Convert an integer to a pointer.
- Load address from memory.
- Store address to memory.

(total: 8)

Additionally, we need to consider conditional branches that depend on pointers:

- Branch if equal
- Branch if not equal
- Branch if less than
- Branch if greater than or equal
- Branch if zero (null)
- Branch if non-zero (not null)

(total: 6)

Finally, we need to introduce a `jalz` control instruction because there are no
hardcoded-to-zero pointer registers.

(total: 1)

So, for the price of 15 new instructions, which share most of their hardware
with already existing instructions, we can double the number of general-purpose
registers.

## Registers

Conceptually, the register file can be seen as a set of 32 general-purpose
registers split into two banks.

| Type    | n  | Register | Alias  | Purpose                                       |
|---------|----|----------|--------|-----------------------------------------------|
| Integer | 0  | 0        | zr     | Zero (hardcoded)                              |
| Integer | 1  | 1        | z6     | Saved register                                |
| Integer | 2  | 2        | z5     | Saved register                                |
| Integer | 3  | 3        | z4     | Saved register                                |
| Integer | 4  | 4        | z3     | Saved register                                |
| Integer | 5  | 5        | z2     | Saved register                                |
| Integer | 6  | 6        | z1     | Saved register                                |
| Integer | 7  | 7        | z0     | Saved register                                |
| Integer | 8  | 8        | y0     | Temporary register                            |
| Integer | 9  | 9        | y1     | Temporary register                            |
| Integer | 10 | a        | x0     | Argument register                             |
| Integer | 11 | b        | x1     | Argument register                             |
| Integer | 12 | c        | x2     | Argument register                             |
| Integer | 13 | d        | x3     | Argument register                             |
| Integer | 14 | e        | x4     | Argument register                             |
| Integer | 15 | f        | x5     | Argument register                             |
| Pointer | 16 | 0        | bp     | Base pointer (start of program memory region) |
| Pointer | 17 | 1        | c6     | Saved register                                |
| Pointer | 18 | 2        | c5     | Saved register                                |
| Pointer | 19 | 3        | c4     | Saved register                                |
| Pointer | 20 | 4        | c3     | Saved register                                |
| Pointer | 21 | 5        | c2     | Saved register                                |
| Pointer | 22 | 6        | c1     | Saved register                                |
| Pointer | 23 | 7        | c0, fp | Saved register, Frame pointer                 |
| Pointer | 24 | 8        | b0, mp | Temporary register, millicode return pointer  |
| Pointer | 25 | 9        | b1     | Temporary register                            |
| Pointer | 26 | a        | a0     | Argument register                             |
| Pointer | 27 | b        | a1     | Argument register                             |
| Pointer | 28 | c        | a2     | Argument register                             |
| Pointer | 29 | d        | a3     | Argument register                             |
| Pointer | 30 | e        | rp     | Return pointer                                |
| Pointer | 31 | f        | sp     | Stack pointer                                 |

The convention for saved registers (S) grows downwards to mitigate the risk of
off-by-one errors between the alias and the register number.

## Conventions

- The stack grows downwards.
- The stack is aligned to 16 bits.

## Instruction encoding

### Formats

Like R16, we still only have three instruction formats conceptually.
The register bank (integer or pointer) will be determined by the operation
(opcode and function).

| Format | Mnemonic   | [28..31] | [24..27] | [20..23] | [16..19] | [12..15] | [8..11] | [4..7] | [0..3] |
|--------|------------|----------|----------|----------|----------|----------|---------|--------|--------|
| A      | Assignment | Opcode   | C/Z      | Func4    | A/X      | imm16    | imm16   | imm16  | imm16  |
| B      | Branch     | Opcode   | Func4    | B/Y      | A/X      | imm16    | imm16   | imm16  | imm16  |
| R      | Register   | Opcode   | C/Z      | B/Y      | A/X      | W        | Func12  | Func12 | Func12 |

Unlike most RISC architectures, we do not need additional formats to
support calls, jumps, or branches because the A and B formats
can already reference any memory position.

The opcode field determines the format of the instruction.
Within an opcode, the function field (Func) determines the specific instruction.
X and Y are integer source registers.
Z and W are integer destination registers.
A and B are source pointer registers.
C is a destination pointer register.


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

| Instruction           | Opcode | Func | Semantics                                              |
|-----------------------|--------|------|--------------------------------------------------------|
| `jalz %Z, %A, offset` | 8      | 0    | Jump and link saving to integer register %Z            |
| `rjump %A, offset`    | 8      | 0    | Pseudo-instruction: `jalz %zr, %A, target`             |
| `jump target`         | 8      | 0    | Pseudo-instruction: `jalz %zr, %bp, target`            |
| `ret`                 | 8      | 0    | Pseudo-instruction: `jalz %zr, %rp, 0`                 |
| `mret`                | 8      | 0    | Pseudo-instruction: `jalz %zr, %mp, 0` (millicode)     |
| `jal %C, %A, offset`  | 8      | 1    | Jump and link                                          |
| `rcall %A, offset`    | 8      | 1    | Pseudo-instruction: `jal %rp, %A, offset`              |
| `rmcall %A, offset`   | 8      | 1    | Pseudo-instruction: `jal %mp, %A, offset`              |
| `call target`         | 8      | 1    | Pseudo-instruction: `jal %rp, %bp, target`             |
| `mcall target`        | 8      | 1    | Pseudo-instruction: `jal %mp, %bp, target` (millicode) |

Compared to R16, we lose a bit of magic here because we can no longer use the
ZR register as a destination register in `jal`. Technically, we could use,
e.g., %b1 as a discard register and keep a single instruction for unconditional
flow control. However, we can simply add a new instruction `jalz`, which is
identical to `jal`, except that it writes to an integer register. If you set
it to ZR, it behaves as a `jump`.

The flexibility of being able to choose the destination register
allows the use of millicode, as in RISC-V. The most common use case is to
reduce the tedium of function prologues and epilogues. The `mcall` and `mret`
pseudo-instructions are used for this purpose.

### Conditional control flow

We use RISC-V style conditional with fused test-and-branch instructions due to
their simplicity and excellent ergonomics. Where necessary, we have enforced
signedness suffixes (s, u) to mitigate accidental misuse.

| Instruction            | Opcode | Func | Binary | Semantics                              |
|------------------------|--------|------|--------|----------------------------------------|
| `beq %Y, %X, target`   | 4      | 0    | 0000   | Branch to target if %Y = %X            |
| `bne %Y, %X, target`   | 4      | 1    | 0001   | Branch to target if %Y ≠ %X            |
| `blt.s %Y, %X, target` | 4      | 4    | 0100   | Branch to target if %Y < %X (signed)   |
| `blt.u %Y, %X, target` | 4      | 5    | 0101   | Branch to target if %Y < %X (unsigned) |
| `bge.s %Y, %X, target` | 4      | 6    | 0110   | Branch to target if %Y ≥ %X (signed)   |
| `bge.u %Y, %X, target` | 4      | 7    | 0111   | Branch to target if %Y ≥ %X (unsigned) |
| `beq.a %B, %A, target` | 4      | 8    | 1000   | Branch to target if %Y = %X            |
| `bne.a %B, %A, target` | 4      | 9    | 1001   | Branch to target if %Y ≠ %X            |
| `bzr.a %B, %A, target` | 4      | a    | 1010   | Branch to target if %Y = 0             |
| `bnz.a %B, %A, target` | 4      | b    | 1011   | Branch to target if %Y ≠ 0             |
| `blt.a %B, %A, target` | 4      | c    | 1100   | Branch to target if %Y < %X            |
| `bge.a %B, %A, target` | 4      | e    | 1110   | Branch to target if %Y ≥ %X            |

### Load from memory

There are two instructions to load bytes, depending on whether it should be
treated as a signed or unsigned value.
We have enforced a suffix (s, u) on both to mitigate accidental misuse.
This is unnecessary for halfwords because they match the register size.

| Instruction              | Opcode | Func | Binary | Semantics                                              |
|--------------------------|--------|------|--------|--------------------------------------------------------|
| `load.sb %Z, %A, offset` | 9      | 0    | 0000   | Read byte at [%A+offset], sign extend, and write to %Z |
| `load.h %Z, %A, offset`  | 9      | 1    | 0001   | Read half at [%A+offset] and write to %Z               |
| `load.ub %Z, %A, offset` | 9      | 4    | 0100   | Read byte at [%A+offset], zero extend, and write to %Z |
| `load.a %C, %A, offset`  | 9      | 9    | 1001   | Read address at [%A+offset] and write to %C            |

### Store to memory

Note that signedness is irrelevant for stores, so a single instruction per
operand size suffices.

| Instruction              | Opcode | Func | Binary | Semantics                                   |
|--------------------------|--------|------|--------|---------------------------------------------|
| `store.b %Y, %A, offset` | 5      | 0    | 0000   | Read byte at %Y and write to [%A+offset]    |
| `store.h %Y, %A, offset` | 5      | 1    | 0001   | Read half at %Y and write to [%A+offset]    |
| `store.a %B, %A, offset` | 5      | 9    | 1001   | Read address at %B and write to [%A+offset] |

### Arithmetic with immediates

Analogous instructions are provided for bytes and halfwords.
When we talk about logical instead of bitwise operations below,
we assume that Booleans are stored as 0 or 1 exactly.

For bytes:

| Instruction          | Opcode | Func | Semantics                                            |
|----------------------|--------|------|------------------------------------------------------|
| `and.bi %Z, %X, imm` | e      | 0    | Bitwise AND / Logical AND                            |
| `or.bi  %Z, %X, imm` | e      | 1    | Bitwise OR / Logical OR                              |
| `xor.bi %Z, %X, imm` | e      | 2    | Bitwise XOR                                          |
| `inv.b  %Z, %X`      | e      | 2    | Pseudo-instruction: bitwise NOT (`xor.b %Z, %X, -1`) |
| `not.b  %Z, %X`      | e      | 2    | Pseudo-instruction: logical NOT (`xor.b %Z, %X, 1`)  |
| `sra.bi %Z, %X, imm` | e      | 3    | Read half at %Y and write to (%X+offset)             |
| `srl.bi %Z, %X, imm` | e      | 4    | Read half at %Y and write to (%X+offset)             |
| `sll.bi %Z, %X, imm` | e      | 5    | Read half at %Y and write to (%X+offset)             |
| `add.bi %Z, %X, imm` | e      | 6    | Read half at %Y and write to (%X+offset)             |

For halfwords:

| Instruction          | Opcode | Func | Semantics                                            |
|----------------------|--------|------|------------------------------------------------------|
| `and.hi %Z, %X, imm` | f      | 0    | Bitwise AND / Logical AND                            |
| `or.hi  %Z, %X, imm` | f      | 1    | Bitwise OR / Logical OR                              |
| `xor.hi %Z, %X, imm` | f      | 2    | Bitwise XOR                                          |
| `inv.h  %Z, %X`      | f      | 2    | Pseudo-instruction: bitwise NOT (`xor.h %Z, %X, -1`) |
| `not.h  %Z, %X`      | f      | 2    | Pseudo-instruction: logical NOT (`xor.h %Z, %X, 1`)  |
| `sra.hi %Z, %X, imm` | f      | 3    | Shift right (arithmetic)                             |
| `srl.hi %Z, %X, imm` | f      | 4    | Shift right (logic)                                  |
| `sll.hi %Z, %X, imm` | f      | 5    | Shift left (logic)                                   |
| `add.hi %Z, %X, imm` | f      | 6    | Addition                                             |

Additionally, the below instructions work on full registers and/or are suitable
for any operand size.

| Instruction          | Opcode | Func | Semantics                                  |
|----------------------|--------|------|--------------------------------------------|
| `slt.si %Z, %X, imm` | b      | 0    | Set %Z to 1 if %X < imm (signed), else 0   |
| `slt.ui %Z, %X, imm` | b      | 1    | Set %Z to 1 if %X < imm (unsigned), else 0 |
| `add.ai %C, %A, imm` | b      | 8    | Add integer to pointer                     |

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

| Instruction        | Opcode | Func  | Semantics                                 |
|--------------------|--------|-------|-------------------------------------------|
| `slt.s %Z, %Y, %X` | 0      | 200   | Set %Z to 1 if %Y < %X (signed), else 0   |
| `slt.u %Z, %Y, %X` | 0      | 201   | Set %Z to 1 if %Y < %X (unsigned), else 0 |
| `add.a %C, %B, %X` | 0      | 208   | Add integer offset to pointer             |
| `sub.a %C, %B, %X` | 0      | 209   | Subtract integer offset from pointer      |
| `diff %Z, %B, %A`  | 0      | 210   | Difference between two pointers           |
| `ptoi %Z, %A`      | 0      | 211   | Pointer to integer conversion             |
| `itop %A, %X`      | 0      | 212   | Integer to pointer conversion             |

## Extensibility

Even if it is not a primary goal, the design is quite extensible.
We do not foresee running out of the R format minor opcode space
(4096 operations), so we will probably only ever need one of those.
Across the other two formats, A and B, we can define up to
(15 opcodes) * (16 functions) = 240 operations that support 16-bit immediates.
