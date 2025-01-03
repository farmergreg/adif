# World's Fastest ADI Parser?

[![Tests](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml/badge.svg)](https://github.com/hamradiolog-net/adif/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamradiolog-net/adif)](https://goreportcard.com/report/github.com/hamradiolog-net/adif)
[![Go Reference](https://pkg.go.dev/badge/github.com/hamradiolog-net/adif.svg)](https://pkg.go.dev/github.com/hamradiolog-net/adif)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/go.mod)
[![License](https://img.shields.io/github/license/hamradiolog-net/adif)](https://github.com/hamradiolog-net/adif/blob/main/LICENSE)

This ADI parser is my attempt to create a fast, efficient ADIF parser for ADI formatted data.

This library outperforms other ADI libraries that I've tested to date by a wide margin.

## Usage

This library provides two ways to parse ADI files: using [ADIFParser](https://github.com/hamradiolog-net/adif/blob/main/adiparser.go) ([examples](https://github.com/hamradiolog-net/adif/blob/main/adiparser_test.go)) to stream records, or using the [Document](https://github.com/hamradiolog-net/adif/blob/main/document.go) ([examples](https://github.com/hamradiolog-net/adif/blob/main/document_test.go)) type to load the entire file into memory.
Both are simple to use and implement standard interfaces that make them easy to use with the go library.

## Benchmarks

- Reading ADI: 300% - 2100% faster
- Writing ADI: 190% - 1100% faster

I've included JSON Marshaling as a baseline for comparison.
The benchmarks marshal the same data to and from JSON and ADIF ADI formatted data.
Note that JSON tends to be significantly smaller than ADI.
Therefore, the performance gains over JSON are even more significant when considering the size of the data.

| Benchmark  (AMD Ryzen 9 7950X)  | Iterations | Time/op (ns) | Bytes/op    | Allocs/op |
|-----------------------------------------|-----------:|-------------:|------------:|-----------:|
| **Read Operations**             |            |              |             |            |
| This Library                    |      1,666 |      710,337 |     455,364 |      8,691 |
| JSON                            |        423 |    2,805,019 |     406,133 |     16,488 |
| Matir                           |        434 |    2,761,990 |   2,037,276 |     66,536 |
| Eminlin                         |         76 |   15,834,453 |  13,143,713 |    193,085 |
| **Write Operations**            |            |              |             |            |
| This Library                    |      5,181 |      231,122 |     515,679 |         20 |
| JSON                            |      1,766 |      675,473 |     669,011 |          5 |
| Matir                           |        422 |    2,845,072 |   1,490,899 |     28,673 |
| Eminlin                         |        N/A |          N/A |         N/A |        N/A |

## Future Work Ideas

Experiments whereby the entire ADI is read into memory increased performance by about 20%.
I've not yet implemented this cleanly, but the gains are clearly there.
We'd loose the ability to stream files though...

This library attempts to take advantage of the go stdlib's use of simd.
I think there is an opportunity to use simd directly to further speed up parsing.
