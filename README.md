# World's Fastest ADI Parser?

[![Tests](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml/badge.svg)](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamradiolog-net/adif)](https://goreportcard.com/report/github.com/hamradiolog-net/adif)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamradiolog-net/adif.svg)](https://pkg.go.dev/github.com/hamradiolog-net/adif)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/LICENSE)

This ADI parser is my attempt to create a fast, efficient ADIF parser for ADI formatted data.

This library also outperforms other ADI libraries that I've tested to date.
Additionally, this library is able to convert to and from ADI faster than the built-in go json library can convert the same data to and from JSON.
This even though JSON is a much more compact format than ADI.

## Usage

This library provides two ways to parse ADI files: using the low-level `ADIFParser` directly or using the higher-level `Document` type.
The interfaces and implementations are designed to be idiomatic go and work well with the go standard library.

The unit tests provide examples of how to use the library.

- [ADIFParser](https://github.com/hamradiolog-net/adif/blob/main/adiparser_test.go)
- [Document](https://github.com/hamradiolog-net/adif/blob/main/document_test.go)

## Benchmarks

- Reading ADI: 300% - 2180% Faster
- Writing ADI: 180% - 1100% Faster

```
Benchmark                                                                 Iterations          Time/op

cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkReadThisLibrary-32                                                 1626            714559 ns/op
BenchmarkReadJSON-32                                                         409           2903384 ns/op
BenchmarkReadMatir-32                                                        416           2873248 ns/op
BenchmarkReadEminlin-32                                                       70          16312785 ns/op

BenchmarkWriteThisLibrary-32                                                4488            245656 ns/op
BenchmarkWriteJSON-32                                                       1666            694251 ns/op
BenchmarkWriteMatir-32                                                       408           2921377 ns/op
BenchmarkWriteEminlin-32                                                     N/A               N/A
```

## Future Work

How could this library be faster yet?
Internally, this library attempts to take advantage of the go stdlib's use of simd.
However, I think there is an opportunity to use SIMD directly to further speed up parsing.
