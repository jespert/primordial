# Executable format v1 (16-bit version)

Difficulties with ELF:

- Fairly complex.
- Way too general for statically linked 16-bit systems.
- Interpreting it correctly and securely requires a lot of code.
- NUL-terminated strings.

Difficulties with COFF:

- The presence of a timestamp is an antipattern.
- Weak architecture and ABI identification.

Difficulties with a.out:

- Weak architecture and ABI identification.

Idea:

- Common format header that ensures forward and backward compatibility.
- Use versions of increasing complexity to support different goals.
- v1 will be close to the a.out format.
- v2 will be closer in functionality to ELF.

Note that v2 doesn't completely replace v1. It just has different goals.

## Features

- Simple.
- Secure.
- Strict.
- Static.

## File header

| Position | Content     | Type    | Description                                      |
|----------|-------------|---------|--------------------------------------------------|
| 0        | magic       | U8[4]   | "EXE\0"                                          |
| 4        | fmt_version | U8      | Major format version                             |
| 5        | fmt_size    | U8      | 0,...,4 = 8, 16, 32, 64, 128 bits                |
| 6        | endianness  | U8      | 0 = little-endian, 1 = big-endian                |
| 7        | file_type   | U8      | 0 = static executable, 1 = static library        |
| 8        | arch        | U8[4]   | Architecture (packed short string, e.g., "SR16") |
| 12       | abi         | U8[4]   | ABI (packed short string, e.g., "PRIM")          |
| 16       | arch_flags  | Bit[64] | Architecture flags                               |
| 24       | abi_flags   | Bit[64] | ABI flags                                        |
| 32       | END         |         |                                                  |

While this could change in the future, this header is expected to remain
constant across all versions of the format.

Until we parse the endianness, all field types must be based on u8, which
works identically on both little-endian and big-endian machines.

We decided to splurge on the ABI and architecture fields because
4-byte packed strings are less likely to collide over time as
independent actors define their own ABI and machine names.
It allows a modicum of decentralisation.

We also splurged on the flags to better support sprawling extension sets
(think, for example, RISC-V).

## Layout (16-bits v1)

The file layout is strict and its components cannot be reordered:

1. File header
2. Main header
3. Program payload (code, ro_data, pi_data)
4. Relocation table
5. Symbol table
6. String table
7. String values.

## Main header (16-bits v1)

Very similar to the a.out format.

| Position | Field        | Type | Description                                 |
|----------|--------------|------|---------------------------------------------|
| 0        | code_size    | S16  | Size of text segment (X)                    |
| 2        | ro_data_size | S16  | Size of read-only data segment (R)          |
| 4        | pi_data_size | S16  | Size of pre-initialised writable data (RW)  |
| 6        | zi_data_size | S16  | Size of zero-initialised writable data (RW) |
| 8        | entrypoint   | S16  | Entrypoint address                          |
| 10       | num_relocs   | S16  | Number of entries in the relocation table   |
| 12       | num_symbols  | S16  | Number of entries in the symbol table       |
| 14       | num_strings  | S16  | Number of entries in the strings table      |
| 16       | END          |      |                                             |

The size of the read-write data segment is pi_data_size + zi_data_size.

## Relocation table entry (16-bits v1)

TODO.

## Symbol table entry (16-bits v1)

| Position | Field     | Type    | Description    |
|----------|-----------|---------|----------------|
| 0        | address   | S16     | Symbol address |
| 2        | type      | U16     | Symbol type    |
| 4        | flags     | Bit[16] | Symbol flags   |
| 6        | string_id | S16     | String ID      |
| 8        | END       |         |                |

## String table entry (16-bits v1)

| Position | Field      | Type | Description   |
|----------|------------|------|---------------|
| 0        | string_end | S16  | End of string |
| 2        | END        |      |               |

A string's ID is its index in the string table

The size of a string is the difference between its end and the end of the
previous entry (or 0 if it is the first entry).

The table is followed by the strings, concatenated one after another.
Delimiters are unnecessary because the size of each string is known.

Strings must be sorted with shortlex, which ensures that binary search can be
used to identify strings.
