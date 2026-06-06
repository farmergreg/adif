# High Performance ADI Parser for Go

This library provides high-performance stream based processing of [ADIF](https://adif.org/) (Amateur Data Interchange Format) ADI files used for ham radio logs.

[![Tests](https://github.com/farmergreg/adif/actions/workflows/test.yml/badge.svg)](https://github.com/farmergreg/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/farmergreg/adif/v5)](https://goreportcard.com/report/github.com/farmergreg/adif/v5)
[![Go Reference](https://pkg.go.dev/badge/github.com/farmergreg/adif/v5.svg)](https://pkg.go.dev/github.com/farmergreg/adif/v5)
[![Go Version](https://img.shields.io/github/go-mod/go-version/farmergreg/adif)](https://github.com/farmergreg/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/farmergreg/adif)](https://github.com/farmergreg/adif/blob/main/LICENSE)

## Features

- **Tested**: 100% test coverage!
- **Developer Friendly**: Clean, idiomatic Go API with three usage patterns: streaming, in-memory, and writing
- **Fast**: 3x-11x faster than [other libraries](https://github.com/farmergreg/adif-benchmark)
- **Memory Efficient**: Uses 2.5x less memory and makes 21x fewer allocations than the nearest competitor.

## Quick Start

```bash
go get github.com/farmergreg/adif/v5
```

### Usage Patterns

| Type | Use when |
|------|----------|
| [`Scanner`](./scanner.go) | Streaming large files record-by-record without loading them fully into memory |
| [`Document`](./document.go) | Loading a complete ADI file into memory for random access |
| [`Writer`](./writer.go) | Writing ADI records to any `io.Writer` |

See [example_test.go](./example_test.go) for runnable examples of all three patterns.

## Benchmarks

Please see the [Go ADIF Parser Benchmarks](https://github.com/farmergreg/adif-benchmark) project for benchmarks.

TLDR, this library processes ADI data 3x faster than the go standard library can process the same data in json format.
This library is nearly 3x faster than the nearest competing ADI parser.

## Related Projects

If you found this library useful, you may also be interested in the following projects:

- [Go ADIF Parser Benchmarks](https://github.com/farmergreg/adif-benchmark)
- [Go ADIF Specification](https://github.com/farmergreg/spec)
- [g3zod/CreateADIFTestFiles](https://github.com/g3zod/CreateADIFTestFiles) ADI Test Files
- [g3zod/CreateADIFExportFiles](https://github.com/g3zod/CreateADIFExportFiles) ADIF Workgroup Specification Export Tool

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.
