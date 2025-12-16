# SRX Computer

**Warning: this highly experimental and very much a work in progress.**

Specification and emulator for the SRX computer architecture.

SRX evolves from SR16 to support 16, 32, and 64-bit architectures.
It uses a variable-length instruction set,
whose instructions can be 16, 32, or 48-bits-long.
Its immediates can be 8, 16, and 32-bits-long.

128-bit support is desirable but not required.

Compared to SR16, SRX supports aims to efficiently support much bigger address
spaces and data types.

Compared to RISC-V, SRX should ideally support more natural immediates,
longer control offsets, and ideally even higher code density.

SRX is roughly based on SR16, but the instruction sets are not compatible.

## Versions

Because long immediates take a sizeable chunk of the encoding space, which is
why most RISC architectures with 32 registers have avoided them, we will have
to try a few things to see what works best. The main two avenues are:

- [Test-and-branch](srx_tab.md), like RISC-V.
- [Condition flags](srx_flags.md), like ARM, x86/x64, and PowerPC.

These designs could end up being quite different, so we will explore them in
separate files.

This file describes their common features.

## Registers

The convention for saved registers (S) grows downwards to mitigate the risk of
off-by-one errors between the alias and the register number.

| Type | Register | Alias | Purpose                                        |
|------|----------|-------|------------------------------------------------|
| D    | 0        | ZR    | Zero (hardcoded)                               |
| D    | 1        |       | Argument register 4                            |
| D    | 2        |       | Argument register 3                            |
| D    | 3        |       | Argument register 2                            |
| D    | 4        |       | Argument register 1                            |
| D    | 5        |       | Argument register 0                            |
| D    | 6        |       | Temporary register 1                           |
| D    | 7        |       | Temporary register 0                           |
| D    | 8        |       | Saved register 0                               |
| D    | 9        |       | Saved register 1                               |
| D    | a        |       | Saved register 2                               |
| D    | b        |       | Saved register 3                               |
| D    | c        |       | Saved register 4                               |
| D    | d        |       | Saved register 5                               |
| D    | e        |       | Saved register 6                               |
| D    | f        |       | Saved register 7                               |
| A    | 0        | SP    | Stack pointer                                  |
| A    | 1        | BP    | Base pointer                                   |
| A    | 2        | TP    | Thread pointer                                 |
| A    | 3        |       | Argument register 2                            |
| A    | 4        |       | Argument register 1                            |
| A    | 5        |       | Argument register 0                            |
| A    | 6        |       | Temporary register 1                           |
| A    | 7        | MP    | Temporary register 0, millicode return pointer |
| A    | 8        | FP    | Saved register 0, Frame pointer                |
| A    | 9        |       | Saved register 1                               |
| A    | a        |       | Saved register 2                               |
| A    | b        |       | Saved register 3                               |
| A    | c        |       | Saved register 4                               |
| A    | d        |       | Saved register 5                               |
| A    | e        |       | Saved register 6                               |
| A    | f        | RP    | Return pointer                                 |
