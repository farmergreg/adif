# World's Fastest ADI Parser?

[![Tests](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml/badge.svg)](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamradiolog-net/adif)](https://goreportcard.com/report/github.com/hamradiolog-net/adif)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamradiolog-net/adif.svg)](https://pkg.go.dev/github.com/hamradiolog-net/adif)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/LICENSE)

This ADI parser is an attempt to create a fast, memory efficient ADIF parser for ADI formatted data.

This library outperforms other ADI libraries that I've tested to date by a wide margin.

## Usage

This library provides three ways to work with ADI files:

1) [ADIFReader](./examples/adireader_test.go): Parse a stream of ADI records via an io.Reader.
2) [Document](./examples/document_test.go): Read all ADI records from an io.Reader into a Document.
3) [Record](./examples/record_test.go): Read a single ADI record from an io.Reader.

## Benchmarks

- Reading ADI: 134% - 1907% faster, significantly fewer allocations
- Writing ADI: 107% -  314% faster, significantly fewer allocations

JSON marshaling is included as a baseline for comparison.
JSON formatted data tends to be significantly smaller the same data in ADI format.
This gives the JSON marshaler an advantage over the ADI parsers.

| Benchmark  (AMD Ryzen 9 7950X)             | Iterations | Time/op (ns) | Bytes/op    | Allocs/op |
|--------------------------------------------|----------:|-------------:|------------:|-----------:|
| ▲ Higher is better / ▼ Lower is better     |         ▲ |            ▼ |           ▼ |          ▼ |
| **Read Operations**                        |           |              |             |            |
| This Library                               |     1,461 |      819,922 |     673,421 |      8,757 |
| JSON                                       |       622 |    1,915,720 |     402,803 |     25,601 |
| Matir                                      |       417 |    2,895,274 |   2,037,004 |     66,535 |
| Eminlin                                    |        68 |   16,453,839 |  13,127,877 |    193,083 |
| **Write Operations**                       |           |              |             |            |
| This Library                               |     1,671 |      723,542 |     515,507 |         21 |
| JSON                                       |       800 |    1,495,712 |     973,083 |     17,805 |
| Matir                                      |       399 |    2,994,459 |   1,490,840 |     28,673 |
| Eminlin                                    |       N/A |          N/A |         N/A |        N/A |

## Technical Implementation

This parser achieves high performance through the following optimizations:

### Architecture

- Implements an iterative state machine parser with O(n) time complexity.
- Processes data in a single pass through the input stream.
- Uses zero-copy techniques where possible to minimize memory allocations and copies.
- Employs buffer reuse for memory efficiency.

### Optimizations

- Utilizes standard library io operations that are likely to process multiple bytes at a time via sse/simd.
- Efficient buffer management with pre-allocation strategies to adapt to discovered record sizes.
- String interning for field names reduces memory allocations and improves comparison speed.
- Custom ASCII case conversion using bitwise operations.
- Specialized base-10 integer parsing optimized for ADIF field lengths.

### Memory Characteristics

- String interning ensures common ADI fields string keys avoids additional allocations, reduces overall memory use.
- Constant memory overhead for streaming operations.
- Minimal temporary allocations during field parsing.
- Peak memory usage scales linearly with record size, not file size.
- Allocate buffers based on learned record field counts.

## Going Faster

Experiments whereby the entire ADI is read into memory and io.Reader is replaced with a byte slice increased performance by about 20% when compared to the current implementation.
I thought it important to maintain the ability to stream files and therefore did not include this optimization.

An alternative approach using a list of field names and a field struct with a name and value field is about 15% faster than the current parser.
It also improved upon adi write performance by about 30%.
I think it is faster because it is more L1/2/3 cache friendly.
I kept the current implementation which uses maps because I think they are a better _user_ interface for this library.
It also serializes easily to a nice looking json format.
I am also concerned that there may be cases where the O(n) lookup time becomes a problem with large records.
Bechmarks showed that the map implementation beats list lookup by 2x for records with ~50 fields when looking up a worst case scenario field.
With the map, we ensure better average performance of O(1) for large records that contain many fields.
This is one of those FUN situations where an O(n) algorithm is faster than O(1) for small values of n!

This library attempts to take advantage of the go stdlib's use of simd.
Using simd directly to further speed up parsing is an opportunity worth exploring.

Writing performance can be improved by about 50% by removing the field sorting in WriteTo.
I feel that sorting the fields provides a better debugging / user experience and is worth the trade-off.
We are still 2x faster than the next fastest benchmarked library.
If needed, we could implement a WriteToFast method that does not sort the fields.
