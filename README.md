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

This library provides two ways to parse ADI files: using [ADIFParser](https://github.com/hamradiolog-net/adif/blob/main/adiparser_test.go) to stream records, or using the [Document](https://github.com/hamradiolog-net/adif/blob/main/document_test.go) type to load the entire file into memory.
Both are easy to use and implement standard interfaces that make them easy to use with the go library.

## Benchmarks

- Reading ADI: 300% - 2180% Faster
- Writing ADI: 180% - 1100% Faster

| Benchmark  (AMD Ryzen 9 7950X)          | Iterations | Time/op (ns) |
|-----------------------------------------|-----------:|-------------:|
| **Read Operations**                     |            |              |
| This Library                            | 1,626      | 714,559      |
| JSON                                    | 409        | 2,903,384    |
| Matir                                   | 416        | 2,873,248    |
| Eminlin                                 | 70         | 16,312,785   |
| **Write Operations**                    |            |              |
| This Library                            | 4,488      | 245,656      |
| JSON                                    | 1,666      | 694,251      |
| Matir                                   | 408        | 2,921,377    |
| Eminlin                                 | N/A        | N/A          |

## Future Work

How could this library be faster yet?
Internally, this library attempts to take advantage of the go stdlib's use of simd.
However, I think there is an opportunity to use SIMD directly to further speed up parsing.
