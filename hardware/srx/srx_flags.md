# SRX Computer (flags)

**Warning: this highly experimental and very much a work in progress.**

This is the version of SRX with flags and separate comparison instructions
(like most RISC architectures other than MIPS and RISC-V).
This design is a work in progress and may be ultimately discarded.

## Flags

This design uses the standard NZCV set of flags.
To reduce excessive data dependencies,
only specific operations can read or write flags (similar to ARM64).
