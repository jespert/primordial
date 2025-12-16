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

## 48-bit instructions

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

| Instruction               | [20..47] | [16..19] | [12..15] | [8..11] | [2..7] | [0..1] |
|---------------------------|----------|----------|----------|---------|--------|--------|
| `call.x imm32`            | imm32    | imm32    | 0        | 0       | 0      | 3      |
| `jmp.x imm32`             | imm32    | imm32    | 0        | 1       | 0      | 3      |
| `jeq.x imm32`             | imm32    | imm32    | 0        | 2       | 0      | 3      |
| `jne.x imm32`             | imm32    | imm32    | 0        | 3       | 0      | 3      |
| `jlt.x imm32`             | imm32    | imm32    | 0        | 4       | 0      | 3      |
| `jle.x imm32`             | imm32    | imm32    | 0        | 5       | 0      | 3      |
| `jge.x imm32`             | imm32    | imm32    | 0        | 6       | 0      | 3      |
| `jgt.x imm32`             | imm32    | imm32    | 0        | 7       | 0      | 3      |
| `jlo.x imm32`             | imm32    | imm32    | 0        | 8       | 0      | 3      |
| `jls.x imm32`             | imm32    | imm32    | 0        | 9       | 0      | 3      |
| `jhs.x imm32`             | imm32    | imm32    | 0        | a       | 0      | 3      |
| `jhi.x imm32`             | imm32    | imm32    | 0        | b       | 0      | 3      |
| `jpi.x imm32`             | imm32    | imm32    | 0        | c       | 0      | 3      |
| `jni.x imm32`             | imm32    | imm32    | 0        | d       | 0      | 3      |
| `jov.l imm32`             | imm32    | imm32    | 0        | e       | 0      | 3      |
| `jno.x imm32`             | imm32    | imm32    | 0        | f       | 0      | 3      |
| `load.sbx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 0      | 3      |
| `beqz.x %X, imm32`        | imm32    | imm32    | 0        | X       | 1      | 3      |
| `load.shx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 1      | 3      |
| `bnez.x %X, imm32`        | imm32    | imm32    | 0        | X       | 2      | 3      |
| `load.swx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 2      | 3      |
| `bltz.x %X, imm32`        | imm32    | imm32    | 0        | X       | 3      | 3      |
| `load.sdx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 3      | 3      |
| `bgez.x %X, imm32`        | imm32    | imm32    | 0        | X       | 4      | 3      |
| `load.ubx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 4      | 3      |
| `beqz.ax %A, imm32`       | imm32    | imm32    | 0        | A       | 5      | 3      |
| `load.uhx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 5      | 3      |
| `bnez.ax %A, imm32`       | imm32    | imm32    | 0        | A       | 6      | 3      |
| `load.uwx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 6      | 3      |
| `bltz.ax %A, imm32`       | imm32    | imm32    | 0        | A       | 7      | 3      |
| `load.udx %Z, %A, imm32`  | imm32    | imm32    | Z != 0   | A       | 7      | 3      |
| `bgez.ax %A, imm32`       | imm32    | imm32    | 0        | A       | 8      | 3      |
| `load.qx %Z, %A, imm32`   | imm32    | imm32    | Z != 0   | A       | 8      | 3      |
| `load.ax %C, %A, imm32`   | imm32    | imm32    | C        | A       | 9      | 3      |
| `store.ax %Y, %A, imm32`  | imm32    | B        | imm32    | A       | 10     | 3      |
| `store.qx %Y, %A, imm32`  | imm32    | Y        | imm32    | A       | 11     | 3      |
| `store.bx %Y, %A, imm32`  | imm32    | Y        | imm32    | A       | 12     | 3      |
| `store.hx %Y, %A, imm32`  | imm32    | Y        | imm32    | A       | 13     | 3      |
| `store.wx %Y, %A, imm32`  | imm32    | Y        | imm32    | A       | 14     | 3      |
| `store.dx %Y, %A, imm32`  | imm32    | Y        | imm32    | A       | 15     | 3      |
| `load.fh %FY, %A, imm32`  | imm32    | imm32    | F0       | A       | 16     | 3      |
| `load.fh %F1, %A, imm32`  | imm32    | imm32    | F1       | A       | 17     | 3      |
| `load.fw %F0, %A, imm32`  | imm32    | imm32    | F0       | A       | 18     | 3      |
| `load.fw %F1, %A, imm32`  | imm32    | imm32    | F1       | A       | 19     | 3      |
| `load.fd %F0, %A, imm32`  | imm32    | imm32    | F0       | A       | 20     | 3      |
| `load.fd %F1, %A, imm32`  | imm32    | imm32    | F1       | A       | 21     | 3      |
| `load.fq %F0, %A, imm32`  | imm32    | imm32    | F0       | A       | 22     | 3      |
| `load.fq %F1, %A, imm32`  | imm32    | imm32    | F1       | A       | 23     | 3      |
| `store.fh %FY, %A, imm32` | imm32    | F0       | imm32    | A       | 16     | 3      |
| `store.fh %F1, %A, imm32` | imm32    | F1       | imm32    | A       | 17     | 3      |
| `store.fw %F0, %A, imm32` | imm32    | F0       | imm32    | A       | 18     | 3      |
| `store.fw %F1, %A, imm32` | imm32    | F1       | imm32    | A       | 19     | 3      |
| `store.fd %F0, %A, imm32` | imm32    | F0       | imm32    | A       | 20     | 3      |
| `store.fd %F1, %A, imm32` | imm32    | F1       | imm32    | A       | 21     | 3      |
| `store.fq %F0, %A, imm32` | imm32    | F0       | imm32    | A       | 22     | 3      |
| `store.fq %F1, %A, imm32` | imm32    | F1       | imm32    | A       | 23     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 24     | 3      |
| `add.wl %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 24     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 25     | 3      |
| `and.wl %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 25     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 26     | 3      |
| `or.wl %Z, %X, imm32`     | imm32    | imm32    | Z != 0   | X       | 26     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 27     | 3      |
| `xor.wl %Z, %A, imm32`    | imm32    | imm32    | Z != 0   | X       | 27     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 28     | 3      |
| `add.dl %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 28     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 29     | 3      |
| `and.dl %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 29     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 30     | 3      |
| `or.dl %Z, %X, imm32`     | imm32    | imm32    | Z != 0   | X       | 30     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 31     | 3      |
| `xor.dl %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 31     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 32     | 3      |
| `add.ql %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 32     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 33     | 3      |
| `and.ql %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 33     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 34     | 3      |
| `or.ql %Z, %X, imm32`     | imm32    | imm32    | Z != 0   | X       | 34     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 35     | 3      |
| `xor.ql %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 35     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 36     | 3      |
| `slt.sl %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 36     | 3      |
| (unused)                  | imm32    | imm32    | 0        | *       | 37     | 3      |
| `slt.ul %Z, %X, imm32`    | imm32    | imm32    | Z != 0   | X       | 37     | 3      |
| `add.al %C, %A, imm32`    | imm32    | imm32    | C        | A       | 38     | 3      |
| ...                       | ...      | ...      | ...      | ...     | 63     | 3      |

F0 refers to a floating-point register in the low bank and F1 in the high bank.

