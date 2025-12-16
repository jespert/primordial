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

