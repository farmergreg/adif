# ⚡ World's Fastest ADI Parser for Go

This library provides high-performance processing of [ADIF](https://adif.org/) (Amateur Data Interchange Format) ADI files used for ham radio logs.
It's idiomatic, developer-friendly API seamlessly integrates with your codebase and the go standard library.

[![Tests](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml/badge.svg)](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamradiolog-net/adif)](https://goreportcard.com/report/github.com/hamradiolog-net/adif)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamradiolog-net/adif.svg)](https://pkg.go.dev/github.com/hamradiolog-net/adif)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/LICENSE)

## ✨ Key Features

- 🚀 **Blazing Fast**: 3-20x faster than other ADI libraries; more than 2x faster than Go's native JSON marshaling!
- 💡 **Memory Efficient**: Up to 1400x fewer memory allocations than alternatives
- 🔧 **Developer Friendly**: Clean, idiomatic Go interfaces
- 🔬 **Battle Tested**: Comprehensive test coverage ensures reliability

## 🚀 Quick Start

```bash
go get github.com/hamradiolog-net/adif@latest
```

1) [Document](./example_document_test.go): Complete ADI file operations using `io.Reader`/`io.Writer`
2) [ADIFReader](./example_adireader_test.go): Stream-based parsing of ADI records using `io.Reader`
3) [Record](./example_record_test.go): Single record operations using `io.Reader`/`io.Writer`

## Benchmarks

JSON marshaling is included as a baseline for comparison.

| Benchmark  (AMD Ryzen 9 7950X)             | Iterations | Time/op (ns) | Bytes/op    | Allocs/op   |
|--------------------------------------------|----------:|---------------:|------------:|-----------:|
| ▲ Higher is better / ▼ Lower is better     |         ▲ |              ▼ |           ▼ |          ▼ |
| **Read Operations**                        |           |                |             |            |
| This Library                               | **1,461** |    **819,922** |   673,421   | **8,757**  |
| JSON                                       |     622   |    1,915,720   | **402,803** |   25,601   |
| Matir                                      |     417   |    2,895,274   | 2,037,004   |   66,535   |
| Eminlin                                    |      68   |   16,453,839   |13,127,877   |  193,083   |
| **Write Operations**                       |           |                |             |            |
| This Library                               | **1,800** |    **666,157** | **514,418** |     **20** |
| JSON                                       |     796   |    1,488,265   |   966,487   |   17,805   |
| Matir                                      |     399   |    2,994,459   | 1,490,840   |   28,673   |
| Eminlin                                    |     N/A   |          N/A   |       N/A   |      N/A   |

## 🔧 Technical Deep Dive

This parser achieves high performance through the following optimizations:

### Architecture

- Implements an O(n) time complexity streaming parser
- Zero-copy techniques to minimize memory operations
- Efficient buffer reuse patterns
- Simple API

### Performance Optimizations

- Leverages stdlib I/O operations with SSE/SIMD acceleration depending upon your CPU architecture
- Smart buffer pre-allocation based on discovered record sizes
- Optimized ASCII case conversion using bitwise operations
- Custom base-10 integer parsing for ADIF field lengths

### Memory Management

- Minimal temporary allocations during field parsing
- String interning for common ADI field names to reduce allocations and memory use
- Constant memory overhead during streaming operations
- Dynamic buffer allocation based on learned field counts

## Performance Considerations

Alternative implementations explored:

1. Non-Streaming / Read ADI File into Memory:
   - 20% performance improvement over current streaming implementation
   - Rejected to maintain streaming capability for large data sets

2. slice of fields instead of map:
   - L1/2/3 cache-friendly design
   - O(n) lookup time instead of map O(1) lookup time.
   - For small values of n where n is less than ~40, this is faster than the map based implementation.
   - 15% faster parsing than current map-based implementation on test files.
   - 30% faster writing performance (it is much faster to iterate over a list than a map)
   - Rejected in favor of better API ergonomics
   - Still considering switching to this implementation because it is common to see ADI records with less than 40 fields.
   - Downside is that above ~50 fields, the O(n) map lookup is faster than the O(n) list lookup.

Future optimization possibilities:

- A perfect hash function with a custom map implementation.
- Direct SIMD implementation for parsing operations.

## Related Projects

If you found this library useful, you may also be interested in the following projects:

- [hamradiolog-net/adif-spec](https://github.com/hamradiolog-net/adif-spec): ADIF Specification For Go
- [g3zod/CreateADIFTestFiles](https://github.com/g3zod/CreateADIFTestFiles) ADI Test Files
- [g3zod/CreateADIFExportFiles](https://github.com/g3zod/CreateADIFExportFiles) ADIF Workgroup Specification Export Tool

## 📝 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.
