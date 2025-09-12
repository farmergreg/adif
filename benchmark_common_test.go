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
				p := NewADIRecordReader(reader, true)
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
	src := NewADIRecordReader(strings.NewReader(benchmarkFile), false)
	dst := NewJSONRecordWriter(&buffer)
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
	p := NewADIRecordReader(strings.NewReader(benchmarkFile), false)
	for {
		record, err := p.Next()
		if err == io.EOF {
			break
		}
		qsoListNative = append(qsoListNative, record)
	}
	return qsoListNative
}

func BenchmarkInternalParseDataLength(b *testing.B) {
	testData := []struct {
		input []byte
		want  int
	}{
		{[]byte("1"), 1},
		{[]byte("12"), 12},
		{[]byte("123"), 123},
		{[]byte("XYZMORE"), -1},
	}

	for _, td := range testData {
		b.Run(string(td.input), func(b *testing.B) {
			for b.Loop() {
				v, err := parseDataLength(td.input)
				if td.want == -1 {
					if err == nil {
						b.Fatalf("Expected error for input '%s', got nil", td.input)
					}
				} else if v != td.want {
					b.Fatalf("Expected %d, got %d", td.want, v)
				}
			}
		})
	}
}
