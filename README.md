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

1) [ADIFParser](./examples/parser_test.go): Stream records from an io.Reader.
2) [Document](./examples/document_test.go): Read all ADI records from an io.Reader into a Document.
3) [Record](./examples/record_test.go): Read a single ADI record from an io.Reader.

## Benchmarks

- Reading ADI: 329% - 2280% faster
- Writing ADI: 185% - 1088% faster

JSON marshaling is included as a baseline for comparison.
JSON formatted data tends to be significantly smaller the same data in ADI format.
So, the ADI parsers are actually doing more work than the JSON marshaler to process the same data.

| Benchmark  (AMD Ryzen 9 7950X)  | Iterations | Time/op (ns) | Bytes/op    | Allocs/op |
|---------------------------------|-----------:|-------------:|------------:|-----------:|
| **Read Operations**             |            |              |             |            |
| This Library                    |      1,748 |      690,571 |     455,364 |      8,691 |
| JSON                            |        396 |    2,963,982 |     406,286 |     16,488 |
| Matir                           |        417 |    2,895,274 |   2,037,004 |     66,535 |
| Eminlin                         |         68 |   16,453,839 |  13,127,877 |    193,083 |
| **Write Operations**            |            |              |             |            |
| This Library                    |      4,596 |      252,071 |     515,663 |         20 |
| JSON                            |      1,621 |      718,302 |     679,896 |          5 |
| Matir                           |        399 |    2,994,459 |   1,490,840 |     28,673 |
| Eminlin                         |        N/A |          N/A |         N/A |        N/A |

## Technical Implementation

This parser achieves its performance through several optimizations strategies:

### Architecture

- Implements an iterative state machine parser with O(n) time complexity.
- Processes data in a single pass through the input stream.
- Uses zero-copy techniques where possible to minimize memory allocations and copies.
- Employs buffer reuse for memory efficiency.

### Optimizations

- String interning for field names reduces memory allocations and improves comparison speed.
- Takes advantage of CPU cache locality to increase field lookup performance.
- Custom ASCII case conversion using bitwise operations.
- Specialized base-10 integer parsing optimized for ADIF field lengths.
- Utilizes standard library io operations that are likely to process multiple bytes at a time via sse/simd.
- Efficient buffer management with pre-allocation strategies to adapt to discovered record sizes.

### Memory Characteristics

- String interning ensures common ADI fields avoid allocations, ensures that keys are stored only once, and ensures lookups require fewer comparisons.
- Constant memory overhead for streaming operations.
- Minimal temporary allocations during field parsing.
- Peak memory usage scales linearly with record size, not file size.
- Allocate buffers based on learned record field counts.

These design choices result in significantly lower allocation counts and better CPU cache utilization compared to more generic parsing approaches.

## Future Improvement Thoughts

Experiments whereby the entire ADI is read into memory and io.Reader is replaced with a byte slice increased performance by about 20%.
We'd lose the ability to stream files though...

This library attempts to take advantage of the go stdlib's use of simd.
Using simd directly to further speed up parsing is an opportunity worth exploring.
