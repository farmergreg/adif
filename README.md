# World's Fastest ADI Parser?

[![Tests](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml/badge.svg)](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamradiolog-net/adif)](https://goreportcard.com/report/github.com/hamradiolog-net/adif)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamradiolog-net/adif.svg)](https://pkg.go.dev/github.com/hamradiolog-net/adif)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/LICENSE)

This ADI parser is my attempt to create a fast, memory efficient ADIF parser for ADI formatted data.

This library outperforms other ADI libraries that I've tested to date by a wide margin.

## Usage

This library provides two ways to parse ADI files: using [ADIFParser](https://github.com/hamradiolog-net/adif/blob/main/adiparser.go) ([examples](https://github.com/hamradiolog-net/adif/blob/main/adiparser_test.go)) to stream records, or using the [Document](https://github.com/hamradiolog-net/adif/blob/main/document.go) ([examples](https://github.com/hamradiolog-net/adif/blob/main/document_test.go)) type to load the entire file into memory.
Both are simple to use and implement standard interfaces that make them easy to use with the go library.

## Benchmarks

- Reading ADI: 329% - 2280% faster
- Writing ADI: 185% - 1088% faster

I've included JSON marshaling as a baseline for comparison.
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

## Future Improvement Thoughts

Experiments whereby the entire ADI is read into memory increased performance by about 20%.
I've not yet implemented this cleanly, but the gains are clearly there.
We'd loose the ability to stream files though...

This library attempts to take advantage of the go stdlib's use of simd.
I think there is an opportunity to use simd directly to further speed up parsing.
