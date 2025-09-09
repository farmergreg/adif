package adif

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
	"testing"
)

//go:embed testdata/N3FJP-AClogAdif.adi
var benchmarkFile string

func BenchmarkParseAllADIFiles(b *testing.B) {
	fs, err := testFileFS.ReadDir("testdata")
	if err != nil {
		b.Fatal(err)
	}

	for _, f := range fs {
		b.Run(f.Name(), func(b *testing.B) {
			for b.Loop() {
				reader, _ := testFileFS.Open("testdata/" + f.Name())
				p := NewADIReader(reader, true)
				for {
					_, err := p.Next()
					if err == io.EOF {
						break
					}
				}
				reader.Close()
			}
		})

	}
}

func benchmarkFileAsJSON() []byte {
	var buffer bytes.Buffer
	src := NewADIReader(strings.NewReader(benchmarkFile), false)
	dst := NewADIJWriter(&buffer)
	srcRecords := make([]ADIFRecord, 0, 10000)
	for {
		record, err := src.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		srcRecords = append(srcRecords, record)
	}
	dst.Write(srcRecords)
	return buffer.Bytes()
}

func loadTestData() []ADIFRecord {
	var qsoListNative []ADIFRecord
	p := NewADIReader(strings.NewReader(benchmarkFile), false)
	for {
		record, err := p.Next()
		if err == io.EOF {
			break
		}
		qsoListNative = append(qsoListNative, record)
	}
	return qsoListNative
}
