# World's Fastest ADIF / ADI Parser?

A high-performance Go library for parsing [ADIF](https://adif.org/) (Amateur Data Interchange Format) files used in Ham Radio logging. This implementation focuses on speed and memory efficiency while providing idiomatic Go interfaces.

[![Tests](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml/badge.svg)](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamradiolog-net/adif)](https://goreportcard.com/report/github.com/hamradiolog-net/adif)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamradiolog-net/adif.svg)](https://pkg.go.dev/github.com/hamradiolog-net/adif)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/LICENSE)

Performance testing shows this library is:

- 3x to 20x faster than comparable ADI libraries
- 2x faster than standard Go JSON marshaling
- More memory efficient than other tested ADI libraries

## Usage

This library provides three ways to work with ADI files:

1) [ADIFReader](./examples/adireader_test.go): Stream-based parsing of ADI records using `io.Reader`
2) [Record](./examples/record_test.go): Single record operations using `io.Reader`/`io.Writer`
3) [Document](./examples/document_test.go): Complete ADI file operations using `io.Reader`/`io.Writer`

## Installation

```bash
go get github.com/hamradiolog-net/adif@latest
```

## Benchmarks

JSON marshaling is included as a baseline for comparison.
Note: JSON formatted data is significantly _smaller_ the same data in ADI format.
This gives the JSON marshaler a significant advantage over the ADI parsers because it has less work to perform.

| Benchmark  (AMD Ryzen 9 7950X)             | Iterations | Time/op (ns) | Bytes/op    | Allocs/op |
|--------------------------------------------|----------:|-------------:|------------:|-----------:|
| ▲ Higher is better / ▼ Lower is better     |         ▲ |            ▼ |           ▼ |          ▼ |
| **Read Operations**                        |           |              |             |            |
| This Library                               |     1,461 |      819,922 |     673,421 |      8,757 |
| JSON                                       |       622 |    1,915,720 |     402,803 |     25,601 |
| Matir                                      |       417 |    2,895,274 |   2,037,004 |     66,535 |
| Eminlin                                    |        68 |   16,453,839 |  13,127,877 |    193,083 |
| **Write Operations**                       |           |              |             |            |
| This Library                               |     2,304 |      519,218 |     514,436 |         20 |
| JSON                                       |       800 |    1,495,712 |     973,083 |     17,805 |
| Matir                                      |       399 |    2,994,459 |   1,490,840 |     28,673 |
| Eminlin                                    |       N/A |          N/A |         N/A |        N/A |

## Technical Implementation

This parser achieves high performance through the following optimizations:

### Architecture

- Implements an O(n) time complexity iterative state machine parser
- Single-pass input stream processing
- Zero-copy techniques to minimize memory operations
- Efficient buffer reuse patterns

### Performance Optimizations

- Leverages stdlib I/O operations with potential SSE/SIMD acceleration depending upon your CPU architecture
- Smart buffer pre-allocation based on discovered record sizes
- String interning for field names to reduce memory allocation and improve comparison speed
- Optimized ASCII case conversion using bitwise operations
- Custom base-10 integer parsing for ADIF field lengths

### Memory Management

- String interning for common ADI field keys to reduce allocations
- Constant memory overhead during streaming operations
- Minimal temporary allocations during field parsing
- Linear memory scaling based on record size (not file size)
- Dynamic buffer allocation based on learned field counts

## Performance Considerations

Alternative implementations were explored:

1. Full-memory reading approach:
   - 20% performance improvement over current streaming implementation
   - Rejected to maintain streaming capability for large files

2. Field list-based approach:
   - 15% faster parsing than current map-based implementation
   - 30% faster writing performance
   - More cache-friendly design
   - Rejected in favor of better API ergonomics

Future optimization opportunities include direct SIMD implementation for parsing operations.
